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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		results, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the order item collection",
				"error": err.Error(),
			})
			return
		}

		var allOrderItems []bson.M

		if err := results.All(ctx, &allOrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating througuh all order items",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allOrderItems)

	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.Food
		orderItemId := c.Param("order_item_id")

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the order item collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, orderItem)


	}
}


func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		orderId := c.Param("order_id")

		results, err := orderItemCollection.Find(ctx, bson.M{"order_id": orderId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the order item collection",
				"error": err.Error(),
			})
		}

		var allOrderItems []bson.M

		if err := results.All(ctx, &allOrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through all the results",
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Query Successful",
			"results": allOrderItems,
		})
	}
}

// func ItemsByOrder(id string) (Order Items []primitives, err error){

// }

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if err := validate.Struct(orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		orderItem.ID = primitive.NewObjectID()
		orderItem.Order_item_id = orderItem.ID.Hex()
		orderItem.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		insertResults, insertErr := orderItemCollection.InsertOne(ctx, bson.M{"order_item_id": orderItem.Order_item_id})
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the order items collection",
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

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		var orderItemUpdate primitive.D
		orderItemId := c.Param("order_item_id")

		if orderItem.Quantity != nil {
			orderItemUpdate = append(orderItemUpdate, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}

		if orderItem.Unit_Price != nil {
			orderItemUpdate = append(orderItemUpdate, bson.E{Key: "unit_price", Value: orderItem.Unit_Price})
		}

		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemUpdate = append(orderItemUpdate, bson.E{Key: "updated_at", Value: updatedAt})

		filter := bson.M{"order_item_id": orderItemId}
		opt := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: orderItemUpdate}}

		updateResults, updateErr := orderItemCollection.UpdateOne(ctx, filter, update, opt)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the order item collection",
				"error": updateErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Update Successful",
			"results": updateResults,
		})

	}
}

func DeleteOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if admin := helpers.IsAdmin(c); !admin {
			defer cancel()
			return
		}

		orderItemId := c.Param("order_item_id")

		deleteResult, deleteErr := orderItemCollection.DeleteOne(ctx, bson.M{"order_item_id": orderItemId})
		defer cancel()
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object in orderitems collection",
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