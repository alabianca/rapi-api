package controllers

import (
	"github.com/alabianca/rapi-api/models"
)

type API struct {
	DAL models.DAL
}

type DefaultDAL struct{}

func (d DefaultDAL) Keys() models.APIKeyDAL {
	return models.APIKeySource{}
}

func (d DefaultDAL) Education() models.EducationDAL {
	return models.EducationSource{}
}

func (d DefaultDAL) Experience() models.ExperienceDAL {
	return models.ExperienceSource{}
}

func (d DefaultDAL) Logs() models.LogDAL {
	return models.LogSource{}
}

func (d DefaultDAL) Personal() models.PersonalDAL {
	return models.PersonalSource{}
}

func (d DefaultDAL) Projects() models.ProjectDAL {
	return models.ProjectSource{}
}

func (d DefaultDAL) Resume() models.ResumeDAL {
	return models.ResumeSource{}
}

func (d DefaultDAL) Tokens() models.TokenDAL {
	return models.TokenSource{}
}

func (d DefaultDAL) Users() models.UserDAL {
	return models.UserSource{}
}
