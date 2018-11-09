package creator_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/creator"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	sampleCreated = make(chan model.Sample)
)

func TestCreateSample(t *testing.T) {
	go func(sampleCreated chan model.Sample) {
		for {
			select {
			case sample := <-sampleCreated:
				fmt.Printf("sample created: %s\n", sample.UUID)
			}
		}
	}(sampleCreated)
	tests := []struct {
		documentFile string
	}{
		{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
		},
	}
	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		c := creator.NewCreator(
			log.New(os.Stdout, "TestCreateSample ", log.LstdFlags|log.Lshortfile),
			sampleCreated)
		sample := c.CreateSample(string(document))
		assert.Regexp(t, `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, sample.UUID)
	}
}
