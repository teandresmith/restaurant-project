package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/controllers"
)

func LoginRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.SignUp())
	incomingRoutes.POST("/login", controllers.Login())
} 