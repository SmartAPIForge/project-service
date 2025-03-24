package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
	"log/slog"
	"project-service/internal/config"
	projectservice "project-service/internal/services/project"
	"sync"
)

type KafkaConsumer struct {
	consumer       *kafka.Consumer
	log            *slog.Logger
	topics         []string
	projectService *projectservice.ProjectService
}

func NewKafkaConsumer(
	cfg *config.Config,
	log *slog.Logger,
	schemaManager *SchemaManager,
	projectService *projectservice.ProjectService,
) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": cfg.KafkaHost})
	if err != nil {
		panic("Error creating kafka consumer")
	}

	return &KafkaConsumer{
		consumer:       consumer,
		topics:         getTopics(schemaManager.schemas),
		log:            log,
		projectService: projectService,
	}
}

func (kc *KafkaConsumer) InitConsumer() {
	err := kc.consumer.SubscribeTopics(kc.topics, nil)
	if err != nil {
		kc.log.Error("Error subscribing to topics", err)
		return
	}

	var wg sync.WaitGroup

	for _, topic := range kc.topics {
		wg.Add(1)

		go func(t string) {
			defer wg.Done()
			switch t {
			case "ProjectStatus":
				kc.consumeProjectStatus()
				break
			default:
				kc.log.Error("SchemaManager provided unknown topic %s", t)
			}
		}(topic)
	}

	wg.Wait()
}

// ProjectStatus topic
func (kc *KafkaConsumer) consumeProjectStatus() {
	kc.log.Info("Started consuming ProjectStatus")

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("Error reading from topic ProjectStatus: %v", err)
			continue
		}

		kc.log.Info("New message from topic ProjectStatus: %s", string(msg.Value))

		// TODO LOGIC
	}
}

func getTopics(codecs map[string]*goavro.Codec) []string {
	topics := make([]string, 0, len(codecs))

	for topic := range codecs {
		topics = append(topics, topic)
	}

	return topics
}
