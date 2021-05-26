package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

// Status godoc
// @Summary Check server health
// @Description get Ok
// @ID status
// @Produce  text/plain
// @Success 200 {string} Ok Ok
// @Router /health [get]
func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Ok")
}
