package config

import (
	"encoding/json"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	logger     *log.Logger
	Checklists []model.Checklist `json:"checklists"`
	Plans      []model.Plan      `json:"plans"`
}

func NewConfig(logger *log.Logger, configFile string) *Config {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		logger.Panic(err)
	}
	defer jsonFile.Close()
	var config Config
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &config)
	config.logger = logger
	return &config
}
