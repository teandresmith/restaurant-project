package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/controllers"
)


func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/:foods_id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())
	incomingRoutes.DELETE("/foods/:foods_id", controllers.DeleteFood())
}