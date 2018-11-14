package curator_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/curator"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	logger             = log.New(os.Stdout, "TestCurate ", log.LstdFlags|log.Lshortfile)
	sampleInterrogated = make(chan model.InterrogationResult)
	planCompleted      = make(chan model.PlanResult)
	certificateIssued  = make(chan model.Certificate)
	checklistMap       = make(map[string]model.Checklist)
	c                  = config.NewConfig(logger, "../../res/config.json")
)

func init() {
	for _, checklist := range c.Checklists {
		checklistMap[checklist.ID()] = checklist
	}
}

func TestCurate(t *testing.T) {
	tests := []struct {
		documentFile        string
		curatedDocumentFile string
	}{
		{
			documentFile:        "../../res/json/ncbi-SAMN03894263.json",
			curatedDocumentFile: "../../res/json/ncbi-SAMN03894263-curated.json",
		},
	}
	for _, test := range tests {
		curator.NewCurator(
			logger,
			sampleInterrogated,
			planCompleted,
			certificateIssued,
			c.Plans,
		)
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}
		sample := model.Sample{UUID: "test-uuid", Document: string(document)}
		sampleInterrogated <- model.InterrogationResult{
			Sample:              sample,
			CandidateChecklists: []model.Checklist{checklistMap["ncbi-0.0.1"]},
		}
		pr := <-planCompleted
		curatedDocument, err := ioutil.ReadFile(test.curatedDocumentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.curatedDocumentFile)))
		}
		assert.Equal(t, string(curatedDocument), pr.Sample.Document)
	}
}
