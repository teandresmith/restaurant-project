package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		// Beginning of Pagination Logic (will return to it later)

		// recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))
		// if err != nil || recordsPerPage < 1 {
		// 	recordsPerPage = 10
		// }

		// pages, err := strconv.Atoi(c.Query("pages"))
		// if err != nil || pages < 1 {
		// 	pages = 1
		// }

		// startingIndex, err := strconv.Atoi(c.Query("startingIndex"))
		// if err != nil || startingIndex < 1 {
		// 	startingIndex = (pages - 1) * recordsPerPage
		// }

		results, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu items"})
			return
		}

		var allMenus []bson.M
		if err = results.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		menu_Id := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menu_Id}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		menu.Created_At,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_Id = menu.ID.Hex()
		
		result, insertErr := menuCollection.InsertOne(ctx, menu)
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var reqMenuData models.Menu

		if err := c.BindJSON(&reqMenuData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		var updatedMenu primitive.D
		menu_id := c.Param("menu_id")

		if reqMenuData.Name != "" {
			updatedMenu = append(updatedMenu, bson.E{Key: "name", Value: reqMenuData.Name})
		}

		if reqMenuData.Category != "" {
			updatedMenu = append(updatedMenu, bson.E{Key: "category", Value: reqMenuData.Category})
		}

		if reqMenuData.Start_Date != nil {
			updatedMenu = append(updatedMenu, bson.E{Key: "start_date", Value: reqMenuData.Start_Date})
		}

		if reqMenuData.End_Date != nil {
			updatedMenu = append(updatedMenu, bson.E{Key: "end_date", Value: reqMenuData.End_Date})
		}

		updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedMenu = append(updatedMenu, bson.E{Key: "updated_at", Value: updated_at})

		filter := bson.M{"menu_id": menu_id}
		opt := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: updatedMenu}}

		results, err := menuCollection.UpdateOne(ctx, filter, update, opt)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the menu collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Update Successful",
			"updated_object": results,
		})

	}
}

func DeleteMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		menuId := c.Param("menu_id")

		deleteResult, err := menuCollection.DeleteOne(ctx, bson.M{"menu_id": menuId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object in the menu collection",
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"number_of_deleted_object": deleteResult.DeletedCount,
		})
	}
}