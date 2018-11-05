package coordinator

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/validator"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
)

//Coordinator controls the workflow of the pipeline
type Coordinator struct {
	Logger        *log.Logger
	Validator     *validator.Validator
	SampleCreated chan string
	Checklists    map[string]string
}

func (c *Coordinator) handleEvents(sampleCreated chan string) {
	go func() {
		for {
			select {
			case s := <-sampleCreated:
				c.FindCandidates(s)
			}
		}
	}()
}

func (c *Coordinator) FindCandidates(sample string) []string {
	var candidates = make([]string, 0)
	for key, checklist := range c.Checklists {
		schema, err := ioutil.ReadFile(checklist)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := c.Validator.Validate(string(schema), sample)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}
		if vr.Valid {
			candidates = append(candidates, key)
		}
	}
	return candidates
}
