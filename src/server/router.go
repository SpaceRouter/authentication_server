package server

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/spacerouter/authentication_server/config"
	"github.com/spacerouter/authentication_server/controllers"
	_ "github.com/spacerouter/authentication_server/docs"
	"github.com/spacerouter/authentication_server/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	pam := controllers.PamController{
		Key:    config.GetSecretKey(),
		Issuer: "spacerouter",
	}

	router.POST("/login", pam.Authenticate)
	router.POST("/tea", controllers.GetTea)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("v1")
	{
		v1.Use(middleware.Auth(config.GetSecretKey()))
		v1.GET("/info", controllers.GetInfo)

		v1.GET("/roles", pam.GetUserRule)
		v1.POST("/update_password", pam.UpdatePassword)
	}
	return router

}
