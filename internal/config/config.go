package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env               string // dev || prod
	GRPC              GRPCConfig
	SchemaRegistryUrl string
	KafkaHost         string
	MongoURL          string
	MongoDB           string
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
}

func MustLoad() *Config {
	loadEnvFile()

	env := getEnv("ENV", "dev")
	grpcPort := getEnvAsInt("GRPC_PORT", 50051)
	grpcTimeout := getEnvAsDuration("GRPC_TIMEOUT", 10*time.Second)
	schemaRegistryUrl := getEnv("SCHEMA_REGISTRY_URL", "http://localhost:6767")
	kafkaHost := getEnv("KAFKA_HOST", "http://localhost:9092")
	mongoUrl := buildMongoURL()
	mongoDb := getEnv("MONGO_DB", "project-service-db")

	return &Config{
		Env: env,
		GRPC: GRPCConfig{
			Port:    grpcPort,
			Timeout: grpcTimeout,
		},
		SchemaRegistryUrl: schemaRegistryUrl,
		KafkaHost:         kafkaHost,
		MongoURL:          mongoUrl,
		MongoDB:           mongoDb,
	}
}

func loadEnvFile() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func buildMongoURL() string {
	user := getEnv("MONGO_USER", "root")
	password := getEnv("MONGO_PASSWORD", "password")
	db := getEnv("MONGO_DB", "project-service-db")
	port := getEnv("MONGO_PORT", "27017")
	host := getEnv("MONGO_HOST", "localhost")

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=project-service-db", user, password, host, port, db)
}
