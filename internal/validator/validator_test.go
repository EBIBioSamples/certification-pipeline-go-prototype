package validator_test

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/validator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		documentFile string
		schemaFile   string
		valid        bool
		errors       []string
		err          error
	}{
		{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
			schemaFile:   "../../res/schemas/ncbi-candidate-schema.json",
			valid:        true,
			err:          nil,
		},
		{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
			schemaFile:   "../../res/schemas/biosamples-schema.json",
			valid:        false,
			errors:       []string{`characteristics.INSDC status.0.text must be one of the following: "public"`},
			err:          nil,
		},
	}

	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		schema, err := ioutil.ReadFile(test.schemaFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.schemaFile)))
		}

		vr, err := validator.Validate(string(schema), string(document))

		assert.Equal(t, test.valid, vr.Valid)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.errors, vr.Errors)
	}
}
