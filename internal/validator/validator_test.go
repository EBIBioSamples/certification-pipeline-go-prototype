package validator_test

import (
	"curation-pipeline/internal/model"
	"curation-pipeline/internal/validator"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestValidate(t *testing.T) {
	checklist := model.Checklist{
		Name: "NCBI Schema",
		URL:  "res/schema/ncbi-schema.json",
	}

	b, err := ioutil.ReadFile("./res/json/ncbi-SAMN03894263.json")
	if err != nil {
		fmt.Println(err)
	}
	str := string(b)
	validator := validator.Validator{}

	validator.Validate(checklist, str)
}
