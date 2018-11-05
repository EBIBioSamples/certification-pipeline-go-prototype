package validator

import (
	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
}

func (v *Validator) Validate(schema string, document string) (string, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err.Error(), err
	}

	var message string
	var validationErrors []string
	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.Description())
		}
	}
	return message, nil
}
