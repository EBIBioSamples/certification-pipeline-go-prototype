package curator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"log"
)

type Curator struct {
	logger                      *log.Logger
	sampleInterrogated          chan model.InterrogationResult
	planCompleted               chan model.PlanResult
	certificateIssued           chan model.Certificate
	plansByCandidateChecklistID map[string]model.Plan
}

func NewCurator(logger *log.Logger, sampleInterrogated chan model.InterrogationResult,
	planCompleted chan model.PlanResult, certificateIssued chan model.Certificate, plans []model.Plan) *Curator {
	plansByCandidateChecklistID := make(map[string]model.Plan)
	for _, p := range plans {
		plansByCandidateChecklistID[p.CandidateChecklistID] = p
	}
	curator := Curator{
		logger:                      logger,
		sampleInterrogated:          sampleInterrogated,
		planCompleted:               planCompleted,
		certificateIssued:           certificateIssued,
		plansByCandidateChecklistID: plansByCandidateChecklistID,
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
	if _, ok := c.plansByCandidateChecklistID[checklist.ID()]; !ok {
		c.logger.Printf("no curation plans for %s", checklist.ID())
		return
	}
	p := c.plansByCandidateChecklistID[checklist.ID()]
	s = p.Execute(s)
	pr := model.PlanResult{
		Sample: s,
		Plan:   p,
	}
	c.planCompleted <- pr
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
