package certifier

import (
	"crypto/md5"
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/validator"
	"github.com/pkg/errors"
	"hash"
	"io/ioutil"
	"log"
)

//Certifier validates a sample against a checklists and issues a certificate if successful
type Certifier struct {
	logger            *log.Logger
	certificateIssued chan model.Certificate
	checklists        []model.Checklist
	hash              hash.Hash
}

func (c *Certifier) certify(cpr model.PlanResult) {
	for _, checklist := range c.checklists {
		schema, err := ioutil.ReadFile(checklist.File)
		if err != nil {
			c.logger.Panic(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := validator.Validate(string(schema), cpr.Sample.Document)
		if err != nil {
			c.logger.Panic(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}

		if vr.Valid {
			cert := model.Certificate{
				Sample:        cpr.Sample,
				SampleHash:    fmt.Sprintf("%x", md5.Sum([]byte(cpr.Sample.Document))),
				Checklist:     checklist,
				ChecklistHash: fmt.Sprintf("%x", md5.Sum(schema)),
			}
			c.certificateIssued <- cert
		}
	}
}

//NewCertifier returns a new instance of an Certifier with the specified checklists
func NewCertifier(logger *log.Logger, in chan model.PlanResult, checklists []model.Checklist) chan model.Certificate {
	c := Certifier{
		logger:            logger,
		certificateIssued: make(chan model.Certificate),
		checklists:        checklists,
		hash:              md5.New(),
	}
	c.handleEvents(in)
	return c.certificateIssued
}

func (c *Certifier) handleEvents(in chan model.PlanResult) {
	go func() {
		for {
			select {
			case cpr := <-in:
				c.certify(cpr)
			}
		}
	}()
}
