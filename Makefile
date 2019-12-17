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
	# mv internal/data/gen/* internal/data/
	# $(GOGEN) $(REPO_NAME)/internal/data
	# mv internal/data/*_gen.go internal/data/gen/

	mv internal/naos/server/gen/* internal/naos/server/
	$(GOGEN) $(REPO_NAME)/internal/naos/server 
	mv internal/naos/server/*_gen.go internal/naos/server/gen/ 

	@for module in $(MODULES) ; do \
		$(GOBUILD) -o $(TARGET_DIR)/$$module -v $(REPO_NAME)/cmd/$$module ; \
	done
