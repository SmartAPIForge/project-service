package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ComposeId string             `bson:"composeId" json:"composeId"`
	Owner     string             `bson:"owner" json:"owner"`
	Name      string             `bson:"name" json:"name"`
	Data      string             `bson:"data" json:"data"`
	Status    string             `bson:"status" json:"status"`
	UrlZip    string             `bson:"urlZip" json:"urlZip"`
	UrlDeploy string             `bson:"urlDeploy" json:"urlDeploy"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
