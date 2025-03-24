package project

import "go.mongodb.org/mongo-driver/mongo"

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(client *mongo.Client, dbName, collectionName string) *ProjectRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProjectRepository{collection: collection}
}
