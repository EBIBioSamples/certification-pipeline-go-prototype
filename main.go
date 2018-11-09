package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/certifier"
	"github.com/EBIBioSamples/certification-pipeline/internal/creator"
	"github.com/EBIBioSamples/certification-pipeline/internal/curator"
	"github.com/EBIBioSamples/certification-pipeline/internal/interrogator"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/reporter"
	"github.com/EBIBioSamples/certification-pipeline/internal/validator"
	"github.com/gorilla/mux"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	logger                = log.New(os.Stdout, "Certification Pipeline ", log.LstdFlags|log.Lshortfile)
	serverPort            = os.Getenv("SERVER_PORT")
	graylogAddr           = os.Getenv("GRAYLOG_URL")
	sampleCreated         = make(chan model.Sample)
	sampleInterrogated    = make(chan model.InterrogationResult)
	curationPlanCompleted = make(chan model.CurationPlanResult)
	certificateIssued     = make(chan model.Certificate)
	checklists            = []model.Checklist{
		{Name: "NCBI Candidate Checklist", File: "./res/schemas/ncbi-candidate-schema.json"},
		{Name: "BioSamples Checklist", File: "./res/schemas/biosamples-schema.json"},
	}
	curationPlans []model.CurationPlan
	cr            = creator.NewCreator(logger, sampleCreated)
	rep           *reporter.Reporter
)

func interrogateHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sample := cr.CreateSample(buf.String())
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(fmt.Sprintf("http://%s/sample/%s", r.Host, sample.UUID))
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
		logger.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
		logger.Printf("logging to stderr & graylog2@'%s'", graylogAddr)
	}

	checklistMap := make(map[string]model.Checklist)
	for _, checklist := range checklists {
		checklistMap[checklist.Name] = checklist
	}
	curationPlans = []model.CurationPlan{
		{
			Logger:        logger,
			Name:          "NCBI to BioSamples",
			FromChecklist: checklistMap["NCBI Candidate Checklist"],
			ToChecklist:   checklistMap["BioSamples Checklist"],
			Curations: []model.Curation{
				{
					Characteristic: "INSDC status",
					NewValue:       "public",
				},
			},
		},
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
		certificateIssued,
		curationPlans,
	)
	certifier.NewCertifier(
		logger,
		&validator.Validator{},
		curationPlanCompleted,
		certificateIssued,
		checklists,
	)
	rep = reporter.NewReporter(
		logger,
		certificateIssued,
	)
}

func main() {
	logger.Printf("starting curation pipeline service")
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/interrogate", interrogateHandler).Methods("POST")
	r.HandleFunc("/sample/{uuid}", sampleHandler).Methods("GET")
	logger.Printf("server starting on port %s", serverPort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r))
}
