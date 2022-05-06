package repositories

import "github.com/german-muzquiz/factory-crd/pkg/domain"

type InMemoryFactoryRepository struct {
	Factories []domain.Factory
}

func (r *InMemoryFactoryRepository) GetFactories() []domain.Factory {
	return r.Factories
}
