package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("x-admin-secret")

		if secret == "" || secret != os.Getenv("ADMIN_DASHBOARD_SECRET") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}