package routes

import ("github.com/gin-gonic/gin"
"github.com/TheMikeKaisen/Restaurant_Management/controllers")


func FoodRoutes(foodRoute *gin.Engine){
	foodRoute.GET("/food", controllers.GetFoods())
	foodRoute.GET("/food/:food_id", controllers.GetFood())
	foodRoute.POST("/foods", controllers.CreateFood())
	foodRoute.PATCH("foods/:food_id", controllers.UpdateFood())
	
}