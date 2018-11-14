package config_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	logger = log.New(os.Stdout, "TestLoader ", log.LstdFlags|log.Lshortfile)
)

func TestConfigWithValidFile(t *testing.T) {
	validConfigFile := "../../res/config.json"
	config, err := config.NewConfig(logger, validConfigFile, "../../res/schemas/config-schema.json")
	assert.Nil(t, err)
	assert.NotEmpty(t, config.Checklists)
	for _, c := range config.Checklists {
		assert.NotEmpty(t, c.Name)
		assert.NotEmpty(t, c.Version)
		assert.NotEmpty(t, c.File)
	}
	assert.NotEmpty(t, config.Plans)
	for _, p := range config.Plans {
		assert.NotEmpty(t, p.CandidateChecklistID)
		assert.NotEmpty(t, p.CertificateChecklistID)
		assert.NotEmpty(t, p.Curations)
		fmt.Println(p.Curations)
	}
}

func TestConfigWithInvalidFile(t *testing.T) {
	validConfigFile := "../../res/invalid_config.json"
	_, err := config.NewConfig(logger, validConfigFile, "../../res/schemas/config-schema.json")
	assert.NotNil(t, err)
}
