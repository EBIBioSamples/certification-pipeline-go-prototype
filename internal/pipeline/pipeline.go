package pipeline

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/certifier"
	"github.com/EBIBioSamples/certification-pipeline/internal/config"
	"github.com/EBIBioSamples/certification-pipeline/internal/creator"
	"github.com/EBIBioSamples/certification-pipeline/internal/curator"
	"github.com/EBIBioSamples/certification-pipeline/internal/interrogator"
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"github.com/EBIBioSamples/certification-pipeline/internal/recorder"
	"log"
	"strings"
)

var (
	creatorIn               = make(chan string)
	interrogatorIn          = make(chan model.Sample)
	curatorIn               = make(chan model.ChecklistMatches)
	certifierIn             = make(chan model.PlanResult)
	recordSampleCreated     = make(chan model.Sample)
	recordPlanCompleted     = make(chan model.PlanResult)
	recordCertificateIssued = make(chan model.Certificate)
)

type Pipeline struct {
	logger             *log.Logger
	sampleCreated      chan model.Sample
	sampleInterrogated chan model.ChecklistMatches
	planCompleted      chan model.PlanResult
	curationCompeted   chan model.CurationEnd
	certificateIssued  chan model.Certificate
}

func NewPipeline(c *config.Config, in chan string) *Pipeline {
	p := Pipeline{
		logger:             c.Logger,
		sampleCreated:      creator.NewCreator(creatorIn),
		sampleInterrogated: interrogator.NewInterrogator(c.Logger, interrogatorIn, c.Checklists),
		certificateIssued:  certifier.NewCertifier(c.Logger, certifierIn, c.Checklists),
	}
	p.planCompleted, p.curationCompeted = curator.NewCurator(c.Logger, curatorIn, c.Plans)
	recorder.NewRecorder(c.Logger, recordSampleCreated, recordPlanCompleted, recordCertificateIssued)
	p.handleEvents(in, p.sampleCreated, p.sampleInterrogated, p.planCompleted, p.certificateIssued)
	return &p
}

func (p *Pipeline) handleEvents(
	in chan string, sampleCreated chan model.Sample, sampleInterrogated chan model.ChecklistMatches, planCompleted chan model.PlanResult, certificateIssued chan model.Certificate) {
	go func() {
		for {
			select {
			case input := <-in:
				p.onIn(input)
			case s := <-sampleCreated:
				p.onSampleCreated(s)
			case cm := <-sampleInterrogated:
				p.onSampleInterrogated(cm)
			case pr := <-planCompleted:
				p.onPlanCompleted(pr)
			case c := <-certificateIssued:
				p.onCertificateIssued(c)
			}
		}
	}()
}

func (p *Pipeline) onIn(input string) {
	p.logger.Println()
	p.logger.Printf("Input\t\t\t\t\t\t | len:%v", len(input))
	creatorIn <- input
}

func (p *Pipeline) onSampleCreated(sample model.Sample) {
	p.logger.Printf("Sample Created\t\t\t\t | sample:%s", sample.UUID)
	interrogatorIn <- sample
	recordSampleCreated <- sample
}

func (p *Pipeline) onSampleInterrogated(cm model.ChecklistMatches) {
	var ids []string
	for _, c := range cm.Checklists {
		ids = append(ids, c.ID())
	}
	p.logger.Printf("Sample Interrograted\t\t | sample:%s matched:%s", cm.Sample.UUID, strings.Join(ids, ", "))
	curatorIn <- cm
}

func (p *Pipeline) onPlanCompleted(pr model.PlanResult) {
	p.logger.Printf("Plan Completed\t\t\t\t | sample:%s plan:%s ", pr.Sample.UUID, pr.Plan.Describe())
	certifierIn <- pr
	recordPlanCompleted <- pr
}

func (p *Pipeline) onCertificateIssued(c model.Certificate) {
	p.logger.Printf("Certificate Issued\t\t\t | sample:%s certificate:%s", c.Sample.UUID, c.Checklist.ID())
	cm := model.ChecklistMatches{
		Sample:     c.Sample,
		Checklists: []model.Checklist{c.Checklist},
	}
	curatorIn <- cm
	recordCertificateIssued <- c
}
