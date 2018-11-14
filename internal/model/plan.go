package model

import (
	"fmt"
	"github.com/tidwall/sjson"
	"log"
)

//Plan is a series of curations to attempt to take a sample complying to one checklist to another
type Plan struct {
	Logger        *log.Logger
	Name          string
	FromChecklist Checklist
	ToChecklist   Checklist
	Curations     []Curation
}

func (p *Plan) Execute(s Sample) Sample {
	p.Logger.Printf("running plan '%s' on sample %s", p.Name, s.UUID)
	for _, c := range p.Curations {
		s = p.applyCuration(s, c)
	}
	return s
}

func (p *Plan) applyCuration(sample Sample, curation Curation) Sample {
	curatedDocument, _ := sjson.Set(sample.Document, fmt.Sprintf("characteristics.%s.0.text", curation.Characteristic), curation.NewValue)
	return Sample{
		UUID:     sample.UUID,
		Document: string(curatedDocument),
	}
}
