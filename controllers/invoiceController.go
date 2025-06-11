package controllers

import (
	"context"
	"net/http"
	"time"

	"golang-restrogo/database"
	"golang-restrogo/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

var validate = validator.New()

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error listing invoices"})
			return
		}

		var invoices []bson.M
		if err = cursor.All(ctx, &invoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing invoices"})
			return
		}

		c.JSON(http.StatusOK, invoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceID := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice not found"})
			return
		}

		// Mocking ItemsByOrder - replace this with actual logic if you have it
		allOrderItems := []bson.M{
			{
				"payment_due":  250.0,
				"table_number": "T01",
				"order_items":  []string{"Burger", "Fries"},
			},
		}

		view := InvoiceViewFormat{
			Invoice_id:       invoice.Invoice_id,
			Order_id:         invoice.Order_id,
			Payment_due_date: invoice.Payment_due_date,
			Payment_method:   getStringValue(invoice.Payment_method),
			Payment_status:   invoice.Payment_status,
			Payment_due:      allOrderItems[0]["payment_due"],
			Table_number:     allOrderItems[0]["table_number"],
			Order_details:    allOrderItems[0]["order_items"],
		}

		c.JSON(http.StatusOK, view)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate order existence
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order not found"})
			return
		}

		// Default values
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}
		now := time.Now()
		invoice.Payment_due_date = now.Add(24 * time.Hour)
		invoice.Created_at = now
		invoice.Updated_at = now
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		// Validation
		if err := validate.Struct(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D
		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
		}
		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
		}

		invoice.Updated_at = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

		upsert := true
		opts := options.UpdateOptions{Upsert: &upsert}
		filter := bson.M{"invoice_id": invoiceID}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opts,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return "null"
}
