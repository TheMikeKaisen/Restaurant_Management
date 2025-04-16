package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/TheMikeKaisen/Restaurant_Management/database"
	"github.com/TheMikeKaisen/Restaurant_Management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// create/refer food collection
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate *validator.Validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		if val := c.Query("startIndex"); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil && parsed >= 0 {
				startIndex = parsed
			}
		}
		matchStage := bson.D{
			{Key: "$match", Value: bson.D{}},
		}
		
		facetStage := bson.D{
			{Key: "$facet", Value: bson.D{
				{Key: "metadata", Value: bson.A{
					bson.D{{Key: "$count", Value: "total_count"}},
				}},
				{Key: "data", Value: bson.A{
					bson.D{{Key: "$skip", Value: startIndex}},
					bson.D{{Key: "$limit", Value: recordPerPage}},
				}},
			}},
		}
		
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "total_count", Value: bson.D{
					{Key: "$arrayElemAt", Value: bson.A{"$metadata.total_count", 0}},
				}},
				{Key: "food_items", Value: "$data"},
			}},
		}

		cursor, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, facetStage, projectStage})
		if err != nil {
			log.Printf("Aggregation error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing food items"})
			return
		}

		var resultData []bson.M
		if err := cursor.All(ctx, &resultData); err != nil {
			log.Fatal(err)
		}

		if len(resultData) == 0 {
			c.JSON(http.StatusOK, gin.H{"total_count": 0, "food_items": []interface{}{}})
			return
		}
		c.JSON(http.StatusOK, resultData[0])
	}
}
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodId := c.Param("food_id")
		if foodId == "" {
			c.JSON(400, gin.H{"error": "foodId cannot be empty"})
			return
		}

		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(500, gin.H{"error": "cannot fetch food at the moment"})
			return
		}
		c.JSON(200, food)

	}
}
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
			return
		}

		err := validate.Struct(food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": "validation error"})
			return
		}

		var menu models.Menu
		err = menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu not found!"})
			return
		}

		food.ID = primitive.NewObjectID()
		food.Created_at = time.Now().UTC()
		food.Updated_at = time.Now().UTC()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, err := foodCollection.InsertOne(ctx, food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": "could not create food item"})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")
		filter := bson.M{"food_id":foodId}

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj bson.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value:  food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value:  food.Price})
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value:  food.Food_image})
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			
			if err != nil {
				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu", Value:  food.Menu_id})
		}

		food.Updated_at = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value:  food.Updated_at})

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
		)

		if err != nil {
			msg := fmt.Sprint("food item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
