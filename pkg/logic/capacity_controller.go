package logic

import (
	"github.com/german-muzquiz/factory-crd/pkg/ports"
	"math/rand"
	"time"
)

type CapacityController struct {
	FactoryRepository ports.FactoryRepository
}

func (cc *CapacityController) Init() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				cc.updateCapacity()
			}
		}
	}()
}

func (cc *CapacityController) updateCapacity() {
	for name := range cc.FactoryRepository.GetFactories() {
		c := rand.Intn(100)
		cc.FactoryRepository.UpdateCapacity(name, c)
	}
}
