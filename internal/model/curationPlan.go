package model

import (
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

func (c *CurationPlan) Execute(s Sample) Sample {
	c.Logger.Printf("running curation plan %s on sample %s", c.Name, s.UUID)
	return s
}
