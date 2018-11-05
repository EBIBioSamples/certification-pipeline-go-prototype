package creator

import (
	"github.com/EBIBioSamples/curation-pipeline/internal/model"
	"github.com/google/uuid"
	"log"
)

type Creator struct {
	Logger *log.Logger
}

func (c *Creator) CreateSample(json string) *model.Sample {
	sample := model.Sample{
		UUID:     uuid.Must(uuid.NewUUID()).String(),
		Document: json,
	}
	c.Logger.Printf("created new sample: %s", sample.UUID)
	return &sample
}
