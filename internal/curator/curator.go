package curator

import (
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"log"
)

type Curator struct {
	logger                   *log.Logger
	sampleInterrogated       chan model.InterrogationResult
	curationPlanCompleted    chan model.CurationPlanResult
	curationPlansByChecklist map[model.Checklist]model.CurationPlan
}

func NewCurator(logger *log.Logger, sampleInterrogated chan model.InterrogationResult,
	curationPlanCompleted chan model.CurationPlanResult, curationPlans []model.CurationPlan) *Curator {
	curationPlansByChecklist := make(map[model.Checklist]model.CurationPlan)
	for _, cp := range curationPlans {
		curationPlansByChecklist[cp.FromChecklist] = cp
	}
	c := Curator{
		logger:                   logger,
		sampleInterrogated:       sampleInterrogated,
		curationPlanCompleted:    curationPlanCompleted,
		curationPlansByChecklist: curationPlansByChecklist,
	}
	c.handleEvents(sampleInterrogated)
	return &c
}

func (c *Curator) runCurationPlans(ir model.InterrogationResult) {
	c.logger.Printf("at the point a curation plan would run for %s", ir.CandidateChecklists)
}

func (c *Curator) handleEvents(sampleInterrogated chan model.InterrogationResult) {
	go func() {
		for {
			select {
			case ir := <-sampleInterrogated:
				c.runCurationPlans(ir)
			}
		}
	}()
}
