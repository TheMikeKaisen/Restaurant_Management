package routes

import (
	"github.com/TheMikeKaisen/Restaurant_Management/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(foodRoute *gin.Engine) {
	foodRoute.GET("/food", controllers.GetFoods())
	foodRoute.GET("/food/:food_id", controllers.GetFood())
	foodRoute.POST("/foods", controllers.CreateFood())
	foodRoute.PATCH("foods/:food_id", controllers.UpdateFood())

}
