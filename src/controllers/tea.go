package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/forms"
	"net/http"
)

func GetTea(c *gin.Context) {
	c.JSON(http.StatusTeapot, forms.UserLoginResponse{
		Message: "I'm a teapot",
		Ok:      true,
		Token:   "UWANNADRINKTEA",
	})
	return
}
