package pipeline_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/pipeline"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestPipeline(t *testing.T) {
	tests := []struct {
		documentFile string
	}{
		/*{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
		},*/
		{
			documentFile: "../../res/json/SAMEA3774859.json",
		},
	}
	for _, test := range tests {
		jsonSubmitted := make(chan string)
		logger := log.New(os.Stdout, "TestPipeline", log.LstdFlags|log.Lshortfile)

		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			logger.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		c, err := config.NewConfig(logger, "../../res/config_test.json", "../../res/schemas/config-schema.json")
		if err != nil {
			logger.Fatal(errors.Wrap(err, fmt.Sprintf("failed to create config")))
		}

		piplineFinished := pipeline.NewPipeline(c, jsonSubmitted)

		jsonSubmitted <- string(document)
		<-piplineFinished
	}

}
