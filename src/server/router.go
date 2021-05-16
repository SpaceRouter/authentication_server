package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/config"
	"github.com/spacerouter/authentication_server/controllers"
	"github.com/spacerouter/authentication_server/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	pam := controllers.PamController{
		Key:    config.GetSecretKey(),
		Issuer: "spacerouter",
	}

	router.POST("/login", pam.Authenticate)
	router.POST("/tea", controllers.GetTea)

	v1 := router.Group("v1")
	{
		v1.Use(middleware.Auth(config.GetSecretKey()))
		v1.GET("/info", controllers.GetInfo)

		v1.GET("/roles", pam.GetUserRule)
		v1.POST("/update_password", pam.UpdatePassword)
	}
	return router

}
