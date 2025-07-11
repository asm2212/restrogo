package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Order_id   string             `json:"order_id"`
	Table_id   *string            `json:"table_id" validate:"required"`
	Order_Date time.Time          `json:"ordered_date"`
	Updated_at time.Time          `json:"updated_at"`
	Created_at time.Time          `json:"created_at"`
}
