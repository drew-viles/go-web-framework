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
