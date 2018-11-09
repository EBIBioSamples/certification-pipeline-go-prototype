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
	logger                *log.Logger
	validator             *validator.Validator
	curationPlanCompleted chan model.CurationPlanResult
	certificateIssued     chan model.Certificate
	checklists            []model.Checklist
	hash                  hash.Hash
}

func (c *Certifier) certify(cpr model.CurationPlanResult) {
	for _, checklist := range c.checklists {
		c.logger.Printf("validating %s against %s\n", cpr.Sample.UUID, checklist.Name)
		schema, err := ioutil.ReadFile(checklist.File)
		if err != nil {
			c.logger.Fatal(errors.Wrap(err, fmt.Sprintf("read failed for: %s", checklist)))
		}
		vr, err := c.validator.Validate(string(schema), cpr.Sample.Document)
		if err != nil {
			c.logger.Fatal(errors.Wrap(err, fmt.Sprintf("failed to validate")))
		}

		if vr.Valid {
			cert := model.Certificate{
				Sample:        cpr.Sample,
				SampleHash:    fmt.Sprintf("%x", md5.Sum([]byte(cpr.Sample.Document))),
				Checklist:     checklist,
				ChecklistHash: fmt.Sprintf("%x", md5.Sum(schema)),
			}
			c.logger.Printf("certificate for %s issued for sample: %s", checklist.Name, cpr.Sample.UUID)
			c.certificateIssued <- cert
		}
	}
}

//NewCertifier returns a new instance of an Certifier with the specified checklists
func NewCertifier(
	logger *log.Logger,
	validator *validator.Validator,
	curationPlanCompleted chan model.CurationPlanResult,
	certificateIssued chan model.Certificate,
	checklists []model.Checklist) *Certifier {
	c := Certifier{
		logger:                logger,
		validator:             validator,
		curationPlanCompleted: curationPlanCompleted,
		certificateIssued:     certificateIssued,
		checklists:            checklists,
		hash:                  md5.New(),
	}
	c.handleEvents(curationPlanCompleted)
	return &c
}

func (c *Certifier) handleEvents(curationPlanCompleted chan model.CurationPlanResult) {
	go func() {
		for {
			select {
			case cpr := <-curationPlanCompleted:
				c.certify(cpr)
			}
		}
	}()
}
