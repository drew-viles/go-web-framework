package auth

import (
	"encoding/json"
	"fmt"
	"gitea.viles.uk/dcp/web-framework/environment"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

// CreateToken Creates a JSON Web Token.
func CreateToken(hashID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["hash_id"] = hashID.String()
	claims["expr"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(environment.GetAPISecret())
}

// ExtractTokenUserID will extract the users id stored within the token
func ExtractTokenUserID(r *http.Request) (uuid.UUID, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return environment.GetAPISecret(), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var uid uuid.UUID
		uid, err = uuid.Parse(claims["hash_id"].(string))
		if err != nil {
			return uuid.Nil, err
		}
		return uid, nil
	}

	return uuid.Nil, nil
}

// TokenValid takes a http request, checks for a token and validates it.
func TokenValid(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return environment.GetAPISecret(), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		pretty(claims)
	}

	return nil
}

// extractToken extracts a token from a http request
func extractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// pretty just prints things in a nice output
func pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("user authenticated: \n%s\n", string(b))
}
