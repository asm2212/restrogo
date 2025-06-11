package controllers

import (
	"context"
	"fmt"
	"golang-restrogo/database"
	"golang-restrogo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while listing orders"})
			return
		}

		var orders []bson.M
		if err = cursor.All(ctx, &orders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the order"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(order); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

	    if order.Table_id!= nil {
			err := tableCollection.FindOne(ctx,bson.M{
				"table_id" : order.Table_id
			}).Decode(&table)
			defer cancel()
			if err!=nil{
				msg:=fmt.Sprintf("message:Table was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
		}

		now := time.Now()
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		order.Created_at = now
		order.Updated_at = now
		order.Ordered_at = now
		order.Price = toFixed(totalPrice, 2)

		_, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order could not be created"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.D{}
		if order.Table_id != "" {
			update = append(update, bson.E{Key: "table_id", Value: order.Table_id})
		}
		if len(order.Items) > 0 {
			var totalPrice float64
			for _, foodID := range order.Items {
				var food models.Food
				err := foodCollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
				if err != nil {
					msg := fmt.Sprintf("food item %s not found", foodID)
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
				totalPrice += *food.Price
			}
			update = append(update, bson.E{Key: "items", Value: order.Items})
			update = append(update, bson.E{Key: "price", Value: toFixed(totalPrice, 2)})
		}

		update = append(update, bson.E{Key: "updated_at", Value: time.Now()})

		filter := bson.M{"order_id": orderId}
		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order update failed"})
			return
		}

		if result.MatchedCount == 1 {
			var updatedOrder models.Order
			err := orderCollection.FindOne(ctx, filter).Decode(&updatedOrder)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch updated order"})
				return
			}
			c.JSON(http.StatusOK, updatedOrder)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	}
}


func OrderItemOrderCreator(order models.Order) string {
	order.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	order.UPdated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx,order)
	defer cancel

	return order.Order_id
}