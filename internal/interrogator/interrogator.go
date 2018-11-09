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
	validator          *validator.Validator
	sampleCreated      chan model.Sample
	sampleInterrogated chan model.InterrogationResult
	checklistMap       map[string]model.Checklist
}

//Interrogate a given sample to find out which checklists it complies to
func (i *Interrogator) interrogate(sample model.Sample) {
	var candidates = make([]model.Checklist, 0)
	for key, checklist := range i.checklistMap {
		i.logger.Printf("checking %s against %s\n", sample.UUID, key)
		schema, err := ioutil.ReadFile(checklist.File)
		if err != nil {
			i.logger.Panic(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := i.validator.Validate(string(schema), sample.Document)
		if err != nil {
			i.logger.Panic(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}
		if vr.Valid {
			candidates = append(candidates, checklist)
		}
	}
	if len(candidates) > 0 {
		ir := model.InterrogationResult{
			Sample:              sample,
			CandidateChecklists: candidates,
		}
		i.logger.Printf("sample %s matches %s", sample.UUID, ir.CandidateChecklists)
		i.sampleInterrogated <- ir
	}
}

//NewInterrogator returns a new instance of an Interrogate with the specified checklists
func NewInterrogator(
	logger *log.Logger,
	validator *validator.Validator,
	sampleCreated chan model.Sample,
	sampleInterrogated chan model.InterrogationResult,
	checklists []model.Checklist) *Interrogator {
	checklistMap := make(map[string]model.Checklist)
	for _, checklist := range checklists {
		checklistMap[checklist.Name] = checklist
	}
	i := Interrogator{
		logger:             logger,
		validator:          validator,
		sampleCreated:      sampleCreated,
		sampleInterrogated: sampleInterrogated,
		checklistMap:       checklistMap,
	}
	i.handleEvents(sampleCreated)
	return &i
}

func (i *Interrogator) handleEvents(sampleCreated chan model.Sample) {
	go func() {
		for {
			select {
			case s := <-sampleCreated:
				i.interrogate(s)
			}
		}
	}()
}
