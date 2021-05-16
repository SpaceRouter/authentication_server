package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/config"
	"github.com/spacerouter/authentication_server/controllers"
	"github.com/spacerouter/sr_auth"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	auth := sr_auth.Auth{
		Key: config.GetSecretKey(),
	}
	pam := controllers.PamController{
		Auth: auth,
	}

	router.POST("/login", pam.Authenticate)
	v1 := router.Group("v1")
	{

		v1.Use(auth.SrAuthMiddleware())
		v1.GET("/info", controllers.GetInfo)
		v1.POST("/update_password", pam.UpdatePassword)

	}
	return router

}
