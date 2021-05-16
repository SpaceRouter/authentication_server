package forms

import "github.com/spacerouter/authentication_server/models"

type UserLogin struct {
	models.Credential
}

type UserLoginResponse struct {
	Message string
	Ok      bool
	Token   string
}

type UserCreate struct {
	User     string
	Password string
	Email    string
	Roles    []string
}

type UserChangeRole struct {
	User string
	Role string
}

type UserChangePassword struct {
	User     string
	Password string
}

type UserChangesResponse struct {
	Message string
	Ok      bool
}
