package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/helpers"
	"github.com/teandresmith/restaurant-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		results, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while searching the invoice collection",
				"error": err.Error(),
			})
			return
		}

		var allOrders []bson.M
		
		if err := results.All(ctx, &allOrders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through all invoices",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allOrders)

	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while searching the invoice collection.",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if err := validate.Struct(order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		insertResults, insertErr := orderCollection.InsertOne(ctx, bson.M{"order_id": order.Order_id})
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the order collection",
				"error": insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Insertion Successful",
			"results": insertResults,
		})
	}
}


// Need to complete logic
func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func DeleteOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if admin := helpers.IsAdmin(c); !admin {
			defer cancel()
			return
		}

		orderId := c.Param("order_id")

		deleteResult, deleteErr := orderCollection.DeleteOne(ctx, bson.M{"order_id": orderId})
		defer cancel()
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object in the order collection",
				"error": deleteErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"results": deleteResult,
		})
	}
}