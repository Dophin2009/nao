# Golang variables
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOGEN=$(GOCMD) generate

# Project variables
TARGET_DIR=bin
REPO_NAME=gitlab.com/Dophin2009/nao
MODULES=naos

SRC_FILES=find . -name '*.go' ! -name '*.gen.go'

.PHONY: check nakedret nargs

default: check

clean:
	rm -rf $(TARGET_DIR)/
	find . -type f -name '*.gen.go' -delete

build: generate test
	$(foreach module,$(MODULES),$(GOBUILD) -o $(TARGET_DIR)/$(module) -v $(REPO_NAME)/cmd/$(module))

generate: clean
	$(GORUN) scripts/gqlgen.go --verbose

test:
	$(GOTEST) ./...

check:
	$(GORUN) github.com/alexkohler/nakedret -l 0 $$($(SRC_FILES))
	$(GORUN) github.com/alexkohler/nargs/cmd/nargs $$($(SRC_FILES))

