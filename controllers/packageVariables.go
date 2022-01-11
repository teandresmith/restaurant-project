package controllers

import (
	"github.com/go-playground/validator"
	"github.com/teandresmith/restaurant-project/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderitems")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
var validate = validator.New()

