package curator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"log"
)

type Curator struct {
	logger                      *log.Logger
	checklistMatched            chan model.ChecklistMatches
	planCompleted               chan model.PlanResult
	certificateIssued           chan model.Certificate
	plansByCandidateChecklistID map[string]model.Plan
}

func NewCurator(logger *log.Logger, checklistMatched chan model.ChecklistMatches, plans []model.Plan) chan model.PlanResult {
	plansByCandidateChecklistID := make(map[string]model.Plan)
	for _, p := range plans {
		plansByCandidateChecklistID[p.CandidateChecklistID] = p
	}
	c := Curator{
		logger:                      logger,
		checklistMatched:            checklistMatched,
		planCompleted:               make(chan model.PlanResult),
		plansByCandidateChecklistID: plansByCandidateChecklistID,
	}
	c.handleEvents(checklistMatched)
	return c.planCompleted
}

func (c *Curator) runCurationPlans(ir model.ChecklistMatches) {
	s := ir.Sample
	for _, cc := range ir.Checklists {
		c.runCurationPlan(cc, s)
	}
}

func (c *Curator) runCurationPlan(checklist model.Checklist, s model.Sample) {
	if _, ok := c.plansByCandidateChecklistID[checklist.ID()]; !ok {
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

func (c *Curator) handleEvents(in chan model.ChecklistMatches) {
	go func() {
		for {
			select {
			case ir := <-in:
				c.runCurationPlans(ir)
			case cert := <-c.certificateIssued:
				c.runCurationPlan(cert.Checklist, cert.Sample)
			}
		}
	}()
}
