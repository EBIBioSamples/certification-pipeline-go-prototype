package interrogator

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/validator"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
)

//Interrogator controls the workflow of the pipeline
type Interrogator struct {
	Logger        *log.Logger
	Validator     *validator.Validator
	SampleCreated chan string
	Checklists    map[string]string
}

func (i *Interrogator) handleEvents(sampleCreated chan string) {
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
func (i *Interrogator) Interrogate(sample string) []string {
	var candidates = make([]string, 0)
	for key, checklist := range i.Checklists {
		schema, err := ioutil.ReadFile(checklist)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := i.Validator.Validate(string(schema), sample)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}
		if vr.Valid {
			candidates = append(candidates, key)
		}
	}
	return candidates
}
