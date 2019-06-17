package models

type Registration struct {
	Email     string `json:email`
	Password  string `json:password`
	Verify    string `json:verify`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
