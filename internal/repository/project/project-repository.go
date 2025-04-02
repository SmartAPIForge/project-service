package project

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project-service/internal/domain/models"
	"time"
)

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(client *mongo.Client, dbName, collectionName string) *ProjectRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProjectRepository{collection: collection}
}

func (r *ProjectRepository) GetAllUserProjects(ctx context.Context, owner string, page, limit int64) ([]*models.Project, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.M{"composeId": 1})

	cursor, err := r.collection.Find(ctx, bson.M{"owner": owner}, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	var projects []*models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) InitProject(ctx context.Context, composeId, owner, name string) (*models.Project, error) {
	existingProject, err := r.getProjectByComposeId(ctx, composeId)
	if err != nil {
		return nil, err
	}
	if existingProject != nil {
		return nil, errors.New("проект с таким названием уже существует для данного пользователя")
	}

	project := &models.Project{
		ComposeId: composeId,
		Owner:     owner,
		Name:      name,
		Data:      "",
		Status:    "NEW",
		UrlZip:    "",
		UrlDeploy: "",
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = r.collection.InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, composeId string, data string) (*models.Project, error) {
	existingProject, err := r.getProjectByComposeId(ctx, composeId)
	if err != nil {
		return nil, err
	}
	if existingProject == nil {
		return nil, errors.New("проект не найден")
	}

	update := bson.M{"$set": bson.M{
		"data":      data,
		"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
	}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"composeId": composeId}, update)
	if err != nil {
		return nil, err
	}

	return r.getProjectByComposeId(ctx, composeId)
}

func (r *ProjectRepository) UpdateProjectStatus(ctx context.Context, composeId string, status string) (*models.Project, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"status": status}}

	var updatedProject models.Project
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"composeId": composeId}, update, opts).Decode(&updatedProject)
	if err != nil {
		return nil, err
	}

	return &updatedProject, nil
}

func (r *ProjectRepository) UpdateProjectUrlZip(ctx context.Context, composeId string, url string) (*models.Project, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"urlZip": url}}

	var updatedProject models.Project
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"composeId": composeId}, update, opts).Decode(&updatedProject)
	if err != nil {
		return nil, err
	}

	return &updatedProject, nil
}

func (r *ProjectRepository) UpdateProjectUrlDeploy(ctx context.Context, composeId string, url string) (*models.Project, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"urlDeploy": url}}

	var updatedProject models.Project
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"composeId": composeId}, update, opts).Decode(&updatedProject)
	if err != nil {
		return nil, err
	}

	return &updatedProject, nil
}

func (r *ProjectRepository) getProjectByComposeId(ctx context.Context, composeId string) (*models.Project, error) {
	var project *models.Project
	err := r.collection.FindOne(ctx, bson.M{"composeId": composeId}).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}
