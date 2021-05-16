package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}

func GetInfo(c *gin.Context) {
	info, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, info)
}
