package model_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	logger = log.New(os.Stdout, "TestPlan ", log.LstdFlags|log.Lshortfile)
)

func TestPlan(t *testing.T) {
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
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}
		curatedDocument, err := ioutil.ReadFile(test.curatedDocumentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.curatedDocumentFile)))
		}

		cp := model.Plan{
			Curations: []model.Curation{
				{
					Characteristic: "INSDC status",
					Value:          "public",
				},
			},
		}

		sample := model.Sample{UUID: "test-uuid", Document: string(document)}
		curatedSample := cp.Execute(sample)
		assert.NotEqual(t, sample.Document, curatedSample.Document)
		assert.Equal(t, string(curatedDocument), curatedSample.Document)
	}
}
