package adapters

import (
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/domain"
)

type InMemoryFactoryRepository struct {
	Factories map[string]*domain.Factory
}

func (r *InMemoryFactoryRepository) GetFactories() map[string]domain.Factory {
	result := map[string]domain.Factory{}
	for n := range r.Factories {
		result[n] = *r.Factories[n]
	}
	return result
}

func (r *InMemoryFactoryRepository) UpdateCapacity(name string, newCapacity int) {
	r.Factories[name].Status.CurrentCapacity = fmt.Sprintf("%d vehicles per minute", newCapacity)
}
