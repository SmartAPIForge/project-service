package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner  string             `bson:"owner" json:"owner"`
	Name   string             `bson:"name" json:"name"`
	Status string             `bson:"status" json:"status"`
	Data   string             `bson:"data" json:"data"`
}
