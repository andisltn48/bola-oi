package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Field struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
    Name     string             `json:"name,omitempty" validate:"required"`
    Location string             `json:"location,omitempty" validate:"required"`
    Type     string             `json:"type,omitempty" validate:"required"`
    Price    string             `json:"price,omitempty" validate:"required"`
    Unit     string             `json:"unit,omitempty" validate:"required"`
}