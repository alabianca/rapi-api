package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

func (a *API) CreateToken(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "JSON Decode Error"))
		return
	}

	resp := a.DAL.Tokens().GetToken(user.Email, user.Password)

	utils.Respond(w, resp)
}
