package controllers

import (
	"fmt"
	"net/http"

	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) GetResume(w http.ResponseWriter, r *http.Request) {

	resumeID, err := GetResumeIdFromContext(r)

	if err != nil {
		utils.Respond(w, utils.Message(http.StatusBadRequest, "Resume ID is required"))
	}

	res := a.DAL.Resume().GetResumeByID(resumeID)

	utils.Respond(w, res)
}

func GetResumeIdFromContext(r *http.Request) (primitive.ObjectID, error) {
	id, ok := r.Context().Value("resume").(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, fmt.Errorf("Could Not Parse Object ID")
	}

	return id, nil
}
