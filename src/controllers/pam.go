package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/forms"
	"github.com/spacerouter/authentication_server/models"
	"github.com/spacerouter/authentication_server/utils"
	"github.com/spacerouter/pam"
	"net/http"
	"strings"
)

var BadRequestMessage = "Your request does not meet the expectations of the server"

type PamController struct {
	Key    string
	Issuer string
}

type Credential struct {
	models.Credential
}

func (c *Credential) RespondPAM(s pam.Style, msg string) (string, error) {
	switch s {
	case pam.PromptEchoOn:
		return c.Login, nil
	case pam.PromptEchoOff:
		return c.Password, nil
	}
	return "", errors.New("unexpected")
}

// Authenticate godoc
// @Summary Get authentication token
// @Description Get authentication token with login and password
// @ID authenticate
// @Param UserLogin body forms.UserLogin false "User credentials"
// @Accept  json
// @Produce  json
// @Success 200 {object} forms.UserLoginResponse
// @Failure 500,400,401 {object} forms.UserLoginResponse
// @Router /login [post]
func (p *PamController) Authenticate(c *gin.Context) {
	loginForm := forms.UserLogin{}
	err := c.BindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequestMessage)
		c.Abort()
		return
	}

	t, err := pam.Start("", "", &Credential{loginForm.Credential})
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserLoginResponse{
			Message: fmt.Sprintf("UserInfo doesn't exist \nError : %s", err),
			Ok:      false,
		})
		c.Abort()
		return
	}

	err = t.Authenticate(0)
	if err != nil {
		c.JSON(http.StatusUnauthorized, forms.UserLoginResponse{Message: fmt.Sprintf("Bad password or login \nError : %s", err), Ok: false, Token: ""})
		c.Abort()
		return
	}

	infos, err := pam.GetUserInfos(loginForm.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserLoginResponse{
			Message: fmt.Sprintf("Cannot get user infos \nError : %s", err),
			Ok:      false,
		})
		c.Abort()
		return
	}

	token, err := utils.CreateToken(infos.Username, p.Issuer, p.Key)

	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserLoginResponse{
			Message: fmt.Sprintf("Cannot create token \nError : %s", err),
			Ok:      false,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, forms.UserLoginResponse{Ok: true, Token: token})
}

// UpdatePassword godoc
// @Summary Change user password
// @Description Update user password using username and new password
// @ID update_password
// @Security ApiKeyAuth
// @Param UserLogin body forms.UserChangePassword false "User password change"
// @Accept  json
// @Produce  json
// @Security
// @Success 200 {object} forms.UserChangesResponse
// @Failure 500,400,401 {object} forms.UserChangesResponse
// @Router /v1/update_password [post]
func (p *PamController) UpdatePassword(c *gin.Context) {
	chgPwdForm := forms.UserChangePassword{}
	err := c.BindJSON(&chgPwdForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, forms.UserChangesResponse{
			Ok:      false,
			Message: BadRequestMessage,
		})
		c.Abort()
		return
	}

	uI, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, forms.UserChangesResponse{
			Ok:      false,
			Message: fmt.Sprintf("Cannot get session user \nError : %s", err),
		})
		c.Abort()
		return
	}

	username := uI.(string)

	err = pam.ChangePassword(username, chgPwdForm.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserChangesResponse{
			Ok:      false,
			Message: fmt.Sprintf("Cannot change password \nError : %s", err),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusUnauthorized, forms.UserChangesResponse{
		Ok:      true,
		Message: "",
	})
}

// GetInfo godoc
// @Summary Get user info
// @Description Get connected user information
// @ID get_info
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} forms.UserInfoResponse
// @Failure 500,400,401 {object} forms.UserInfoResponse
// @Router /v1/info [get]
func GetInfo(c *gin.Context) {
	info, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(500)
		return
	}

	username := info.(string)

	userInfo, err := pam.GetUserInfos(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserInfoResponse{
			Ok:      false,
			Message: fmt.Sprintf("Cannot change password \nError : %s", err),
		})
		c.Abort()
		return
	}

	u := UserInfoToUser(userInfo)

	c.JSON(200, forms.UserInfoResponse{
		Ok:       true,
		Message:  "Ok",
		UserInfo: u,
	})
}

// GetUserRole godoc
// @Summary Get user role
// @Description Get connected user role
// @ID get_role
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} forms.UserRolesResponse
// @Failure 500,400,401 {object} forms.UserRolesResponse
// @Router /v1/role [get]
func (p *PamController) GetUserRole(c *gin.Context) {
	uI, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, forms.UserRolesResponse{
			Ok:      false,
			Message: fmt.Sprintf("Cannot get session user \nError : UserInfo doesn't exist"),
		})
		c.Abort()
		return
	}

	username := uI.(string)
	role, err := utils.GetUserRole(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserRolesResponse{
			Message: fmt.Sprintf("Cannot get roles \nError : %s", err),
			Ok:      false,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, forms.UserRolesResponse{
		Role:    *role,
		Ok:      true,
		Message: "Ok",
	})
}

// GetUserRole2 godoc
// @Summary Get user role
// @Description Get connected user role
// @ID get_roles2
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} forms.UserRolesResponse
// @Failure 500,400,401 {object} forms.UserRolesResponse
// @Router /v1/user/{username}/role [get]
func (p *PamController) GetUserRole2(c *gin.Context) {
	p.GetUserRole(c)
}

// GetUserPermissions godoc
// @Summary Get user permissions
// @Description Get connected user permissions
// @ID get_permissions
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} forms.UserPermissionsResponse
// @Failure 500,400,401 {object} forms.UserPermissionsResponse
// @Router /v1/permissions [get]
func (p *PamController) GetUserPermissions(c *gin.Context) {

	username := c.Param("name")

	if username == "" {
		info, exist := c.Get("user")
		if !exist {
			c.AbortWithStatus(500)
			return
		}
		username = info.(string)
	}

	role, err := utils.GetUserRole(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserPermissionsResponse{
			Ok:      false,
			Message: fmt.Sprintf("Can't get user role \nError : %s", err),
		})
		c.Abort()
		return
	}

	permissions, err := utils.GetRolePermissions(*role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserPermissionsResponse{
			Ok:      false,
			Message: fmt.Sprintf("Can't get user permissions \nError : %s", err),
		})
		c.Abort()
		return
	}

	c.JSON(200, forms.UserPermissionsResponse{
		Ok:          true,
		Message:     "Ok",
		Role:        *role,
		Permissions: permissions,
	})
}

// GetUserPermissions2 godoc
// @Summary Get user permissions
// @Description Get connected user permissions
// @ID get_permissions2
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} forms.UserPermissionsResponse
// @Failure 500,400,401 {object} forms.UserPermissionsResponse
// @Router /v1/user/{username}/permissions [get]
func (p *PamController) GetUserPermissions2(c *gin.Context) {
	p.GetUserPermissions(c)
}

func UserInfoToUser(info *pam.UserInfo) models.UserInfo {
	userI := models.UserInfo{
		Login: info.Username,
	}

	uInfo := strings.Split(info.UserInformation, ",")
	names := strings.Split(uInfo[0], " ")
	userI.FirstName = names[0]
	userI.LastName = names[0]
	if len(uInfo) > 4 {
		userI.Email = uInfo[4]
	}
	return userI
}
