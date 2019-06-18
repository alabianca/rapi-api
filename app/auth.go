package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

type notAuth struct {
	url    string
	method string
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// list of endpoinst that don't require auth tokens
		notAuth := []notAuth{
			{
				url:    "/v1/api/user",
				method: "POST",
			},
			{
				url:    "/v1/api/token",
				method: "POST",
			},
		}
		requestPath := r.URL.Path // current request path
		requestMethod := r.Method

		// check if request does not need auth, serve the request if it does

		for _, value := range notAuth {
			if requestMethod == "OPTIONS" || (value.url == requestPath && value.method == requestMethod) {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = utils.Message(http.StatusUnauthorized, "Missing auth token")
			utils.Respond(w, response)
			return
		}

		split := strings.Split(tokenHeader, " ") // The token normaly comes in format `Bearer {token}`

		if len(split) != 2 {
			response = utils.Message(http.StatusUnauthorized, "Invalid/Malformed auth token")
			utils.Respond(w, response)
			return
		}

		tokenPart := split[1] // grab the token part
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			msg := validationErrorMessage(err)
			response = utils.Message(http.StatusUnauthorized, msg)
			utils.Respond(w, response)
			return
		}

		if !token.Valid {
			response = utils.Message(http.StatusUnauthorized, "Invalid Token")
			utils.Respond(w, response)
			return
		}

		// Everything went well, proceed with the request and set the caller
		// to the user retrieved from teh parsed token
		log.Printf("User %s", tk.UserID)
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r) // proceed the middleware chain

	})
}

func validationErrorMessage(err error) string {
	valErr := err.(*jwt.ValidationError)

	var msg string
	switch valErr.Errors {
	case jwt.ValidationErrorExpired:
		msg = "Token is expired"
	default:
		msg = "Token is malformed"
	}

	return msg
}
