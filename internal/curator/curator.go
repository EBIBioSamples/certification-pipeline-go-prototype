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
	curator := Curator{
		logger:                   logger,
		sampleInterrogated:       sampleInterrogated,
		curationPlanCompleted:    curationPlanCompleted,
		curationPlansByChecklist: curationPlansByChecklist,
	}
	curator.handleEvents(sampleInterrogated)
	return &curator
}

func (c *Curator) runCurationPlans(ir model.InterrogationResult) {
	for _, cc := range ir.CandidateChecklists {
		cp := c.curationPlansByChecklist[cc]
		s := cp.Execute(ir.Sample)
		c.logger.Printf("curated sample %s", s.UUID)
	}
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
