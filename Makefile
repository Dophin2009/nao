# Golang variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Project variables
TARGET_DIR=bin
BINARY_NAME=goni

all: build

build:
	$(GOBUILD) -o $(TARGET_DIR)/$(BINARY_NAME) -v

run: build
	$(TARGET_DIR)/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(TARGET_DIR)/$(BINARY_NAME)