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
	"golang.org/x/crypto/bcrypt"
)



func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if admin := helpers.IsAdmin(c); !admin {
			defer cancel()
			return
		}

		results, err := userCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the user collection",
				"error": err.Error(),
			})
			return
		}
		

		var allUsers []bson.M

		if err := results.All(ctx, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error during the iteration of all users",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allUsers)

	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		userId := c.Param("user_id")

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user); 
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}



func EditUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)


		userType, err := c.Get("uid")
		if !err {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User_type key value was not sent. Cannot determined privileges.",
			})
			defer cancel()
			return
		}

		userID := c.Param("user_id")

		if userType != userID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authorized.",
			})
			defer cancel()
			return
		}

		var userEdits models.User

		if bindErr := c.BindJSON(&userEdits); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": bindErr.Error(),
			})
			defer cancel()
			return
		}

		var userUpdateObject primitive.D

		if userEdits.User_Type != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authorized to change roles",
			})
			defer cancel()
			return
		}

		if userEdits.First_Name != nil {
			userUpdateObject = append(userUpdateObject, bson.E{Key: "first_name", Value: userEdits.First_Name})
		}

		if userEdits.Last_Name != nil {
			userUpdateObject = append(userUpdateObject, bson.E{Key: "last_name", Value: userEdits.Last_Name})
		}

		if userEdits.Password != nil {
			newPassword, hashErr := HashPassword(*userEdits.Password)
			if hashErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "There was an issue while hashing the new password",
					"error": hashErr.Error(),
				})
				defer cancel()
				return
			}
			userUpdateObject = append(userUpdateObject, bson.E{Key: "Password", Value: newPassword})
		}

		if userEdits.Email != nil {
			userUpdateObject = append(userUpdateObject, bson.E{Key: "email", Value: userEdits.Email})
		}

		if userEdits.Avatar != nil {
			userUpdateObject = append(userUpdateObject, bson.E{Key: "avatar", Value: userEdits.Avatar})
		}

		
		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userUpdateObject = append(userUpdateObject, bson.E{Key: "updated_at", Value: updatedAt})


		filter := bson.M{"user_id": userID}
		opt := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: userUpdateObject}}

		results, updateErr := userCollection.UpdateOne(ctx, filter, update, opt)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the user collection",
				"error": updateErr.Error(),
			})
		}

		c.JSON(http.StatusOK, results)
		
	}
}

func DeleteUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if IsAdmin := helpers.IsAdmin(c); !IsAdmin {
			defer cancel()
			return
		}

		userID := c.Param("user_id")

		result, deleteErr := userCollection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userID}})
		defer cancel()
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object in the user collection",
				"error": deleteErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"result": result,
		})
	}
}

func HashPassword(password string) (string, error) {
	passwordInBytes := []byte(password)
	salt := 10
	hashedPassword, hashErr := bcrypt.GenerateFromPassword(passwordInBytes, salt )
	if hashErr != nil {
		return ("Error during Hashing of password"), hashErr
	}
	return string(hashedPassword), nil
}

func VerifyPassword(userPassword string, providerPassword string)(error){
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providerPassword))
	if err != nil {
		return err 
	}
	return nil
}


