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


func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)


		// Beginning of Pagination ( Will return to it later )
		// 
		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }

		// pages, err := strconv.Atoi(c.Query("pages"))
		// if err != nil || pages < 1 {
		// 	pages = 1
		// }

		// startIndex := (pages - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		results, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while searching the invoice collection",
				"error": err.Error(),
			})
			return
		}

		var allInvoices []bson.M
		
		if err := results.All(ctx, &allInvoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through all invoices",
				"error": err.Error(),
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		invoiceId := c.Param("invoice_id")
		var invoice models.Food

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while searching the invoice collection.",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, invoice)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var invoice models.Invoice
		var order models.Order

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error when attempting to bind request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": validationErr.Error(),
			})
			defer cancel()
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the order collection. Order with that order_id not found",
				"error": err.Error(),
			})
			return
		}

		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()
		invoice.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := invoiceCollection.InsertOne(ctx, invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the invoice collection.",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var reqInvoiceData models.Invoice

		if err := c.BindJSON(&reqInvoiceData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while attempting to bind request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}
		
		var updatedInvoice primitive.D
		invoiceID := c.Param("invoice_id")

		
		var menu models.Menu
		
		if reqInvoiceData.Order_id != "" {
			err := invoiceCollection.FindOne(ctx, bson.M{"order_id": reqInvoiceData.Order_id}).Decode(&menu)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "There was an error while querying the order collection. Please check the correct inputted Order ID",
					"error": err.Error(),
				})
				return
			}

			updatedInvoice = append(updatedInvoice, bson.E{Key: "order_id", Value: reqInvoiceData.Order_id})
		}

		if reqInvoiceData.Payment_Method != nil {
			updatedInvoice = append(updatedInvoice, bson.E{Key: "payment_method", Value: reqInvoiceData.Payment_Method})
		}

		if reqInvoiceData.Payment_status != nil {
			updatedInvoice = append(updatedInvoice, bson.E{Key: "payment_status", Value: reqInvoiceData.Payment_status})
		}

		if reqInvoiceData.Payment_due_date.String() != "" {
			updatedInvoice = append(updatedInvoice, bson.E{Key: "payment_due_date", Value: reqInvoiceData.Payment_due_date})
		}

		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedInvoice = append(updatedInvoice, bson.E{Key: "updated_at", Value: updatedAt})

		filter := bson.M{"invoice_id": invoiceID}
		update := bson.D{{Key: "$set", Value: updatedInvoice}}
		opt := options.Update().SetUpsert(true)

		results, err := invoiceCollection.UpdateOne(ctx, filter, update, opt)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the invoice collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Updated Successful",
			"updated_object": results,
		})

	}
}

func DeleteInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if admin := helpers.IsAdmin(c); !admin{
			defer cancel()
			return
		}
		
		invoiceId := c.Param("invoice_id")

		deleteCount, err := invoiceCollection.DeleteOne(ctx, bson.M{"invoice_id": invoiceId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while attempting to delete an object from collection invoice",
				"error": err.Error(),
				
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"number_of_objects_deleted": deleteCount.DeletedCount,
		})


	}
}