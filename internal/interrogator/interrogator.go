package interrogator

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/validator"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
)

//Interrogator checks samples against checklists to find candidate samples for curation
type Interrogator struct {
	logger             *log.Logger
	sampleInterrogated chan model.ChecklistMatches
	checklistMap       map[string]model.Checklist
}

//Interrogate a given sample to find out which checklists it complies to
func (i *Interrogator) interrogate(sample model.Sample) {
	var candidates = make([]model.Checklist, 0)
	for _, checklist := range i.checklistMap {
		schema, err := ioutil.ReadFile(checklist.File)
		if err != nil {
			i.logger.Panic(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := validator.Validate(string(schema), sample.Document)
		if err != nil {
			i.logger.Panic(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}
		if vr.Valid {
			candidates = append(candidates, checklist)
		}
	}
	if len(candidates) > 0 {
		ir := model.ChecklistMatches{
			Sample:     sample,
			Checklists: candidates,
		}
		i.sampleInterrogated <- ir
	}
}

//NewInterrogator returns a new instance of an Interrogate with the specified checklists
func NewInterrogator(logger *log.Logger, in chan model.Sample, checklists []model.Checklist) chan model.ChecklistMatches {
	checklistMap := make(map[string]model.Checklist)
	for _, checklist := range checklists {
		checklistMap[checklist.Name] = checklist
	}
	i := Interrogator{
		logger:             logger,
		sampleInterrogated: make(chan model.ChecklistMatches),
		checklistMap:       checklistMap,
	}
	i.handleEvents(in)
	return i.sampleInterrogated
}

func (i *Interrogator) handleEvents(in chan model.Sample) {
	go func() {
		for {
			select {
			case sample := <-in:
				i.interrogate(sample)
			}
		}
	}()
}
