package controllers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

var CreateResume = func(w http.ResponseWriter, r *http.Request) {
	resume := &models.Resume{}
	user := r.Context().Value("user").(string)

	if err := json.NewDecoder(r.Body).Decode(resume); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error"))
		return
	}

	userId, err := primitive.ObjectIDFromHex(user)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "Not Authorized"))
		return
	}

	response := models.GenerateRandomURL(userId)

	utils.Respond(w, response)
}
