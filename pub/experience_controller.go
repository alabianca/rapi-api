package pub

import (
	"net/http"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
)

var GetExperience = func(w http.ResponseWriter, r *http.Request) {
	resumeID, err := GetResumeIdFromContext(r)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Resume ID is required"))
	}

	res := models.GetExperience(resumeID)

	utils.Respond(w, res)
}
