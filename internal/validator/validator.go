package validator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/xeipuuv/gojsonschema"
)

//Validate a given document against a given schema
func Validate(schema string, document string) (model.ValidationResult, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return model.ValidationResult{}, err
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

	vr := model.ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  validationErrors,
	}

	return vr, nil
}
