package models

import (
	"net/http"
	"time"

	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExperienceDAL interface {
	GetExperience(id primitive.ObjectID) map[string]interface{}
}

type Experience struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Company         string             `json:"company"`
	JobTitle        string             `json:"jobTitle"`
	StartDate       time.Time          `json:"startDate"`
	EndDate         time.Time          `json:"endDate"`
	Accomplishments []string           `json:"accomplishments"`
	Current         bool               `json:"current"`
}

type ExperienceSource struct{}

func (e ExperienceSource) GetExperience(id primitive.ObjectID) map[string]interface{} {
	resume, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res := utils.Message(http.StatusOK, "Experience Found")
	res["data"] = resume.Experiences

	return res

}
