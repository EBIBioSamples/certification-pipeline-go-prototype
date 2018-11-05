package interrogator

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/EBIBioSamples/curation-pipeline/internal/validator"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
)

//Interrogator controls the workflow of the pipeline
type Interrogator struct {
	logger             *log.Logger
	validator          *validator.Validator
	sampleCreated      chan model.Sample
	sampleInterrogated chan string
	checklists         map[string]string
}

func (i *Interrogator) handleEvents(sampleCreated chan model.Sample) {
	go func() {
		for {
			select {
			case s := <-sampleCreated:
				i.Interrogate(s)
			}
		}
	}()
}

//Interrogate a given sample to find out which checklists it complies to
func (i *Interrogator) Interrogate(sample model.Sample) []string {
	var candidates = make([]string, 0)
	for key, checklist := range i.checklists {
		i.logger.Printf("checking %s against %s\n", sample.UUID, key)
		schema, err := ioutil.ReadFile(checklist)
		if err != nil {
			i.logger.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := i.validator.Validate(string(schema), sample.Document)
		if err != nil {
			i.logger.Fatal(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}
		if vr.Valid {
			candidates = append(candidates, key)
		}
	}
	return candidates
}

func NewInterrogator(
	logger *log.Logger,
	validator *validator.Validator,
	sampleCreated chan model.Sample,
	sampleInterrogated chan string,
	checklists map[string]string) *Interrogator {
	i := Interrogator{
		logger:             logger,
		validator:          validator,
		sampleCreated:      sampleCreated,
		sampleInterrogated: sampleInterrogated,
		checklists:         checklists,
	}
	i.handleEvents(sampleCreated)
	return &i
}
