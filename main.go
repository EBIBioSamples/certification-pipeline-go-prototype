package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/pipeline"
	"github.com/EBIBioSamples/certification-pipeline/internal/reporter"
	"github.com/gorilla/mux"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	logger        = log.New(os.Stdout, "Certification Pipeline ", log.LstdFlags|log.Lshortfile)
	serverPort    = os.Getenv("SERVER_PORT")
	graylogAddr   = os.Getenv("GRAYLOG_URL")
	c, _          = config.NewConfig(logger, "./res/config.json", "./res/schemas/config-schema.json")
	jsonSubmitted = make(chan string)
	rep           *reporter.Reporter
)

func interrogateHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	jsonSubmitted <- buf.String()
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(fmt.Sprintf("http://%s/sample/%s", r.Host, sample.UUID))
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	result, ok := rep.SampleInfo(uuid)
	if !ok {
		http.Error(w, "Sample not found", 404)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func init() {
	if serverPort == "" {
		logger.Fatal("$SERVER_PORT not set")
	}
	if graylogAddr != "" {
		gelfWriter, err := gelf.NewWriter(graylogAddr)
		if err != nil {
			log.Fatalf("gelf.NewWriter: %s", err)
		}
		// log to both stderr and graylog2
		logger.SetOutput(io.MultiWriter(os.Stdout, gelfWriter))
		logger.Printf("logging to stdout & graylog2@'%s'", graylogAddr)
	}

	checklistMap := make(map[string]model.Checklist)
	for _, checklist := range c.Checklists {
		checklistMap[checklist.Name] = checklist
	}

}

func main() {
	logger.Printf("creating certification pipeline")
	pipeline.NewPipeline(c, jsonSubmitted)
	logger.Printf("starting service")
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/interrogate", interrogateHandler).Methods("POST")
	r.HandleFunc("/sample/{uuid}", sampleHandler).Methods("GET")
	logger.Printf("server starting on port %s", serverPort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r))
}
