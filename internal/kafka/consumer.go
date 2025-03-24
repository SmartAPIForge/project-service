package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"project-service/internal/config"
	projectservice "project-service/internal/services/project"
)

var topics = []string{"ProjectStatus"}

type KafkaConsumer struct {
	log            *slog.Logger
	consumer       *kafka.Consumer
	projectService *projectservice.ProjectService
}

func NewKafkaConsumer(
	cfg *config.Config,
	log *slog.Logger,
	schemaManager *SchemaManager,
	projectService *projectservice.ProjectService,
) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaHost,
		"group.id":          "project-service-group",
	})
	if err != nil {
		panic(fmt.Sprintf("Error creating kafka consumer %v", err))
	}

	return &KafkaConsumer{
		consumer:       consumer,
		log:            log,
		projectService: projectService,
	}
}

func (kc *KafkaConsumer) InitConsumer() {
	err := kc.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		kc.log.Error("Error subscribing to topics", err)
		return
	}

	go func() {
		kc.consumeProjectStatus()
	}()
}

func (kc *KafkaConsumer) consumeProjectStatus() {
	kc.log.Info("Started consuming ProjectStatus")

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("Error reading from topic ProjectStatus:", err)
			continue
		}

		kc.log.Info("New message from topic ProjectStatus:", string(msg.Value))

		// TODO LOGIC
	}
}
