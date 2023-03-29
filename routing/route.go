package routing

import (
	"fmt"
	"gitea.viles.uk/dcp/web-framework/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Route is used to store the information f a single route.
type Route struct {
	Name                   string
	Description            string
	Path                   string
	HandlerFunc            http.HandlerFunc
	RequestMethod          string
	ContentType            string
	IsStaticPath           bool
	RequiresAuthorisation  bool
	RequiresAuthentication bool
	AccessLevel            int
	HasJSONResponse        bool
	QueryParams            []string
}

// SetupRoutes takes an array of Route and creates a set of routes for the mux router. It will add any middleware, static paths and more as required.
// It supports authenticated and unauthenticated routes.
func SetupRoutes(routes *[]Route, router *mux.Router) {
	hasStaticPaths := false
	var staticPaths []Route
	for _, route := range *routes {
		if route.IsStaticPath {
			hasStaticPaths = true
			staticPaths = append(staticPaths, route)
			continue
		} else {
			var routeHandler http.HandlerFunc
			logMessage := "Setting up"

			routeHandler = route.HandlerFunc

			if route.RequiresAuthorisation {
				logMessage = fmt.Sprintf("%s AUTHENTICATED %s Route: %s, on path: %s", logMessage, route.RequestMethod, route.Name, route.Path)
				routeHandler = middleware.SetMiddlewareAuthorisation(routeHandler)
			} else {
				logMessage = fmt.Sprintf("%s UNAUTHENTICATED %s Route: %s, on path: %s", logMessage, route.RequestMethod, route.Name, route.Path)
			}

			if route.AccessLevel > 0 {
				logMessage = fmt.Sprintf("%s WITH access level %d", logMessage, route.AccessLevel)
				routeHandler = middleware.SetMiddlewareAuthentication(routeHandler, route.AccessLevel, route.AccessLevel)
			} else {
				logMessage = fmt.Sprintf("%s WITHOUT an access level", logMessage)
			}

			// HasJSONResponse must be the last check
			if route.HasJSONResponse {
				routeHandler = middleware.SetMiddlewareJSON(routeHandler)
			} else {
				router.Headers("Content-Type", "text/html")
			}

			if len(route.QueryParams) > 0 {
				logMessage = fmt.Sprintf("%s has query params: %s", logMessage, route.QueryParams)
				router.Path(route.Path).Queries(route.QueryParams...).HandlerFunc(routeHandler).Methods(route.RequestMethod)
				log.Println(logMessage)
				continue
			}
			router.HandleFunc(route.Path, routeHandler).Methods(route.RequestMethod)
			log.Println(logMessage)
		}
	}

	if hasStaticPaths {
		rootPath := "/public/"

		for _, route := range staticPaths {

			pathPrefix := "/" + route.Path + "/"
			pathValue := rootPath + route.Path

			handler := http.StripPrefix(pathPrefix,
				http.FileServer(http.Dir("."+pathValue)))
			route.HandlerFunc = handler.ServeHTTP
			contentType := "text/html"

			if route.ContentType != "" {
				contentType = route.ContentType
			}
			router.Headers("Content-Type", contentType)

			logMessage := fmt.Sprintf("Setting up STATIC %s Route: %s, on path: %s", route.RequestMethod, route.Name, pathPrefix)
			log.Println(logMessage)

			router.PathPrefix(pathPrefix).Handler(handler).Methods(route.RequestMethod)
		}
	}
}
