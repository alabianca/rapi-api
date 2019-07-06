package models

import (
	"net/http"
	"time"

	"github.com/alabianca/rapi-api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EducationDAL interface {
	GetEducation(id primitive.ObjectID) map[string]interface{}
}

type Education struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	School    string             `json:"school"`
	GPA       float32            `json:"gpa"`
	StartDate time.Time          `json:"startDate"`
	EndDate   time.Time          `json:"endDate"`
	Degree    string             `json:"degree"`
}

func GetEducation(id primitive.ObjectID) map[string]interface{} {
	resume, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res := utils.Message(http.StatusOK, "Education Found")
	res["data"] = resume.Education

	return res

}
