package app

import (
	"net/http"
	"strings"

	"github.com/alabianca/rapi-api/utils"
)

var CheckKey = func(next http.Handler) http.Handler {
	needsKey := "/v1/record"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		if !strings.Contains(requestPath, needsKey) {
			next.ServeHTTP(w, r)
			return
		}

		key := r.URL.Query().Get("key")

		if key == "" {
			utils.Respond(w, utils.Message(http.StatusUnauthorized, "API Key is required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
