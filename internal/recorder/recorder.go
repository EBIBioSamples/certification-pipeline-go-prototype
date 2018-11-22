package recorder

import (
	"fmt"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"io/ioutil"
	"log"
)

//Recorder records the actions performed on a sample to a file
type Recorder struct {
	logger *log.Logger
}

func NewRecorder(logger *log.Logger, sampleIdentified chan model.Sample, planCompleted chan model.PlanResult, certificateIssued chan model.Certificate) *Recorder {
	r := Recorder{
		logger: logger,
	}
	r.handleEvents(sampleIdentified, planCompleted, certificateIssued)
	return &r
}

func (r *Recorder) handleEvents(
	sampleIdentified chan model.Sample,
	planCompleted chan model.PlanResult,
	certificateIssued chan model.Certificate) {
	go func() {
		for {
			select {
			case s := <-sampleIdentified:
				r.onSampleIdentified(s)
			case pr := <-planCompleted:
				r.onPlanCompleted(pr)
			case c := <-certificateIssued:
				r.onCertificateIssued(c)
			}
		}
	}()
}

func (r *Recorder) onSampleIdentified(sample model.Sample) {
	r.logger.Printf("Recorded Sample Identified\t | sample:%s", sample.Accession)
}

func (r *Recorder) onPlanCompleted(pr model.PlanResult) {
	r.logger.Printf("Recorded Plan Completed\t | sample:%s plan:%s ", pr.Sample.Accession, pr.Plan.Describe())
}

func (r *Recorder) onCertificateIssued(c model.Certificate) {
	r.logger.Printf("Recorded Certificate Issued\t | sample:%s certificate:%s", c.Sample.Accession, c.Checklist.ID())
	err := ioutil.WriteFile(fmt.Sprintf("%s-%s.json", c.Sample.Accession, c.Checklist.ID()), []byte(c.Sample.Document), 0644)
	if err != nil {
		r.logger.Printf("Error writing file: %s", err.Error())
	}
}
