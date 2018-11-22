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

func NewRecorder(logger *log.Logger, sampleCreated chan model.Sample, planCompleted chan model.PlanResult, certificateIssued chan model.Certificate) *Recorder {
	r := Recorder{
		logger: logger,
	}
	r.handleEvents(sampleCreated, planCompleted, certificateIssued)
	return &r
}

func (r *Recorder) handleEvents(
	sampleCreated chan model.Sample,
	planCompleted chan model.PlanResult,
	certificateIssued chan model.Certificate) {
	go func() {
		for {
			select {
			case s := <-sampleCreated:
				r.onSampleCreated(s)
			case pr := <-planCompleted:
				r.onPlanCompleted(pr)
			case c := <-certificateIssued:
				r.onCertificateIssued(c)
			}
		}
	}()
}

func (r *Recorder) onSampleCreated(sample model.Sample) {
	r.logger.Printf("Recorded Sample Created\t | sample:%s", sample.UUID)
}

func (r *Recorder) onPlanCompleted(pr model.PlanResult) {
	r.logger.Printf("Recorded Plan Completed\t | sample:%s plan:%s ", pr.Sample.UUID, pr.Plan.Describe())
}

func (r *Recorder) onCertificateIssued(c model.Certificate) {
	r.logger.Printf("Recorded Certificate Issued\t | sample:%s certificate:%s", c.Sample.UUID, c.Checklist.ID())
	err := ioutil.WriteFile(fmt.Sprintf("%s-%s.json", c.Sample.UUID, c.Checklist.ID()), []byte(c.Sample.Document), 0644)
	if err != nil {
		r.logger.Printf("Error writing file: %s", err.Error())
	}
}
