package app

import (
	"bytes"
	"gitea.viles.uk/dcp/web-framework/environment"
	"gitea.viles.uk/dcp/web-framework/erroring"
	"gitea.viles.uk/dcp/web-framework/routing"
	"gitea.viles.uk/dcp/web-framework/validation"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"path"
)

// Server is used to store router, config and validation info
type Server struct {
	Router    *mux.Router
	Validator *validation.Validator
	Config    *environment.ConfigMap
}

// Initialise create a new Gorilla mux router and initialises the []routing.Route passed into it
func (s *Server) Initialise(routes *[]routing.Route) {
	s.Router = mux.NewRouter()
	s.initialiseRoutes(routes)
}

func (s *Server) initialiseRoutes(routes *[]routing.Route) {
	routing.SetupRoutes(routes, s.Router)
}

// InterfaceWithAPI creates a http client to send a request to another URL returning an array of bytes as the response.
func (s *Server) InterfaceWithAPI(url string, method string, inputData []byte) (result []byte, err error) {
	var req *http.Request
	var res *http.Response

	apiURL := path.Join(s.Config.Api.ApiEndpoint, url)
	client := &http.Client{}

	req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(inputData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return body, erroring.NotFoundError
	}

	return body, nil
}
