package app

import (
	"log"
	"net/http"
	"time"

	"github.com/alabianca/rapi-api/models"
)

var LogKey = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Context().Value("apiKey").(models.APIKey)

		logEvent := &models.Log{
			APIID: apiKey.ID,
			Date:  time.Now().UTC(),
		}

		res := models.CreateLog(logEvent)
		log.Printf("Logged API Request to public endpoint %v", res)

		next.ServeHTTP(w, r)
	})

}
