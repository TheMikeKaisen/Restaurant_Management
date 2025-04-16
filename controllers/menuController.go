package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/TheMikeKaisen/Restaurant_Management/database"
	"github.com/TheMikeKaisen/Restaurant_Management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var menus []bson.M
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": "cannot fetch menu data"})
			return
		}

		err = result.All(ctx, &menus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": "cannot decode the data"})
			return
		}
		c.JSON(200, menus)

	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		if menuId == "" {
			c.JSON(400, gin.H{"error": "menuId cannot be empty"})
			return
		}

		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": "cannot retrieve menu"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var menu models.Menu

		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(500, gin.H{"error: ": "bind failed"})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation error": validationErr})
			return
		}

		menu.ID = primitive.NewObjectID()
		menu.Created_at = time.Now().UTC()
		menu.Updated_at = time.Now().UTC()
		menu.Menu_id = menu.ID.Hex()

		result, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"insert error": err})
			return
		}

		c.JSON(200, result)

	}
}

func inTimeSpan(start, end time.Time) bool {
	// ensures that the start date always happens in the future and end date always lie after start date
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj bson.D

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.Updated_at = time.Now().UTC()
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			res, err := menuCollection.UpdateOne(ctx, filter, bson.D{
				{Key: "$set", Value: updateObj},
			})

			if err != nil {
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}

			c.JSON(http.StatusOK, res)

		}
	}
}
