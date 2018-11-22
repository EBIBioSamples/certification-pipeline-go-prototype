package creator

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

//Creator registers a samples document with the system to enable tracking
type Creator struct {
	sampleCreated chan model.Sample
}

//CreateSample assigns a UUID to a sample document to use for tracking
func (c *Creator) createSample(json string) {
	sample := model.Sample{
		UUID:     identify(json),
		Document: json,
	}
	c.sampleCreated <- sample
}

func identify(json string) string {
	value := gjson.Get(json, "accession")
	if !value.Exists() {
		return uuid.Must(uuid.NewUUID()).String()
	}
	return value.String()
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
