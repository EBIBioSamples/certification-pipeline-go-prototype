package model

import (
	"fmt"
	"github.com/tidwall/sjson"
)

//Plan is a series of curations to attempt to take a sample complying to one checklist to another
type Plan struct {
	CandidateChecklistID   string     `json:"candidate_checklist_id"`
	CertificateChecklistID string     `json:"certification_checklist_id"`
	Curations              []Curation `json:"curations"`
}

func (p *Plan) Execute(s Sample) Sample {
	for _, c := range p.Curations {
		s = p.applyCuration(s, c)
	}
	return s
}

func (p *Plan) applyCuration(sample Sample, curation Curation) Sample {
	curatedDocument, _ := sjson.Set(sample.Document, fmt.Sprintf("characteristics.%s.0.text", curation.Characteristic), curation.Value)
	return Sample{
		UUID:     sample.UUID,
		Document: string(curatedDocument),
	}
}

func (p *Plan) Describe() string {
	return fmt.Sprintf("%s to %s", p.CandidateChecklistID, p.CertificateChecklistID)
}
