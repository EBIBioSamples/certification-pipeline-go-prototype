package validator

import "curation-pipeline/internal/model"

type Validator struct {
}

func (v *Validator) Validate(checklist model.Checklist, json string) error {
	return nil
}
