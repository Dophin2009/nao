# Golang variables
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

# Project variables
TARGET_DIR=bin
REPO_NAME=gitlab.com/Dophin2009/nao
MODULES=naos

default: build

clean:
	rm -rf $(TARGET_DIR)/

# Fix this
build: clean
	mv pkg/data/gen/service*.go pkg/data/
	$(GOGEN) $(REPO_NAME)/pkg/data
	mv pkg/data/service*.go pkg/data/gen/

	mv internal/naos/server/gen/base_handlers.go internal/naos/server/base_handlers.go
	$(GOGEN) $(REPO_NAME)/internal/naos/server
	mv internal/naos/server/base_handlers.go internal/naos/server/gen/base_handlers.go

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done
