package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spacerouter/authentication_server/utils"
	"net/http"
	"strings"
)

func Auth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		tokenString = strings.Replace(tokenString, "bearer ", "", 1)

		u, err := utils.GetUsernameFromToken(tokenString, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
			return
		}

		c.Set("user", u)
		c.Next()
	}
}
