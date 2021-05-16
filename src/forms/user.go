package forms

import (
	"github.com/spacerouter/authentication_server/models"
)

type UserLogin struct {
	models.Credential
}

type UserLoginResponse struct {
	Message string
	Ok      bool
	Token   string
}

type UserCreate struct {
	models.User
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

type UserRolesResponse struct {
	Message string
	Ok      bool
	Roles   []models.Role
}

type UserInfoResponse struct {
	Message string
	Ok      bool
	User    models.User
}
