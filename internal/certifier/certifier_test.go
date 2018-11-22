package certifier_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/certifier"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	curationPlanCompleted = make(chan model.PlanResult)
	certificateIssued     = make(chan model.Certificate)
	checklists            = []model.Checklist{
		{Name: "BioSamples Checklist", File: "../../res/schemas/biosamples-schema.json"},
	}
)

func TestCertify(t *testing.T) {
	tests := []struct {
		documentFile string
	}{
		{
			documentFile: "../../res/json/ncbi-SAMN03894263-curated.json",
		},
	}
	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		certificateIssued := certifier.NewCertifier(
			log.New(os.Stdout, "TestInterrogate ", log.LstdFlags|log.Lshortfile),
			curationPlanCompleted,
			checklists,
		)

		sample := model.Sample{Accession: "test-uuid", Document: string(document)}
		cpr := model.PlanResult{
			Sample: sample,
		}
		curationPlanCompleted <- cpr
		c := <-certificateIssued
		assert.Equal(t, sample.Document, c.Sample.Document)
		assert.Equal(t, checklists[0].Name, c.Checklist.Name)
	}
}
