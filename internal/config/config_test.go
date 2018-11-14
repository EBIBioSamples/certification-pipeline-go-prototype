package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	logger     = log.New(os.Stdout, "TestLoader ", log.LstdFlags|log.Lshortfile)
	configFile = "../../res/config.json"
)

func TestLoader(t *testing.T) {
	config := NewConfig(logger, configFile)
	assert.NotEmpty(t, config.Checklists)
	for _, c := range config.Checklists {
		assert.NotNil(t, c.Name)
		assert.NotNil(t, c.Version)
		assert.NotNil(t, c.File)
	}
}
