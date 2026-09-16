.DEFAULT_GOAL := build

STATICCHECK_VERSION ?= v0.8.0
COMMENTCENSOR_VERSION ?= v0.3.1

.PHONY: install-tools format lint comments test-build test build

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	python3 -m pip install --quiet --upgrade git+https://github.com/botforge-pro/commentcensor.git@$(COMMENTCENSOR_VERSION)

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
