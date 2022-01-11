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


func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User


		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error when validating the request body data",
				"error": validationErr.Error(),
			})
		}

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, tokenErr := helpers.GenerateTokens(*user.First_Name, *user.Last_Name, *user.User_Type, *user.Email, user.User_id)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while creating a new token",
				"error": tokenErr.Error(),
			})
		}
		user.Token = &token
		user.Refresh_Token = &refreshToken

		hashPassword, hashErr := HashPassword(*user.Password)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while hashinig the password",
				"error": hashErr.Error(),
			})
		}
		user.Password = &hashPassword

		result, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the user collection",
				"error": insertErr.Error(),
			})
		}
		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var userInput models.User

		if err := c.BindJSON(&userInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return 
		}
		

		err := userCollection.FindOne(ctx, bson.M{"email": userInput.Email}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the user collection",
				"error": err.Error(),
			})
			return
		}

		if userInput.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Password provided"})
			return
		} else {
			if err := VerifyPassword(*user.Password,*userInput.Password); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
		}

		// Create new tokens everytime user login
		newToken, newRefreshToken, err := helpers.GenerateTokens(*user.First_Name, *user.Last_Name, *user.User_Type, *user.Email, user.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while generating new tokens",
				"error": err.Error(),
			})
		}

		updateErr := helpers.UpdateTokens(newToken, newRefreshToken, user.User_id)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating user tokens",
				"error": updateErr.Error(),
			})
		}
		
		
		c.JSON(http.StatusOK, user)
	}
}