package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-chi/chi"

	"github.com/alabianca/rapi-api/utils"
)

var ScopeMapping = map[string]string{
	"POST": "create",
	"GET":  "read",
}

var CheckKey = func(next http.Handler) http.Handler {
	needsKey := "/pub/v1/record"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		method := r.Method

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

		// check if the api key has the required scope for the request
		if !scopeMapping(method, apiKey.Scope) {
			utils.Respond(w, utils.Message(http.StatusUnauthorized, fmt.Sprintf("This API Key is not authorized to execute %s requests", method)))
			return
		}

		log.Printf("Resume %s\n", resumeID.Hex())
		log.Printf("API Key %s\n", apiKey.Key)

		ctx := context.WithValue(r.Context(), "resume", resumeID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func scopeMapping(method string, scope []string) bool {
	requiredScope := ScopeMapping[method]

	for _, s := range scope {
		if s == requiredScope {
			return true
		}
	}

	return false
}
