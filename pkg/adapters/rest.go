package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/ports"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RestAdapter struct {
	FactoryRepository ports.FactoryRepository
}

func (a *RestAdapter) GetFactories() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f := a.FactoryRepository.GetFactories()
		b, err := json.Marshal(f)
		if err != nil {
			log.WithError(err).Errorf("Error serializing factories")
			return
		}
		_, err = fmt.Fprintf(w, string(b))
		if err != nil {
			log.WithError(err).Errorf("Error printing response")
		}
	}
}
