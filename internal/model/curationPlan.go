package model

import (
	"fmt"
	"github.com/tidwall/sjson"
	"log"
)

//CurationPlan is a series of curations to attempt to take a sample complying to one checklist to another
type CurationPlan struct {
	Logger        *log.Logger
	Name          string
	FromChecklist Checklist
	ToChecklist   Checklist
	Curations     []Curation
}

func (cp *CurationPlan) Execute(s Sample) Sample {
	cp.Logger.Printf("running curation plan '%s' on sample %s", cp.Name, s.UUID)
	for _, c := range cp.Curations {
		s = cp.applyCuration(s, c)
	}
	return s
}

func (cp *CurationPlan) applyCuration(sample Sample, curation Curation) Sample {
	curatedDocument, _ := sjson.Set(sample.Document, fmt.Sprintf("characteristics.%s.0.text", curation.Characteristic), curation.NewValue)
	return Sample{
		UUID:     sample.UUID,
		Document: string(curatedDocument),
	}
}
