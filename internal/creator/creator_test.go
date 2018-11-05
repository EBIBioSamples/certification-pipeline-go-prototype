package creator_test

import (
	"fmt"
	"github.com/EBIBioSamples/curation-pipeline/internal/creator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
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

		c := creator.Creator{
			Logger: log.New(os.Stdout, "TestCreateSample ", log.LstdFlags|log.Lshortfile),
		}
		sample := c.CreateSample(string(document))
		assert.Regexp(t, `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, sample.UUID)
	}
}
