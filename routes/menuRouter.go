package routes

import ("github.com/gin-gonic/gin"
"github.com/TheMikeKaisen/Restaurant_Management/controllers")



func MenuRoutes(menuRoute *gin.Engine){
	menuRoute.GET("/menus", controllers.GetMenus())
	menuRoute.GET("/menu/:menu_id", controllers.GetMenu())
	menuRoute.POST("/menu", controllers.CreateMenu())
	menuRoute.PATCH("menu/:menu_id", controllers.UpdateMenu())
}

