package reporter

import (
	"github.com/EBIBioSamples/certification-pipeline/internal/model"
	"log"
	"sync"
)

type Reporter struct {
	logger            *log.Logger
	certificateIssued chan model.Certificate
	sync.RWMutex
	certMap map[string]model.Certificate
}

func NewReporter(logger *log.Logger, certificateIssued chan model.Certificate) *Reporter {
	r := Reporter{
		logger:            logger,
		certificateIssued: certificateIssued,
		certMap:           make(map[string]model.Certificate),
	}
	r.handleEvents(certificateIssued)
	return &r
}

func (r *Reporter) handleEvents(certificateIssued chan model.Certificate) {
	go func() {
		for {
			select {
			case cert := <-certificateIssued:
				r.Lock()
				r.certMap[cert.Sample.UUID] = cert
				r.Unlock()
				r.logger.Printf("recorded %s certificate for %s", cert.Checklist.Name, cert.Sample.UUID)
			}
		}
	}()
}

func (r *Reporter) SampleInfo(uuid string) (cert model.Certificate, ok bool) {
	r.RLock()
	result, ok := r.certMap[uuid]
	r.RUnlock()
	return result, ok
}
