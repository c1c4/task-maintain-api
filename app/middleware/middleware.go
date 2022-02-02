package middleware

import (
	"api/app/authentication"
	"api/app/utils/error_utils"
	"log"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("\n%s %s %s", c.Request.Method, c.Request.RequestURI, c.Request.Host)
		c.Next()
	}
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authentication.ValidateToken(c); err != nil {
			errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
			c.JSON(errUnauthorized.Status(), errUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
