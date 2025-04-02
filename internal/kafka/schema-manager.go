package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/linkedin/goavro/v2"
	"net/http"
	"project-service/internal/config"
	"sync"
)

// topic -> codec
var schemasForThisService = map[string]*goavro.Codec{
	"ProjectStatus": nil,
	"NewZip":        nil,
	"DeployPayload": nil,
}

type SchemaManager struct {
	mu                sync.RWMutex
	Schemas           map[string]*goavro.Codec
	schemaRegistryURL string
}

func NewSchemaManager(cfg *config.Config) *SchemaManager {
	manager := &SchemaManager{
		Schemas:           schemasForThisService,
		schemaRegistryURL: cfg.SchemaRegistryUrl,
	}

	manager.loadSchemasFromRegistry()

	return manager
}

func (sm *SchemaManager) loadSchemasFromRegistry() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for topic := range sm.Schemas {
		schemaData, err := sm.fetchSchemaFromRegistry(topic)
		if err != nil {
			panic(fmt.Sprintf("Failed to load schema for topic %s: %v", topic, err))
		}

		codec, err := goavro.NewCodec(schemaData)
		if err != nil {
			panic(fmt.Sprintf("Failed to create codec for topic %s: %v", topic, err))
		}

		sm.Schemas[topic] = codec
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
