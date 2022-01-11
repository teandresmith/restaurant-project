package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/teandresmith/restaurant-project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SECRET_KEY  = []byte(os.Getenv("SECRET_KEY")) 
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")


type SignedTokenDetails struct {
	First_Name	string
	Last_Name	string
	User_Type	string
	Email		string
	Uid			string
	*jwt.RegisteredClaims
}


func GenerateTokens(firstName string, lastName, userType string, email string, uid string) (signedToken string, signedRefreshToken string, err error) {
	

	claims := &SignedTokenDetails{
		First_Name: firstName,
		Last_Name: 	lastName,
		User_Type: 	userType,
		Email: 		email,
		Uid: 		uid,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Second*10)),
			Issuer: "Restaurant-Management-Backend-Project",
		},
		
	}

	refreshTokenClaims := &SignedTokenDetails{
		First_Name: firstName,
		Last_Name: 	lastName,
		User_Type: 	userType,
		Email: 		email,
		Uid: 		uid,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour*24)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}


	return token, refreshToken, nil
	
}

func UpdateTokens(token string, refreshToken string, userId string) (error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateUserToken primitive.D

	updateUserToken = append(updateUserToken, bson.E{Key: "token", Value: token})
	updateUserToken = append(updateUserToken, bson.E{Key: "refresh_token", Value: refreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateUserToken = append(updateUserToken, bson.E{Key: "updated_at", Value: updated_at})

	filter := bson.M{"user_id": userId}
	opt := options.Update().SetUpsert(true)

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateUserToken}}, opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func VerifyToken(tokenString string) (claims *SignedTokenDetails, message string) {


	token, err := jwt.ParseWithClaims(tokenString, &SignedTokenDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	

	claims, ok := token.Claims.(*SignedTokenDetails) 
	if !ok {
		message = err.Error()
		return
	}

	//the token is expired
	if claims.VerifyExpiresAt(time.Now().Local(), true) {
		message = err.Error()
		return
	}

	return claims, message

}