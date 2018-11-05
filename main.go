package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/interrogator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	serverPort    = os.Getenv("SERVER_PORT")
	sampleCreated = make(chan string)
	checklists    = map[string]string{
		"NCBI Checklist":       "./res/schemas/ncbi-schema.json",
		"BioSamples Checklist": "./res/schemas/biosamples-schema.json",
	}
	i = interrogator.Interrogator{
		Logger:        log.New(os.Stdout, "Curation Pipeline ", log.LstdFlags|log.Lshortfile),
		SampleCreated: sampleCreated,
		Checklists:    checklists,
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	result := i.Interrogate(buf.String())
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func init() {
	if serverPort == "" {
		log.Fatal("$SERVER_PORT not set")
	}
}

func main() {
	i.Logger.Printf("starting curation pipeline service")
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/interrogate", handler).Methods("POST")
	i.Logger.Printf("server starting on port %s", serverPort)
	i.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r))
}
