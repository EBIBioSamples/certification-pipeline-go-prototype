package main

import (
	"flag"
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/pipeline"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	configFile := flag.String("config", "config.json", "the path of the config file")
	schemaFile := flag.String("schema", "config-schema.json", "the path of the config schema file")
	flag.Parse()
	configFileValue := *configFile
	schemaFileValue := *schemaFile
	jsonSubmitted := make(chan string)
	logger := log.New(os.Stdout, "Certification Pipeline ", log.LstdFlags|log.Lshortfile)
	c, _ := config.NewConfig(logger, configFileValue, schemaFileValue)
	piplineFinished := pipeline.NewPipeline(c, jsonSubmitted)
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: json | bscurate")
		return
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	jsonSubmitted <- string(bytes)
	<-piplineFinished
}
