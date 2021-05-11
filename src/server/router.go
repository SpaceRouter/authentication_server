package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/controllers"
	"github.com/spacerouter/authentication_server/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	v1 := router.Group("v1")
	{
		pam := controllers.PamController{}
		v1.POST("/login", pam.Authenticate)

		admin := v1.Group("/admin")
		{
			admin.Use(middlewares.AuthMiddleware())
			admin.GET("/info", middlewares.GetInfo)
		}
	}
	return router

}
