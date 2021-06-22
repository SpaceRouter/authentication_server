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

	main := router.Group("auth")
	{

		health := new(controllers.HealthController)

		main.GET("/health", health.Status)

		pam := controllers.PamController{
			Key:    config.GetSecretKey(),
			Issuer: "spacerouter",
		}

		main.POST("/login", pam.Authenticate)
		main.POST("/tea", controllers.GetTea)
		main.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		v1 := main.Group("v1")
		{
			v1.Use(middleware.Auth(config.GetSecretKey()))
			v1.GET("/info", controllers.GetInfo)

			v1.GET("/role", pam.GetUserRole)
			v1.GET("/permissions", pam.GetUserPermissions)

			v1.POST("/update_password", pam.UpdatePassword)

			user := v1.Group("user")
			{
				user.GET(":name/permissions", pam.GetUserPermissions)
				user.GET(":name/role", pam.GetUserRole)
			}
		}
	}
	return router

}
