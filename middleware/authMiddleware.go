package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/helpers"
)


func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "No Authorization header provided",
			})
			return
		}

		claims, message :=  helpers.VerifyToken(token); 
		if message != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token not valid",
				"error": message,
			})
			return
		}

		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)
		c.Set("user_type", claims.User_Type)
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
			
	}
}
