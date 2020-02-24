.PHONY: build
build:
	go build -v ./cmd/caribou

.DEFAULT_GOAL := build