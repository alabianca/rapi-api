package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/utils"

	"github.com/alabianca/rapi-api/models"
)

var CreateKey = func(w http.ResponseWriter, r *http.Request) {
	var key models.APIKey
	log.Println("in here")
	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, err.Error()))
		return
	}

	userID := r.Context().Value("user").(string)
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "User ID is required"))
		return
	}

	resume := chi.URLParam(r, "resumeID")
	resumeID, err := primitive.ObjectIDFromHex(resume)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "Resume ID is required"))
		return
	}

	key.Resume = resumeID
	key.UserID = id

	res := key.Create()

	utils.Respond(w, res)
}

var GetKeys = func(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(string)
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "User ID is required"))
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Resume ID is required"))
		return
	}

	resume, err := primitive.ObjectIDFromHex(resumeID)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Resume ID is required"))
		return
	}

	response := models.GetKeys(id, resume)

	utils.Respond(w, response)

}
