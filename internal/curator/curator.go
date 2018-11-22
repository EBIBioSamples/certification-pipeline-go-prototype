package curator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"log"
)

type Curator struct {
	logger                      *log.Logger
	checklistMatched            chan model.ChecklistMatches
	planCompleted               chan model.PlanResult
	curationCompleted           chan model.CurationEnd
	certificateIssued           chan model.Certificate
	plansByCandidateChecklistID map[string]model.Plan
}

func NewCurator(logger *log.Logger, checklistMatched chan model.ChecklistMatches, plans []model.Plan) (planCompleted chan model.PlanResult, curationCompleted chan model.CurationEnd) {
	plansByCandidateChecklistID := make(map[string]model.Plan)
	for _, p := range plans {
		plansByCandidateChecklistID[p.CandidateChecklistID] = p
	}
	c := Curator{
		logger:                      logger,
		checklistMatched:            checklistMatched,
		planCompleted:               make(chan model.PlanResult),
		curationCompleted:           make(chan model.CurationEnd),
		plansByCandidateChecklistID: plansByCandidateChecklistID,
	}
	c.handleEvents(checklistMatched)
	return c.planCompleted, c.curationCompleted
}

func (c *Curator) runCurationPlans(ir model.ChecklistMatches) {
	s := ir.Sample
	if len(ir.Checklists) == 0 {
		c.curationCompleted <- model.CurationEnd{Sample: s}
	}
	for _, cc := range ir.Checklists {
		c.runCurationPlan(cc, s)
	}
}

func (c *Curator) runCurationPlan(checklist model.Checklist, s model.Sample) {
	if _, ok := c.plansByCandidateChecklistID[checklist.ID()]; !ok {
		c.curationCompleted <- model.CurationEnd{Sample: s, Checklist: checklist}
	} else {
		p := c.plansByCandidateChecklistID[checklist.ID()]
		s = p.Execute(s)
		pr := model.PlanResult{
			Sample: s,
			Plan:   p,
		}
		c.planCompleted <- pr
	}
}

func (c *Curator) handleEvents(in chan model.ChecklistMatches) {
	go func() {
		for {
			select {
			case ir := <-in:
				c.runCurationPlans(ir)
			}
		}
	}()
}
