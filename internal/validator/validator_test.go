package validator_test

import (
	"curation-pipeline/internal/validator"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	data, err := os.Open("./res/json/ncbi-SAMN03894263.json")
	if err != nil {
		log.Fatal(errors.Wrap(err, "read failed"))
	}

	checklist, err := os.Open("./res/schemas/ncbi-schema.json")
	if err != nil {
		log.Fatal(errors.Wrap(err, "read failed"))
	}

	validator := validator.Validator{}

	result, _ := validator.Validate(checklist, data)
	fmt.Println(result)
}
