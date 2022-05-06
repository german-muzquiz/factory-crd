package main

import (
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/adapters"
	"github.com/german-muzquiz/factory-crd/pkg/domain"
	"github.com/german-muzquiz/factory-crd/pkg/repositories"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	f1 := domain.Factory{Name: "factory1"}
	r := &repositories.InMemoryFactoryRepository{Factories: []domain.Factory{f1}}
	ra := &adapters.RestAdapter{FactoryRepository: r}
	http.HandleFunc("/factories", ra.GetFactories())

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
