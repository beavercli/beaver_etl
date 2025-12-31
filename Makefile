
.PHONY: help run build 

.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

## run: Build and run the server
run: build
	./bin/etl

## build: Build the server binary
build:
	go build -o bin/etl ./cmd/etl

## test: Run all tests in verbose mode
test:
	go test -v ./...

