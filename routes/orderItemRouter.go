package routes

import (
	"github.com/TheMikeKaisen/Restaurant_Management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(orderItemRoutes *gin.Engine) {
	orderItemRoutes.GET("/orderItems", controllers.GetOrderItems())
	orderItemRoutes.GET("/orderItems/:orderItem_id", controllers.GetOrderItem())
	orderItemRoutes.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrder())
	orderItemRoutes.POST("/orderItems", controllers.CreateOrderItem())
	orderItemRoutes.PATCH("/orderItems/:orderItem_id", controllers.UpdateOrderItem())
}
