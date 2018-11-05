package creator

import (
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/google/uuid"
	"log"
)

type Creator struct {
	logger        *log.Logger
	sampleCreated chan model.Sample
}

func (c *Creator) CreateSample(json string) *model.Sample {
	sample := model.Sample{
		UUID:     uuid.Must(uuid.NewUUID()).String(),
		Document: json,
	}
	c.logger.Printf("created new sample: %s", sample.UUID)
	c.sampleCreated <- sample
	return &sample
}

func NewCreator(logger *log.Logger, sampleCreated chan model.Sample) *Creator {
	c := Creator{
		logger:        logger,
		sampleCreated: sampleCreated,
	}
	return &c
}
