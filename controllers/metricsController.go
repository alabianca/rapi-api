package controllers

import (
	"net/http"
	"strings"

	"github.com/alabianca/rapi-api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) GetMetrics(w http.ResponseWriter, r *http.Request) {
	keyArgs := r.URL.Query().Get("keys")
	ids := objectIdsFromStrings(keyArgs)

	response := a.DAL.Logs().GetLogsForKeys(ids)

	utils.Respond(w, response)
}

func objectIdsFromStrings(str string) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, 0)
	keys := strings.Split(str, ",")

	for _, key := range keys {
		id, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return ids
}
