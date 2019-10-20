# Golang variables
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

# Project variables
TARGET_DIR=target
REPO_NAME=gitlab.com/Dophin2009/anisheet
MODULES=anisheet

default: build

clean:
	rm -rf $(TARGET_DIR)/

build: clean
	mv pkg/data/gen/service.go pkg/data/service.go
	$(GOGEN) $(REPO_NAME)/pkg/data
	mv pkg/data/service.go pkg/data/gen/service.go

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done