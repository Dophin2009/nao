# Golang variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Project variables
BINARY_NAME=goni

all: build

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)