package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.DELETE("users/:user_id", controllers.DeleteUser())
	incomingRoutes.PATCH("users/:user_id", controllers.EditUser())
}