package creator

import (
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/google/uuid"
	"log"
)

//Creator registers a samples document with the system to enable tracking
type Creator struct {
	logger        *log.Logger
	sampleCreated chan model.Sample
}

//CreateSample assigns a UUID to a sample document to use for tracking
func (c *Creator) CreateSample(json string) *model.Sample {
	sample := model.Sample{
		UUID:     uuid.Must(uuid.NewUUID()).String(),
		Document: json,
	}
	c.logger.Printf("created new sample: %s", sample.UUID)
	c.sampleCreated <- sample
	return &sample
}

//NewCreator returns a new instance of a creator with an output channel
func NewCreator(logger *log.Logger, sampleCreated chan model.Sample) *Creator {
	c := Creator{
		logger:        logger,
		sampleCreated: sampleCreated,
	}
	return &c
}
