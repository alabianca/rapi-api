package models

import (
	"net/http"

	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PersonalDAL interface {
	GetPersonal(id primitive.ObjectID) map[string]interface{}
	GetSkills(id primitive.ObjectID) map[string]interface{}
}

type Personal struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Objective string             `json:"objective"`
}

type PersonalSource struct{}

func (ps PersonalSource) GetPersonal(id primitive.ObjectID) map[string]interface{} {
	resume, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res := utils.Message(http.StatusOK, "Personal Found")
	res["data"] = resume.Personal

	return res

}

func (ps PersonalSource) GetSkills(id primitive.ObjectID) map[string]interface{} {
	resume, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res := utils.Message(http.StatusOK, "Education Found")
	res["data"] = resume.Skills

	return res

}
