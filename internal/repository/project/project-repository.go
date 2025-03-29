package project

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProjectUniqueIdentifier представляет уникальный идентификатор проекта
type ProjectUniqueIdentifier struct {
	Owner string `bson:"owner"`
	Name  string `bson:"name"`
}

// ProjectStatus представляет статус проекта
type ProjectStatus int

const (
	NEW ProjectStatus = iota
	GENERATE_PENDING = 2
	GENERATE_SUCCESS = 3
	GENERATE_FAIL    = 4
	DEPLOY_PENDING   = 5
	DEPLOY_SUCCESS   = 6
	DEPLOY_FAIL      = 7
	RUNNING          = 8
	STOPPED          = 9
	FAILED           = 10
)

// Project представляет структуру проекта в базе данных
type Project struct {
	ComposeID ProjectUniqueIdentifier `bson:"compose_id"`
	Data      map[string]interface{}  `bson:"data,omitempty"`
	Status    ProjectStatus           `bson:"status"`
}

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(client *mongo.Client, dbName, collectionName string) *ProjectRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProjectRepository{collection: collection}
}

// CreateUniqueIndex создает уникальный индекс для поля compose_id
func (r *ProjectRepository) CreateUniqueIndex(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "compose_id.owner", Value: 1},
				{Key: "compose_id.name", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}

// GetProjectByID получает проект по его уникальному идентификатору
func (r *ProjectRepository) GetProjectByID(ctx context.Context, id ProjectUniqueIdentifier) (*Project, error) {
	filter := bson.M{"compose_id": bson.M{"owner": id.Owner, "name": id.Name}}
	var project Project
	err := r.collection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Проект не найден
		}
		return nil, err
	}
	return &project, nil
}

// GetAllUserProjects получает все проекты пользователя с пагинацией
func (r *ProjectRepository) GetAllUserProjects(ctx context.Context, owner string, page, limit int64) ([]*Project, error) {
	filter := bson.M{"compose_id.owner": owner}
	
	options := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.M{"compose_id.name": 1})
	
	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var projects []*Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	
	return projects, nil
}

// InitProject инициализирует новый проект
func (r *ProjectRepository) InitProject(ctx context.Context, id ProjectUniqueIdentifier) (*Project, error) {
	// Проверяем, что проект с таким идентификатором не существует
	existingProject, err := r.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingProject != nil {
		return nil, errors.New("проект с таким названием уже существует для данного пользователя")
	}
	
	// Создаем новый проект
	project := &Project{
		ComposeID: id,
		Status:    NEW,
		Data:      make(map[string]interface{}),
	}
	
	_, err = r.collection.InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}
	
	return project, nil
}

// UpdateProject обновляет данные проекта
func (r *ProjectRepository) UpdateProject(ctx context.Context, id ProjectUniqueIdentifier, data map[string]interface{}) (*Project, error) {
	// Проверяем, что проект существует
	existingProject, err := r.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingProject == nil {
		return nil, errors.New("проект не найден")
	}
	
	// Обновляем данные проекта
	update := bson.M{
		"$set": bson.M{
			"data": data,
		},
	}
	
	filter := bson.M{"compose_id": bson.M{"owner": id.Owner, "name": id.Name}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	
	// Получаем обновленный проект
	return r.GetProjectByID(ctx, id)
}

// UpdateProjectStatus обновляет статус проекта
func (r *ProjectRepository) UpdateProjectStatus(ctx context.Context, id ProjectUniqueIdentifier, status ProjectStatus) error {
	filter := bson.M{"compose_id": bson.M{"owner": id.Owner, "name": id.Name}}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
