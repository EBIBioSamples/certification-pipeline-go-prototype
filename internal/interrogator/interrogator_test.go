package interrogator_test

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/interrogator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	sampleCreated = make(chan string)
	checklists    = map[string]string{
		"NCBI Checklist":       "../../res/schemas/ncbi-schema.json",
		"BioSamples Checklist": "../../res/schemas/biosamples-schema.json",
	}
)

func TestInterrogate(t *testing.T) {
	tests := []struct {
		documentFile       string
		expectedCandidates []string
	}{
		{
			documentFile:       "../../res/json/ncbi-SAMN03894263.json",
			expectedCandidates: []string{"NCBI Checklist"},
		},
	}
	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		i := interrogator.Interrogator{
			Logger:        log.New(os.Stdout, "TestInterrogate ", log.LstdFlags|log.Lshortfile),
			SampleCreated: sampleCreated,
			Checklists:    checklists,
		}
		candidates := i.Interrogate(string(document))
		assert.Equal(t, test.expectedCandidates, candidates)
	}
}
