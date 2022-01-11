package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func IsAdmin(c *gin.Context) (isAdmin bool) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User_type key value could not be found. Cannot determine user previliges",
			})
			return false
		}

		if userType != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authorized.",
			})
			return false
		}

		return true
}