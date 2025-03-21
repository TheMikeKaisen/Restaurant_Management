package routes

import ("github.com/gin-gonic/gin"
"github.com/TheMikeKaisen/Restaurant_Management/controllers")


func UserRoutes(userRoute *gin.Engine){
	userRoute.GET("/users", controllers.GetUsers())
	userRoute.GET("/user/:user_id", controllers.GetUser())
	userRoute.POST("/user/signup", controllers.SignUp())
	userRoute.POST("/user/login", controllers.Login())
	
}