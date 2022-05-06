package domain

import "k8s.io/apimachinery/pkg/runtime"

type Factory struct {
	runtime.Object `json:"-"`
	Spec           FactorySpec `json:"spec"`
}

type FactorySpec struct {
	Location string `json:"location"`
}
