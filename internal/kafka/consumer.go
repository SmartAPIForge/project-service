package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
	"log/slog"
	"project-service/internal/config"
	"project-service/internal/dto"
	projectservice "project-service/internal/services/project"
)

type KafkaConsumer struct {
	log            *slog.Logger
	consumer       *kafka.Consumer
	topic          string
	codec          *goavro.Codec
	projectService *projectservice.ProjectService
}

func NewKafkaConsumer(
	log *slog.Logger,
	cfg *config.Config,
	topic string,
	codec *goavro.Codec,
	projectService *projectservice.ProjectService,
) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaHost,
		"group.id":           "project-psg",
		"enable.auto.commit": false,
	})
	if err != nil {
		panic(fmt.Sprintf("Error creating kafka consumer %v", err))
	}

	return &KafkaConsumer{
		log:            log,
		consumer:       consumer,
		topic:          topic,
		codec:          codec,
		projectService: projectService,
	}
}

func (kc *KafkaConsumer) Sub() {
	err := kc.consumer.Subscribe(kc.topic, nil)
	if err != nil {
		kc.log.Error("Error subscribing to topic: ", kc.topic, err)
	}
}

// Consume > ~1 minute wait for assigning
func (kc *KafkaConsumer) Consume() {
	switch kc.topic {
	case "ProjectStatus":
		kc.consumeProjectStatus()
		break
	case "NewZip":
		kc.consumeProjectNewZip()
		break
	case "DeployPayload":
		kc.consumeDeployPayload()
		break
	}
}

func (kc *KafkaConsumer) consumeProjectStatus() {
	kc.log.Info("Started consuming ProjectStatus")

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("Error reading from topic ProjectStatus:", err)
			continue
		}

		kc.log.Info("New message from topic ProjectStatus")

		native, _, err := kc.codec.NativeFromTextual(msg.Value)
		if err != nil {
			kc.log.Error("Incorrect message while handling ProjectStatus:", string(msg.Value), err)
			kc.commitMessage(msg)
			continue
		}

		projectStatusDTO := dto.MapNativeToProjectStatusDTO(native)
		canCommit, _ := kc.projectService.UpdateProjectStatus(context.Background(), projectStatusDTO)
		if canCommit {
			kc.commitMessage(msg)
		}
	}
}

func (kc *KafkaConsumer) consumeProjectNewZip() {
	kc.log.Info("Started consuming NewZip")

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("Error reading from topic NewZip:", err)
			continue
		}

		kc.log.Info("New message from topic NewZip")

		native, _, err := kc.codec.NativeFromTextual(msg.Value)
		if err != nil {
			kc.log.Error("Incorrect message while handling NewZip:", string(msg.Value), err)
			kc.commitMessage(msg)
			continue
		}

		newZipDTO := dto.MapNativeToNewZipDTO(native)
		canCommit, _ := kc.projectService.UpdateProjectUrlZip(context.Background(), newZipDTO)
		if canCommit {
			kc.commitMessage(msg)
		}
	}
}

func (kc *KafkaConsumer) consumeDeployPayload() {
	kc.log.Info("Started consuming DeployPayload")

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("Error reading from topic DeployPayload:", err)
			continue
		}

		kc.log.Info("New message from topic DeployPayload")

		native, _, err := kc.codec.NativeFromTextual(msg.Value)
		if err != nil {
			kc.log.Error("Incorrect message while handling DeployPayload:", string(msg.Value), err)
			kc.commitMessage(msg)
			continue
		}

		deployPayloadDTO := dto.MapNativeToDeployPayloadDTO(native)
		canCommit, _ := kc.projectService.UpdateProjectUrlDeploy(context.Background(), deployPayloadDTO)
		if canCommit {
			kc.commitMessage(msg)
		}
	}
}

func (kc *KafkaConsumer) commitMessage(msg *kafka.Message) {
	_, err := kc.consumer.CommitMessage(msg)
	if err != nil {
		kc.log.Error(
			"Failed to commit message",
			"topic", *msg.TopicPartition.Topic,
			"partition", msg.TopicPartition.Partition,
			"offset", msg.TopicPartition.Offset,
			"error", err,
		)
	}
}
