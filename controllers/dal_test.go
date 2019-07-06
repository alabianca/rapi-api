package controllers

import (
	"github.com/alabianca/rapi-api/models"
)

type TestDAL struct{}

func getTestAPI() *API {
	dal := TestDAL{}
	api := &API{
		DAL: dal,
	}

	return api
}

func (d TestDAL) Keys() models.APIKeyDAL {
	return models.APIKeySource{}
}

func (d TestDAL) Education() models.EducationDAL {
	return models.EducationSource{}
}

func (d TestDAL) Experience() models.ExperienceDAL {
	return models.ExperienceSource{}
}

func (d TestDAL) Logs() models.LogDAL {
	return models.LogSource{}
}

func (d TestDAL) Personal() models.PersonalDAL {
	return models.PersonalSource{}
}

func (d TestDAL) Projects() models.ProjectDAL {
	return models.ProjectSource{}
}

func (d TestDAL) Resume() models.ResumeDAL {
	return models.ResumeSource{}
}

func (d TestDAL) Tokens() models.TokenDAL {
	return models.TokenSource{}
}

func (d TestDAL) Users() models.UserDAL {
	return TestUsersDAL{}
}
