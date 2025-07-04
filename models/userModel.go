package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	User_id      string             `json:"user_id"`
	First_name   string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name    string             `json:"last_name" validate:"required,min=2,max=100"`
	Email        string             `json:"email" validate:"email,required"`
	Password     string             `json:"password,omitempty" validate:"required,min=6"`
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refresh_token"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}
