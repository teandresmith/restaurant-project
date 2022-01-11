package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func DBInstance() *mongo.Client{
	// MongoDB := os.Getenv("MONGODB_URI")
	MongoDb := "mongodb+srv://teandre3:83760298@cluster0.ecbfw.mongodb.net/restaurant-project?retryWrites=true&w=majority"
	fmt.Print(MongoDb)
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("restuarant-project").Collection(collectionName)

	return collection
}