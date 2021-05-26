package models

type Credential struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserInfo struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
