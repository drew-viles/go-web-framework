/*
Copyright 2023 Drew Viles.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"bytes"
	"github.com/drew-viles/go-web-framework/environment"
	"github.com/drew-viles/go-web-framework/erroring"
	"github.com/drew-viles/go-web-framework/routing"
	"github.com/drew-viles/go-web-framework/validation"
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
