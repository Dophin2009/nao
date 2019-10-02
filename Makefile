# Golang variables
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Project variables
TARGET_DIR=bin
REPO_NAME=gitlab.com/Dophin2009/anisheet
MODULES=anisheet

default: build

build:
	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done

clean:
	rm -rf $(TARGET_DIR)/