package controllers

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		// Attempt at Pagination and Mongoose Aggregrated Queries
		// Will continue in the future 

		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }

		// page, err := strconv.Atoi(c.Query("page"))
		// if err != nil || page < 1{
		// 	page = 1
		// }

		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// matchStage := bson.D{{"$match", bson.D{{}}}}
		// groupStage := bson.D{
		// 	{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}

		// projectStage := bson.D{
		// 	{
		// 		"$project", bson.D{
		// 			{"_id", 0},
		// 			{"total_count", 1},
		// 			{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
						
		// 		}
		// 	}
		// }

		// aggregatedResults, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
		// 	matchStage, groupStage, projectStage
		// })
		// Mongoose Aggregrate Search Results
		// Match Stage
		// Group Stage
		// Project State


		results, err := foodCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying food collection",
				"error": err.Error(),
			})
			return
		}

		var allFoods []bson.M

		if err := results.All(ctx, &allFoods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error during iteration of all foods",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFood() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the food collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var menu models.Menu
		var food models.Food

		
		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while trying to read request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}
		

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The request body data could not be validated",
				"error": validationErr.Error(),
			})
			defer cancel()
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_ID}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the menu collection.",
				"error": err.Error(),
			})
			return
		}

		food.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_ID = food.ID.Hex()
		var price = toFixed(*food.Price, 2)
		food.Price = &price

		result, insertErr := foodCollection.InsertOne(ctx, food)
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while trying to insert object into food collection",
				"error": insertErr.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		
		var food models.Food
		foodId := c.Param("food_id")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request body data was unable to processed",
				"error": err.Error(),
			})
			defer cancel()
			return
		}
		
		var updatedFood primitive.D
		var menu models.Menu
		

		if food.Name != nil {
			updatedFood = append(updatedFood, bson.E{Key: "name",Value: food.Name})
		}

		if food.Price != nil {
			updatedFood = append(updatedFood, bson.E{Key: "price", Value: toFixed(*food.Price, 2)})
		}

		if food.Image != nil {
			updatedFood = append(updatedFood, bson.E{Key: "image", Value:  food.Image})
		}

		if food.Menu_ID != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_ID}).Decode(&menu)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "There was an error while querying the menu collection. Menu ID not found",
					"error": err.Error(),
				})
			}
			updatedFood = append(updatedFood, bson.E{Key: "menu_id", Value: food.Menu_ID})
		}

		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedFood = append(updatedFood, bson.E{Key: "updated_at", Value: food.Updated_At})

		
		filter := bson.M{"food_id": foodId}
		opt := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: updatedFood}}

		results, err := foodCollection.UpdateOne(ctx, filter, update, opt )
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating object in food collection",
				"error": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, results)
	}
}

func DeleteFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		foodId := c.Param("food_id")
		
		res, err := foodCollection.DeleteOne(ctx, bson.D{{Key: "food_id", Value: foodId}})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while attempting to delete an object from the food collection",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"number_of_objects_deleted": res.DeletedCount,
		})
	}
}



func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64{
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}