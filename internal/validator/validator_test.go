package validator_test

import (
	"curation-pipeline/internal/validator"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"testing"
)

func TestValidate(t *testing.T) {

	tests := []struct {
		documentFile string
		schemaFile   string
		result       string
		err          error
	}{
		{
			documentFile: "../../res/json/ncbi-SAMN03894263.json",
			schemaFile:   "../../res/schemas/ncbi-schema.json",
			err:          nil,
		}}

	for _, test := range tests {
		document, err := ioutil.ReadFile(test.documentFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.documentFile)))
		}

		schema, err := ioutil.ReadFile(test.schemaFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", test.schemaFile)))
		}

		validator := validator.Validator{}

		result, _ := validator.Validate(string(schema), string(document))
		fmt.Println(result)
	}
}
