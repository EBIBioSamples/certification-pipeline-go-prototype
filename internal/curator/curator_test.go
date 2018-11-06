package curator_test

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/curator"
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	logger                = log.New(os.Stdout, "TestCurate ", log.LstdFlags|log.Lshortfile)
	sampleInterrogated    = make(chan model.InterrogationResult)
	curationPlanCompleted = make(chan model.CurationPlanResult)
	checklists            = []model.Checklist{
		{Name: "NCBI Candidate Checklist", File: "../../res/schemas/ncbi-candidate-schema.json"},
		{Name: "BioSamples Checklist", File: "../../res/schemas/biosamples-schema.json"},
	}
	checklistMap  = make(map[string]model.Checklist)
	curationPlans []model.CurationPlan
)

func init() {
	for _, checklist := range checklists {
		checklistMap[checklist.Name] = checklist
	}
	curationPlans = []model.CurationPlan{
		{
			Logger:        logger,
			Name:          "NCBI to BioSamples",
			FromChecklist: checklistMap["NCBI Candidate Checklist"],
			ToChecklist:   checklistMap["BioSamples Checklist"],
			Curations: []model.Curation{
				{
					Characteristic: "INSDC status",
					NewValue:       "public",
				},
			},
		},
	}
}

func TestCurate(t *testing.T) {
	tests := []struct {
		documentFile string
	}{
		{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
		},
	}
	for _, test := range tests {
		curator.NewCurator(
			logger,
			sampleInterrogated,
			curationPlanCompleted,
			curationPlans,
		)
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}
		sample := model.Sample{UUID: "test-uuid", Document: string(document)}
		sampleInterrogated <- model.InterrogationResult{
			Sample:              sample,
			CandidateChecklists: []model.Checklist{checklistMap["NCBI Candidate Checklist"]},
		}
	}
}
