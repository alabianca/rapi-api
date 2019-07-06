package models

type DAL interface {
	Keys() APIKeyDAL
	Education() EducationDAL
	Experience() ExperienceDAL
	Logs() LogDAL
	Personal() PersonalDAL
	Projects() ProjectDAL
	Resume() ResumeDAL
	Tokens() TokenDAL
	Users() UserDAL
}
