package repositories

import "github.com/german-muzquiz/factory-crd/pkg/domain"

type FactoryRepository interface {
	GetFactories() []domain.Factory
}
