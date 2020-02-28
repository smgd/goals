.PHONY: build
build:
	go get -d ./...
	go build -v ./cmd/caribou

.DEFAULT_GOAL := build