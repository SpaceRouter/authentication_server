package forms

import (
	"github.com/spacerouter/authentication_server/models"
)

type UserLogin struct {
	models.Credential
}

type UserLoginResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Token   string `json:"token"`
}

type UserCreate struct {
	models.UserInfo
}

type UserChangeRole struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type UserChangePassword struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type UserChangesResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type UserRolesResponse struct {
	Message string        `json:"message"`
	Ok      bool          `json:"ok"`
	Roles   []models.Role `json:"roles"`
}

type UserInfoResponse struct {
	Message  string
	Ok       bool
	UserInfo models.UserInfo
}
