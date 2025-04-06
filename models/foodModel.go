package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required,min=2,max=100"`
	Price      *float64           `json:"price" validate:"required"`
	Food_image *string            `json:"food_image" validate:"required"`
	Created_At time.Time          `json:"create_at"`
	Food_id    string             `json:"updated_at"`
	Menu_id    *string            `json:"menu_id" validate:"required"`
}
