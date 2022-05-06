package main

import (
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/adapters"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := &adapters.KubeFactoryRepository{}
	r.Init()
	ra := &adapters.RestAdapter{FactoryRepository: r}

	http.HandleFunc("/factories", ra.GetFactories())

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
