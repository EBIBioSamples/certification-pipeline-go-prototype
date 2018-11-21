package config

import (
	"encoding/json"
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Logger     *log.Logger
	Checklists []model.Checklist `json:"checklists"`
	Plans      []model.Plan      `json:"plans"`
}

type ConfigError struct {
	message          string
	validationErrors []string
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("%s - %s", ce.message, ce.validationErrors)
}

func NewConfig(logger *log.Logger, configFilePath string, configSchemaFilePath string) (*Config, error) {
	schemaFile, err := os.Open(configSchemaFilePath)
	defer schemaFile.Close()
	if err != nil {
		logger.Panic(err)
	}
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		logger.Panic(err)
	}
	schemaBytes, _ := ioutil.ReadAll(schemaFile)
	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
	configBytes, _ := ioutil.ReadAll(configFile)
	documentLoader := gojsonschema.NewStringLoader(string(configBytes))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		logger.Panic(err)
	}
	if !result.Valid() {
		ce := ConfigError{message: "The config is not valid"}
		for _, desc := range result.Errors() {
			ce.validationErrors = append(ce.validationErrors, desc.Description())
		}
		return nil, ce
	}
	var config Config
	json.Unmarshal([]byte(configBytes), &config)
	config.Logger = logger
	return &config, nil
}
