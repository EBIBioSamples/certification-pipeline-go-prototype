package creator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/google/uuid"
)

//Creator registers a samples document with the system to enable tracking
type Creator struct {
	sampleCreated chan model.Sample
}

//CreateSample assigns a UUID to a sample document to use for tracking
func (c *Creator) createSample(json string) {
	sample := model.Sample{
		UUID:     uuid.Must(uuid.NewUUID()).String(),
		Document: json,
	}
	c.sampleCreated <- sample
}

//NewCreator returns a new instance of a creator with an output channel
func NewCreator(in chan string) chan model.Sample {
	c := Creator{
		sampleCreated: make(chan model.Sample),
	}
	c.handleEvents(in)
	return c.sampleCreated
}

func (c *Creator) handleEvents(in chan string) {
	go func() {
		for {
			select {
			case input := <-in:
				c.createSample(input)
			}
		}
	}()
}
