.PHONY: build
build:
	go get -d ./...
	go build -v ./cmd/caribo

.DEFAULT_GOAL := build