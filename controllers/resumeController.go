package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

var CreateResume = func(w http.ResponseWriter, r *http.Request) {
	resume := &models.Resume{}

	if err := json.NewDecoder(r.Body).Decode(resume); err != nil {
		utils.Respond(w, utils.Message(http.StatusInternalServerError, "Json Decode Error"))
		return
	}

	response := models.GenerateRandomURL()

	utils.Respond(w, response)
}
