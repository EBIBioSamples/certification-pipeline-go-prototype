package curator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"log"
)

type Curator struct {
	logger                   *log.Logger
	sampleInterrogated       chan model.InterrogationResult
	curationPlanCompleted    chan model.CurationPlanResult
	certificateIssued        chan model.Certificate
	curationPlansByChecklist map[model.Checklist]model.CurationPlan
}

func NewCurator(logger *log.Logger, sampleInterrogated chan model.InterrogationResult,
	curationPlanCompleted chan model.CurationPlanResult, certificateIssued chan model.Certificate, curationPlans []model.CurationPlan) *Curator {
	curationPlansByChecklist := make(map[model.Checklist]model.CurationPlan)
	for _, cp := range curationPlans {
		curationPlansByChecklist[cp.FromChecklist] = cp
	}
	curator := Curator{
		logger:                   logger,
		sampleInterrogated:       sampleInterrogated,
		curationPlanCompleted:    curationPlanCompleted,
		certificateIssued:        certificateIssued,
		curationPlansByChecklist: curationPlansByChecklist,
	}
	curator.handleEvents(sampleInterrogated)
	return &curator
}

func (c *Curator) runCurationPlans(ir model.InterrogationResult) {
	s := ir.Sample
	for _, cc := range ir.CandidateChecklists {
		c.runCurationPlan(cc, s)
	}
}

func (c *Curator) runCurationPlan(checklist model.Checklist, s model.Sample) {
	if _, ok := c.curationPlansByChecklist[checklist]; !ok {
		c.logger.Printf("no curation plans for %s", checklist.Name)
		return
	}
	cp := c.curationPlansByChecklist[checklist]
	s = cp.Execute(s)
	cpr := model.CurationPlanResult{
		Sample:       s,
		CurationPlan: cp,
	}
	c.curationPlanCompleted <- cpr
}

func (c *Curator) handleEvents(sampleInterrogated chan model.InterrogationResult) {
	go func() {
		for {
			select {
			case ir := <-sampleInterrogated:
				c.runCurationPlans(ir)
			case cert := <-c.certificateIssued:
				c.runCurationPlan(cert.Checklist, cert.Sample)
			}
		}
	}()
}
