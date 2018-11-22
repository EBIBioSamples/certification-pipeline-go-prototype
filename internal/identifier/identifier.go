package identifier

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

//Identifier registers a samples document with the system to enable tracking
type Identifier struct {
	sampleIdentified chan model.Sample
}

//CreateSample assigns a Accession to a sample document to use for tracking
func (c *Identifier) identify(json string) {
	var accession string
	value := gjson.Get(json, "accession")
	if !value.Exists() {
		accession = uuid.Must(uuid.NewUUID()).String()
	}
	accession = value.String()
	sample := model.Sample{
		Accession: accession,
		Document:  json,
	}
	c.sampleIdentified <- sample
}

//NewIdentifier returns a new instance of an identifier with an output channel
func NewIdentifier(in chan string) chan model.Sample {
	c := Identifier{
		sampleIdentified: make(chan model.Sample),
	}
	c.handleEvents(in)
	return c.sampleIdentified
}

func (c *Identifier) handleEvents(in chan string) {
	go func() {
		for {
			select {
			case input := <-in:
				c.identify(input)
			}
		}
	}()
}
