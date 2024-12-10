package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/polaris/codesandbox/settings"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	if !settings.RemoteConfig.JwtConfig.NeedAuth {
		return func(c *gin.Context) {}
	}
	return func(c *gin.Context) {
		abortWithAuthFailure := func() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Authorization failed",
			})
		}

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			abortWithAuthFailure()
		}
		user, err := ValidateToken(token)
		if err != nil {
			abortWithAuthFailure()
		}
		c.Set("user", user)
		c.Next()
	}
}
