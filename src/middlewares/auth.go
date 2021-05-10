package middlewares

import (
	"authentication_server/config"
	"authentication_server/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type CustomClaims struct {
	User models.User
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		tokenString = strings.Replace(tokenString, "bearer ", "", 1)

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			configs := config.GetConfig()
			return []byte(configs.GetString("security.secret_key")), nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
			return
		}

		c.Set("user", token.Claims.(*CustomClaims).User)
		c.Next()
	}
}

func CreateToken(user models.User) (string, error) {
	var claim = jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		Issuer:    "SpaceRouter",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{
		user,
		claim,
	})

	configs := config.GetConfig()

	return token.SignedString([]byte(configs.GetString("security.secret_key")))
}

func GetInfo(c *gin.Context) {
	info, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, info)

}
