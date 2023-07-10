# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOFLAGS := -v 
LDFLAGS := -s -w

ifneq ($(shell go env GOOS),darwin)
LDFLAGS := -extldflags "-static"
endif
    
all: build
build:
	$(GOBUILD) $(GOFLAGS) -ldflags '$(LDFLAGS)' -o "evm-cli" main.go
test: 
	$(GOTEST) $(GOFLAGS) ./... -cover -coverprofile=cover.out
cover: test
	$(GOCMD) tool cover -html=cover.out
tidy:
	$(GOMOD) tidy