package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/creator"
	"github.com/EBIBioSamples/curation-pipeline/internal/curator"
	"github.com/EBIBioSamples/curation-pipeline/internal/interrogator"
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/EBIBioSamples/curation-pipeline/internal/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	logger                = log.New(os.Stdout, "Curation Pipeline ", log.LstdFlags|log.Lshortfile)
	serverPort            = os.Getenv("SERVER_PORT")
	sampleCreated         = make(chan model.Sample)
	sampleInterrogated    = make(chan model.InterrogationResult)
	curationPlanCompleted = make(chan model.CurationPlanResult)
	checklists            = []model.Checklist{
		{Name: "NCBI Candidate Checklist", File: "./res/schemas/ncbi-candidate-schema.json"},
		{Name: "BioSamples Checklist", File: "./res/schemas/biosamples-schema.json"},
	}
	curationPlans []model.CurationPlan
	cr            = creator.NewCreator(logger, sampleCreated)
)

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sample := cr.CreateSample(buf.String())
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(sample.UUID)
}

func init() {
	if serverPort == "" {
		log.Fatal("$SERVER_PORT not set")
	}
	interrogator.NewInterrogator(
		logger,
		&validator.Validator{},
		sampleCreated,
		sampleInterrogated,
		checklists,
	)
	curator.NewCurator(
		logger,
		sampleInterrogated,
		curationPlanCompleted,
		curationPlans,
	)
}

func main() {
	logger.Printf("starting curation pipeline service")
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/interrogate", handler).Methods("POST")
	logger.Printf("server starting on port %s", serverPort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r))
}
