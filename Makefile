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
	find . -type f -name '*.gen.go' -delete

# Fix this
build: clean
	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done
