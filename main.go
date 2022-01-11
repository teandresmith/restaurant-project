package main

import (
	"log"
	"os"

	"github.com/teandresmith/restaurant-project/middleware"
	"github.com/teandresmith/restaurant-project/routes"

	"github.com/gin-gonic/gin"
) 



func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()

	// Default Middleware
	router.Use(gin.Logger())
	routes.LoginRoutes(router)

	// Custom Authentication Middleware
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)

	log.Fatal(router.Run(":" + port))
	
}