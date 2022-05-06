
BINARY   := factory-control
OS       ?= $(shell go version | cut -d' ' -f 4 | cut -d'/' -f 1)
ARCH   	 ?= $(shell go version | cut -d' ' -f 4 | cut -d'/' -f 2)
REGISTRY := 192.168.64.2:32000/factory-control
VERSION  ?= latest


.PHONY: clean
clean:
	@rm -rf bin

.PHONY: build
build:
	@echo "Building: $(BINARY)"
	@go build -mod=vendor -o bin/${OS}_${ARCH}/${BINARY} cmd/main.go

.PHONY: docker
docker:
	@docker build -t ${REGISTRY}:${VERSION} .

.PHONY: docker-push
docker-push:
	@docker push ${REGISTRY}:${VERSION}
