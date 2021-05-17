package models

type Credential struct {
	Login    string
	Password string
}

type UserInfo struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
}
