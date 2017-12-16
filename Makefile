GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOLINT=golint
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=go-wc

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
deps-test-all: golint
lint:
	$(GOLINT) -min_confidence 0.1 -set_exit_status
golint:
	$(GOGET) github.com/golang/lint/golint
regression:
	./regression.sh
test:
	$(GOTEST) -v ./...
test-all: deps-test-all vet lint test regression
vet:
	$(GOVET)
