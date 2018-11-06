package interrogator_test

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/interrogator"
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/EBIBioSamples/curation-pipeline/internal/validator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	sampleCreated      = make(chan model.Sample)
	sampleInterrogated = make(chan model.InterrogationResult)
	checklists         = []model.Checklist{
		{Name: "NCBI Candidate Checklist", File: "../../res/schemas/ncbi-candidate-schema.json"},
		{Name: "BioSamples Checklist", File: "../../res/schemas/biosamples-schema.json"},
	}
)

func TestInterrogate(t *testing.T) {
	tests := []struct {
		documentFile           string
		expectedCandidateNames []string
	}{
		{
			documentFile:           "../../res/json/ncbi-SAMN03894263.json",
			expectedCandidateNames: []string{"NCBI Candidate Checklist"},
		},
	}
	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		interrogator.NewInterrogator(
			log.New(os.Stdout, "TestInterrogate ", log.LstdFlags|log.Lshortfile),
			&validator.Validator{},
			sampleCreated,
			sampleInterrogated,
			checklists,
		)

		sample := model.Sample{UUID: "test-uuid", Document: string(document)}
		sampleCreated <- sample
		ir := <-sampleInterrogated
		var candidateNames []string
		for _, checklist := range ir.CandidateChecklists {
			candidateNames = append(candidateNames, checklist.Name)
		}
		assert.Equal(t, test.expectedCandidateNames, candidateNames)
	}
}
