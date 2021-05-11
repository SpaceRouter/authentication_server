package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/middlewares"
	"github.com/spacerouter/authentication_server/models"
	"github.com/spacerouter/pam"
	"net/http"
	"os/user"
)

var BadRequestMessage = gin.H{"message": "your request does not meet the expectations of the server"}

type PamController struct {
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
	cred := models.Credential{}
	err := c.BindJSON(&cred)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequestMessage)
		c.Abort()
		return
	}

	t, err := pam.Start("", "", Credential{cred})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
		c.Abort()
		return
	}

	err = t.Authenticate(0)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Bad password or login", "ok": false, "token": ""})
		c.Abort()
		return
	}
	userInfo, err := user.Lookup(cred.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
		c.Abort()
		return
	}
	groupIds, err := userInfo.GroupIds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
		c.Abort()
		return
	}

	var groups []string
	for _, group := range groupIds {
		gr, err := user.LookupGroupId(group)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
			c.Abort()
			return
		}
		groups = append(groups, gr.Name)
	}
	infos, err := pam.GetUserInfos(cred.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
		c.Abort()
		return
	}

	token, err := middlewares.CreateToken(models.User{
		Login:     cred.Login,
		Roles:     groups,
		FirstName: infos.UserInformation,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to create token", "error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "", "ok": true, "token": token})
}
