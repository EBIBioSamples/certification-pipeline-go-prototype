package validator

import (
	"github.com/xeipuuv/gojsonschema"
	"io"
)

type Validator struct {
}

func (v *Validator) Validate(c io.Reader, d io.Reader) (string, error) {
	schemaLoader, _ := gojsonschema.NewReaderLoader(c)
	documentLoader, _ := gojsonschema.NewReaderLoader(d)
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
	fmt.println
	return message, nil
}
