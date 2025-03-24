package app

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	grpcapp "project-service/internal/app/grpc"
	"project-service/internal/config"
	"project-service/internal/kafka"
	"project-service/internal/repository/project"
	projectservice "project-service/internal/services/project"
)

type App struct {
	GrpcApp *grpcapp.GrpcApp
}

func NewApp(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	mongoOptions := options.Client().ApplyURI(cfg.MongoURL)
	mongoClient, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		log.Error(err.Error())
	}

	projectRepository := project.NewProjectRepository(mongoClient, cfg.MongoDB, "project")
	projectService := projectservice.NewProjectService(log, projectRepository)

	schemaManager := kafka.NewSchemaManager(cfg)
	kafkaConsumer := kafka.NewKafkaConsumer(
		cfg,
		log,
		schemaManager,
		projectService,
	)
	kafkaConsumer.InitConsumer()

	grpcApp := grpcapp.NewGrpcApp(
		log,
		projectService,
		cfg.GRPC.Port,
	)

	return &App{
		GrpcApp: grpcApp,
	}
}
