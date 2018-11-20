package orchestrator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/certifier"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/creator"
	"github.com/EBIBioSamples/certification-pipeline/internal/curator"
	"github.com/EBIBioSamples/certification-pipeline/internal/interrogator"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/reporter"
	"github.com/EBIBioSamples/certification-pipeline/internal/validator"
	"log"
)

var (
	sampleCreated      = make(chan model.Sample)
	sampleInterrogated = make(chan model.InterrogationResult)
	planCompleted      = make(chan model.PlanResult)
	certificateIssued  = make(chan model.Certificate)
)

type Orchestrator struct {
	logger *log.Logger
}

func NewOrchestrator(logger *log.Logger, c *config.Config, jsonSubmitted chan string) *Orchestrator {
	creator.NewCreator(
		logger,
		jsonSubmitted,
		sampleCreated,
	)
	interrogator.NewInterrogator(
		logger,
		&validator.Validator{},
		sampleCreated,
		sampleInterrogated,
		c.Checklists,
	)
	curator.NewCurator(
		logger,
		sampleInterrogated,
		planCompleted,
		certificateIssued,
		c.Plans,
	)
	certifier.NewCertifier(
		logger,
		&validator.Validator{},
		planCompleted,
		certificateIssued,
		c.Checklists,
	)
	reporter.NewReporter(
		logger,
		certificateIssued,
	)
	o := Orchestrator{
		logger: logger,
	}
	return &o
}
