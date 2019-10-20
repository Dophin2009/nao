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
REPO_NAME=gitlab.com/Dophin2009/nao
MODULES=naos

default: build

clean:
	rm -rf $(TARGET_DIR)/

# Fix this
build: clean
	mv pkg/data/gen/service.go pkg/data/service.go
	$(GOGEN) $(REPO_NAME)/pkg/data
	mv pkg/data/service.go pkg/data/gen/service.go

	mv cmd/naos/controller/gen/routers.go cmd/naos/controller/routers.go
	$(GOGEN) $(REPO_NAME)/cmd/naos/controller
	mv cmd/naos/controller/routers.go cmd/naos/controller/gen/routers.go

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done