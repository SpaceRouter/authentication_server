package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/forms"
	"github.com/spacerouter/authentication_server/models"
	"github.com/spacerouter/pam"
	"github.com/spacerouter/sr_auth"
	"net/http"
	"os/user"
	"strings"
)

var BadRequestMessage = "Your request does not meet the expectations of the server"

type PamController struct {
	Auth sr_auth.Auth
}

type Credential struct {
	models.Credential
}

func (c Credential) RespondPAM(s pam.Style, msg string) (string, error) {
	switch s {
	case pam.PromptEchoOn:
		return c.Login, nil
	case pam.PromptEchoOff:
		return c.Password, nil
	}
	return "", errors.New("unexpected")
}

func (p PamController) Authenticate(c *gin.Context) {
	loginForm := forms.UserLogin{}
	err := c.BindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequestMessage)
		c.Abort()
		return
	}

	t, err := pam.Start("", "", Credential{loginForm.Credential})
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserLoginResponse{
			Message: fmt.Sprintf("User doesn't exist \nError : %s", err),
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

	roles, err := GetUserRoles(loginForm.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, forms.UserLoginResponse{
			Message: fmt.Sprintf("Cannot get roles \nError : %s", err),
			Ok:      false,
		})
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

	userInfo := InfoToUser(*infos, roles)

	token, err := p.Auth.CreateToken(userInfo, "spacerouter")

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

func (p PamController) UpdatePassword(c *gin.Context) {
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

	u := uI.(sr_auth.User)
	if !u.HasRole(models.ChangeUserInfo) && chgPwdForm.User != u.Login {
		c.JSON(http.StatusUnauthorized, forms.UserChangesResponse{
			Ok:      false,
			Message: fmt.Sprintf("You are not authorized to modify this user \nError : %s", err),
		})
		c.Abort()
		return
	}

	err = pam.ChangePassword(u.Login, chgPwdForm.Password)
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

func InfoToUser(info pam.UserInfo, roles []sr_auth.Role) sr_auth.User {
	userI := sr_auth.User{
		Login: info.Username,
		Roles: roles,
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

func GetUserRoles(username string) ([]sr_auth.Role, error) {
	lookup, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	groupIds, err := lookup.GroupIds()
	if err != nil {
		return nil, fmt.Errorf("cannot reach group Ids : %s", err.Error())
	}

	var groups []sr_auth.Role
	for _, group := range groupIds {
		gr, err := user.LookupGroupId(group)
		if err != nil {
			return nil, fmt.Errorf("cannot convert id to name %s : %s", group, err.Error())
		}
		groups = append(groups, sr_auth.Role(gr.Name))
	}

	return groups, nil
}
