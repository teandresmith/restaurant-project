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

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		results, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the table collection",
				"error": err.Error(),
			})
			return
		}

		var allTables []bson.M

		if err := results.All(ctx, &allTables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through all table results",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allTables)

	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the table collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, table)

	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}


		if err := validate.Struct(table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There wasn an error while validating the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()
		table.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		

		results, insertErr := tableCollection.InsertOne(ctx, bson.M{"table_id": table.Table_id})
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the table collection",
				"error": insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Insertion Successful",
			"results": results,
		})
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		var tableUpdates primitive.D
		tableId := c.Param("table_id")


		if table.Number_of_guests != nil {
			tableUpdates = append(tableUpdates, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}

		if table.Table_number != nil {
			tableUpdates = append(tableUpdates, bson.E{Key: "table_number", Value: table.Table_number})
		}

		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tableUpdates = append(tableUpdates, bson.E{Key: "updated_at", Value: updatedAt})

		filter := bson.M{"table_id": tableId}
		opt := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: tableUpdates}}

		
		updateResults, updateErr := tableCollection.UpdateOne(ctx, filter, update, opt)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the table collection",
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

func DeleteTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if admin := helpers.IsAdmin(c); !admin {
			defer cancel()
			return
		}

		tableId := c.Param("table_id")

		deleteResults, err := tableCollection.DeleteOne(ctx, bson.M{"table_id": tableId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an issue while deleting an object in the table collection",
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"results": deleteResults,
		})
	}
}