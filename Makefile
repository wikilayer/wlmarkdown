.DEFAULT_GOAL := build

.PHONY: format lint comments test-build test build

format:
	gofmt -w .

comments:
	commentcensor .

lint: comments
	go vet ./...
	gofmt -l . | (! grep .)
	staticcheck ./...

test-build:
	go build ./...
	go test -run '^$$' ./...

test:
	go test ./...

build: lint test-build test
	go build ./...
