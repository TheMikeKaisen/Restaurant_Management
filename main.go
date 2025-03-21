package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/TheMikeKaisen/Restaurant_Management/database"
	"github.com/TheMikeKaisen/Restaurant_Management/routes"
	"github.com/TheMikeKaisen/Restaurant_Management/middleware"
	"go.mongodb.org/mongo-driver/mongo"

)

var foodCollection *mongo.Collection = database.OpenCollection(database.client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)


	router.Run(":"+port)

}
