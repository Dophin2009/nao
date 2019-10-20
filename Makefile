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

# Fix this
build: clean
	mv pkg/data/gen/service.go pkg/data/service.go
	$(GOGEN) $(REPO_NAME)/pkg/data
	mv pkg/data/service.go pkg/data/gen/service.go

	mv cmd/anisheet/controller/gen/routers.go cmd/anisheet/controller/routers.go
	$(GOGEN) $(REPO_NAME)/cmd/anisheet/controller
	mv cmd/anisheet/controller/routers.go cmd/anisheet/controller/gen/routers.go

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done