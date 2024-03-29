package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/utils"

	"github.com/alabianca/rapi-api/models"
)

func (a *API) CreateKey(w http.ResponseWriter, r *http.Request) {
	var key models.APIKey

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

	res := a.DAL.Keys().CreateKey(&key)

	utils.Respond(w, res)
}

func (a *API) GetKeys(w http.ResponseWriter, r *http.Request) {
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

	response := a.DAL.Keys().GetKeys(id, resume)

	utils.Respond(w, response)

}

func (a *API) PatchKey(w http.ResponseWriter, r *http.Request) {
	var key models.APIKey
	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Json Decode Error"))
		return
	}

	response := a.DAL.Keys().UpdateKey(&key)

	utils.Respond(w, response)
}

func (a *API) DeleteKey(w http.ResponseWriter, r *http.Request) {
	keyId := chi.URLParam(r, "keyID")
	if keyId == "" {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Key ID is required"))
		return
	}

	id, err := primitive.ObjectIDFromHex(keyId)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, err.Error()))
		return
	}

	response := a.DAL.Keys().DeleteKey(id)
	utils.Respond(w, response)
}
