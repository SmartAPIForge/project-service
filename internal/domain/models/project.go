package models

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner  string             `bson:"owner" json:"owner"`
	Name   string             `bson:"name" json:"name"`
	Status string             `bson:"status" json:"status"`
	Data   json.RawMessage    `bson:"data" json:"data"`
}
