package controllers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

func (a *API) CreateResume(w http.ResponseWriter, r *http.Request) {
	resume := &models.Resume{}
	user := r.Context().Value("user").(string)

	if err := json.NewDecoder(r.Body).Decode(resume); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error"))
		return
	}

	userId, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "UserID is required"))
		return
	}

	resume.UserID = userId

	response := a.DAL.Resume().CreateResume(resume)

	utils.Respond(w, response)
}

func (a *API) GetResumes(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)

	userId, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		utils.Respond(w, utils.Message(http.StatusUnauthorized, "UserID is required"))
		return
	}

	res := a.DAL.Resume().GetResumes(userId)

	utils.Respond(w, res)

}
