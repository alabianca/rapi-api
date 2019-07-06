package controllers

import (
	"net/http"

	"github.com/alabianca/rapi-api/utils"
)

func (a *API) GetProjects(w http.ResponseWriter, r *http.Request) {
	resumeID, err := GetResumeIdFromContext(r)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Resume ID is required"))
	}

	res := a.DAL.Projects().GetProjects(resumeID)

	utils.Respond(w, res)
}
