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

package utils

import (
	"errors"
	"strings"
)

// FormatError just formats things nicely if certain keywords are picked up.
func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("username is already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("email address is already taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect username or password")
	}
	return errors.New(err)
}
