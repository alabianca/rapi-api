package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

var CreateToken = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "JSON Decode Error"))
		return
	}

	resp := models.GetToken(user.Email, user.Password)

	utils.Respond(w, resp)
}
