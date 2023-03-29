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

package middleware

import (
	"errors"
	"github.com/drew-viles/go-web-framework/auth"
	"github.com/drew-viles/go-web-framework/responses"
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
