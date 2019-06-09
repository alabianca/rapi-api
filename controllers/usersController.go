package controllers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi"

	"github.com/alabianca/rapi-api/utils"

	"github.com/alabianca/rapi-api/models"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error "+err.Error()))
		return
	}

	resp := user.Create()

	utils.Respond(w, resp)
}

var AuthenticateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error"))
		return
	}

	resp := models.Login(user.Email, user.Password)

	utils.Respond(w, resp)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		utils.Message(http.StatusNotFound, "Could Not Retrieve ID "+userID)
	}

	resp := models.GetUserById(id)

	utils.Respond(w, resp)
}
