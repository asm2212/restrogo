package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-restrogo/models"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", nil},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}},
		}}}
		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"menu_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}

		result, err := menuCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing menu items"})
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		objId, _ := primitive.ObjectIDFromHex(menuId)

		err := menuCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&menu)
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{"category", menu.Category})
		}
		if menu.Start_Date != nil {
			updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
		}
		if menu.End_Date != nil {
			updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		objId, _ := primitive.ObjectIDFromHex(menuId)
		filter := bson.M{"_id": objId}
		update := bson.D{{"$set", updateObj}}

		result, err := menuCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu update failed"})
			return
		}

		// Return the updated document
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
