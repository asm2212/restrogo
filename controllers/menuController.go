package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang-restrogo/database"
	"golang-restrogo/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Add pagination parameters
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		// MongoDB query with pagination
		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "menu_items", Value: bson.D{
					{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
				}},
			}},
		}

		cursor, err := menuCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing menu items"})
			return
		}

		var allMenus []bson.M
		if err = cursor.All(ctx, &allMenus); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(allMenus) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No menus found", "data": []interface{}{}})
			return
		}

		c.JSON(http.StatusOK, allMenus[0])
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		// menuId in the database is likely stored as a string (menu_id field), not ObjectId.
		// We will search by menu_id (string hex) instead of _id (ObjectId).
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while fetching the menu",
			})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		now := time.Now()
		menu.Created_at = &now
		menu.Updated_at = menu.Created_at
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj bson.D

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}
		if menu.Start_Date != nil {
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
		}
		if menu.End_Date != nil {
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})
		}

		now := time.Now()
		menu.Updated_at = &now
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		// Query by menu_id string field, NOT _id (ObjectId)
		filter := bson.M{"menu_id": menuId}
		update := bson.D{{Key: "$set", Value: updateObj}}

		result, err := menuCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu update failed"})
			return
		}

		// Return the updated document only if matched
		var updatedMenu models.Menu
		if result.MatchedCount == 1 {
			err := menuCollection.FindOne(ctx, filter).Decode(&updatedMenu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "menu update failed"})
				return
			}
		}

		c.JSON(http.StatusOK, updatedMenu)
	}
}
