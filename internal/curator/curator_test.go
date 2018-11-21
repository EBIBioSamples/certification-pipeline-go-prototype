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
	logger       = log.New(os.Stdout, "TestCurate ", log.LstdFlags|log.Lshortfile)
	checklistMap = make(map[string]model.Checklist)
	c, _         = config.NewConfig(logger, "../../res/config.json", "../../res/schemas/config-schema.json")
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
		checklistMatched := make(chan model.ChecklistMatches)
		planCompleted, _ := curator.NewCurator(
			logger,
			checklistMatched,
			c.Plans,
		)
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}
		sample := model.Sample{UUID: "test-uuid", Document: string(document)}
		checklistMatched <- model.ChecklistMatches{
			Sample:     sample,
			Checklists: []model.Checklist{checklistMap["ncbi-0.0.1"]},
		}
		pr := <-planCompleted
		curatedDocument, err := ioutil.ReadFile(test.curatedDocumentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.curatedDocumentFile)))
		}
		assert.Equal(t, string(curatedDocument), pr.Sample.Document)
	}
}

func TestCurateWithNoMatchingPlans(t *testing.T) {
	checklistMatched := make(chan model.ChecklistMatches)
	_, curationCompleted := curator.NewCurator(
		logger,
		checklistMatched,
		c.Plans,
	)
	sample := model.Sample{UUID: "test-uuid", Document: "{}"}
	checklistMatched <- model.ChecklistMatches{
		Sample:     sample,
		Checklists: []model.Checklist{checklistMap["biosamples-0.0.1"]},
	}
	cc := <-curationCompleted
	assert.Equal(t, "biosamples-0.0.1", cc.Checklist.ID())
}
