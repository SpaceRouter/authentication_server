package models

type Credential struct {
	Login    string
	Password string
}

type Role string

type User struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
	Roles     []Role
}
