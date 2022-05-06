package ports

import "github.com/german-muzquiz/factory-crd/pkg/domain"

type FactoryRepository interface {
	GetFactories() map[string]domain.Factory
	UpdateCapacity(name string, newCapacity int)
}
