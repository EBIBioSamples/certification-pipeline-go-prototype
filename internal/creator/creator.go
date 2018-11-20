package creator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/google/uuid"
	"log"
)

//Creator registers a samples document with the system to enable tracking
type Creator struct {
	logger        *log.Logger
	jsonSubmitted chan string
	sampleCreated chan model.Sample
}

//CreateSample assigns a UUID to a sample document to use for tracking
func (c *Creator) createSample(json string) {
	sample := model.Sample{
		UUID:     uuid.Must(uuid.NewUUID()).String(),
		Document: json,
	}
	c.logger.Printf("created new sample: %s", sample.UUID)
	c.sampleCreated <- sample
}

//NewCreator returns a new instance of a creator with an output channel
func NewCreator(logger *log.Logger, jsonSubmitted chan string, sampleCreated chan model.Sample) *Creator {
	c := Creator{
		logger:        logger,
		jsonSubmitted: jsonSubmitted,
		sampleCreated: sampleCreated,
	}
	c.handleEvents(jsonSubmitted)
	return &c
}

func (c *Creator) handleEvents(jsonSubmitted chan string) {
	go func() {
		for {
			select {
			case j := <-jsonSubmitted:
				c.createSample(j)
			}
		}
	}()
}
