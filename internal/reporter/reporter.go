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

type Result struct {
	Cert  model.Certificate
	Badge string
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
				r.certMap[cert.Sample.Accession] = cert
				r.Unlock()
				r.logger.Printf("recorded %s certificate for %s", cert.Checklist.Name, cert.Sample.Accession)
			}
		}
	}()
}

func (r *Reporter) SampleInfo(uuid string) (result Result, ok bool) {
	r.RLock()
	cert, ok := r.certMap[uuid]
	var res Result
	if ok {
		res = Result{Cert: cert, Badge: cert.Badge()}
	}
	r.RUnlock()
	return res, ok
}
