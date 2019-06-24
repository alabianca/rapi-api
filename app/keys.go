package app

import (
	"context"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-chi/chi"

	"github.com/alabianca/rapi-api/utils"
)

var CheckKey = func(next http.Handler) http.Handler {
	needsKey := "/pub/v1/record"
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

		record := chi.URLParam(r, "resumeID")

		if record == "" {
			utils.Respond(w, utils.Message(http.StatusUnauthorized, "Resume ID is required"))
			return
		}

		resumeID, err := primitive.ObjectIDFromHex(record)
		if err != nil {
			utils.Respond(w, utils.Message(http.StatusUnauthorized, "Resume ID is required "+err.Error()))
		}

		db, err := models.GetDB()
		if err != nil {
			utils.Respond(w, utils.Message(http.StatusInternalServerError, "Could not get a handle on db"))
			return
		}

		var apiKey models.APIKey
		keys := db.Collection(models.KeyCollection)
		filter := bson.D{{"key", key}}

		if err := keys.FindOne(context.TODO(), filter).Decode(&apiKey); err != nil {
			utils.Respond(w, utils.Message(http.StatusInternalServerError, err.Error()))
			return
		}

		if apiKey.Resume != resumeID {
			utils.Respond(w, utils.Message(http.StatusUnauthorized, "Resume IDs don't match up"))
			return
		}

		log.Printf("Resume %s\n", resumeID.Hex())
		log.Printf("API Key %s\n", apiKey.Key)

		ctx := context.WithValue(r.Context(), "resume", resumeID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
