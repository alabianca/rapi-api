package models

import (
	"net/http"

	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectDAL interface {
	GetProjects(id primitive.ObjectID) map[string]interface{}
}

type Project struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string             `json:"string"`
	Link   string             `json:"link"`
	Detail string             `json:"detail"`
}

type ProjectSource struct{}

func (ps ProjectSource) GetProjects(id primitive.ObjectID) map[string]interface{} {
	resume, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res := utils.Message(http.StatusOK, "Projects Found")
	res["data"] = resume.Projects

	return res

}
