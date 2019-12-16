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
INTERNAL_NAOS_SERVER_GEN_FILES = base_handlers.go
build: clean
	mv pkg/data/gen/service*.go pkg/data/
	$(GOGEN) $(REPO_NAME)/pkg/data
	mv pkg/data/service*.go pkg/data/gen/

	@for file in $(INTERNAL_NAOS_SERVER_GEN_FILES) ; do \
		mv internal/naos/server/gen/$$file internal/naos/server/$$file; \
		$(GOGEN) $(REPO_NAME)/internal/naos/server; \
		mv internal/naos/server/$$file internal/naos/server/gen/$$file; \
	done

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done
