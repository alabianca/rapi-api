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
	reg := &models.Registration{}
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error "+err.Error()))
		return
	}

	if reg.Password != reg.Verify {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Passwords Do Not Match"))
		return
	}

	user.Email = reg.Email
	user.Password = reg.Password

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
		utils.Respond(w, utils.Message(http.StatusNotFound, "Could Not Retrieve ID "+userID))
		return
	}

	resp := models.GetUserById(id)

	utils.Respond(w, resp)
}

var PostUser = func(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusNotFound, "Could Not Retrieve ID "+userID))
		return
	}

	record := &models.URLRecord{}

	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "JSON Decode Error"))
		return
	}

	resp := models.AddRecord(id, record.ID)

	utils.Respond(w, resp)
}

var GetRecordsForUser = func(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusNotFound, "Could Not Retrieve ID "+userID))
		return
	}

	resp := models.GetRecords(id)

	utils.Respond(w, resp)

}
