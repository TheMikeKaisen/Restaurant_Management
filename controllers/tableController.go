package controllers

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/TheMikeKaisen/Restaurant_Management/database"
)


var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")