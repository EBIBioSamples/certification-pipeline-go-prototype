package identifier_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/identifier"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

var (
	in = make(chan string)
)

func TestCreateSample(t *testing.T) {
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

		sampleIdentified := identifier.NewIdentifier(in)

		in <- string(document)
		sample := <-sampleIdentified

		assert.Regexp(t, `SAM[END][AG]?[0-9]+`, sample.Accession)
	}
}
