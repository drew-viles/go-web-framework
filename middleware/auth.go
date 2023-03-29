package middleware

import (
	"errors"
	"gitea.viles.uk/dcp/web-framework/auth"
	"gitea.viles.uk/dcp/web-framework/responses"
	"log"
	"net/http"
)

// SetMiddlewareAuthentication allows the request to continue if the provided user access level is greater than the route's access level
func SetMiddlewareAuthentication(next http.HandlerFunc, userAccessLevel, routeAccessLevel int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userAccessLevel < routeAccessLevel {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(w, r)
	}
}

// SetMiddlewareAuthorisation allows the request to continue if a provided jwt is valid.
func SetMiddlewareAuthorisation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			log.Printf("error vaildating token: %s\n", err)
			responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(w, r)
	}
}
