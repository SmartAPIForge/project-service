package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/linkedin/goavro/v2"
	"net/http"
	"os"
	"path"
	"project-service/internal/config"
	"sync"
)

var schemasFromThisService = map[string]string{
	"ProjectStatus": path.Join("avro", "project-status.avsc"),
}

var schemasForThisService = map[string]*goavro.Codec{
	"NewZip":        nil,
	"ProjectStatus": nil,
}

type SchemaManager struct {
	mu                sync.RWMutex
	schemas           map[string]*goavro.Codec
	schemaRegistryURL string
}

func NewSchemaManager(cfg *config.Config) *SchemaManager {
	manager := &SchemaManager{
		schemas:           schemasForThisService,
		schemaRegistryURL: cfg.SchemaRegistryUrl,
	}

	manager.uploadSchemasToRegistry()
	manager.loadSchemasFromRegistry()

	return manager
}

func (sm *SchemaManager) uploadSchemasToRegistry() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for topic, schemaPath := range schemasFromThisService {
		if err := sm.uploadSchemaToRegistry(topic, schemaPath); err != nil {
			panic(fmt.Sprintf("Failed to upload schema for topic %s: %v", topic, err))
		}
		fmt.Printf("Schema for topic %s successfully uploaded to registry\n", topic)
	}
}

func (sm *SchemaManager) uploadSchemaToRegistry(topic, schemaPath string) error {
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	type schemaRequest struct {
		Schema string `json:"schema"`
	}
	requestBody, err := json.Marshal(schemaRequest{Schema: string(schemaData)})
	if err != nil {
		return err
	}

	schemaURL := fmt.Sprintf("%s/subjects/%s-value/versions", sm.schemaRegistryURL, topic)

	resp, err := http.Get(schemaURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		fmt.Printf("Schema for topic %s already exists in registry\n", topic)
		return nil
	}

	req, err := http.NewRequest("POST", schemaURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func (sm *SchemaManager) loadSchemasFromRegistry() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for topic := range sm.schemas {
		schemaData, err := sm.fetchSchemaFromRegistry(topic)
		if err != nil {
			panic(fmt.Sprintf("Failed to load schema for topic %s: %v", topic, err))
		}

		codec, err := goavro.NewCodec(schemaData)
		if err != nil {
			panic(fmt.Sprintf("Failed to create codec for topic %s: %v", topic, err))
		}

		sm.schemas[topic] = codec
		fmt.Printf("Schema for topic %s successfully loaded from registry\n", topic)
	}
}

func (sm *SchemaManager) fetchSchemaFromRegistry(topic string) (string, error) {
	schemaURL := fmt.Sprintf("%s/subjects/%s-value/versions/latest", sm.schemaRegistryURL, topic)
	resp, err := http.Get(schemaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var schemaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		return "", err
	}

	schema, ok := schemaResp["schema"].(string)
	if !ok {
		return "", err
	}

	return schema, nil
}

func (sm *SchemaManager) GetCodec(topic string) (*goavro.Codec, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	codec, exists := sm.schemas[topic]
	if !exists {
		return nil, fmt.Errorf("schema for topic %s not found", topic)
	}

	return codec, nil
}
