package recorder

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
)

//Recorder records the actions performed on a sample to a file
type Recorder struct {
	sampleCreated chan model.Sample
}

func NewRecorder(sampleCreated chan model.Sample) *Recorder {
	r := Recorder{
		sampleCreated: sampleCreated,
	}
	r.handleEvents(sampleCreated)
	return &r
}

func (r *Recorder) handleEvents(sampleCreated chan model.Sample) {
	go func() {
		for {
			select {
			case s := <-sampleCreated:
				r.onSampleCreated(s)
			}
		}
	}()
}

func (r *Recorder) onSampleCreated(sample model.Sample) {
	//fmt.Printf("create sample: %s\n", sample.Document)
}
