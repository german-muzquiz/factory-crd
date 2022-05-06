package domain

import "k8s.io/apimachinery/pkg/runtime"

type Factory struct {
	runtime.Object `json:"-"`
	Spec           FactorySpec   `json:"spec"`
	Status         FactoryStatus `json:"status"`
}

type FactorySpec struct {
	Location string `json:"location"`
}

type FactoryStatus struct {
	CurrentCapacity string `json:"currentCapacity"`
}
