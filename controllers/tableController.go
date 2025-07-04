package controllers

import (
	"context"
	"golang-restrogo/database"
	"golang-restrogo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while listing tables"})
			return
		}

		var tables []bson.M
		if err = cursor.All(ctx, &tables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the table"})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(table); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		now := time.Now()
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()
		table.Created_at = now
		table.Updated_at = now

		_, insertErr := tableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table could not be created"})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.D{}
		if table.Number != 0 {
			update = append(update, bson.E{Key: "number", Value: table.Number})
		}
		if table.Capacity != 0 {
			update = append(update, bson.E{Key: "capacity", Value: table.Capacity})
		}
		update = append(update, bson.E{Key: "updated_at", Value: time.Now()})

		filter := bson.M{"table_id": tableId}
		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table update failed"})
			return
		}

		if result.MatchedCount == 1 {
			var updatedTable models.Table
			err := tableCollection.FindOne(ctx, filter).Decode(&updatedTable)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch updated table"})
				return
			}
			c.JSON(http.StatusOK, updatedTable)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
	}
}
