.DEFAULT_GOAL := build

STATICCHECK_VERSION ?= v0.8.1

.PHONY: install-tools format lint comments test-build test build sync-corpus

sync-corpus:
	cp corpus/rules.yaml ../wlmarkdown-swift/Sources/WLMarkdown/Resources/
	cp corpus/dialect.yaml corpus/plain_text.yaml ../wlmarkdown-swift/Tests/WLMarkdownTests/Resources/
	cp corpus/rules.yaml ../wlmarkdown-kotlin/src/main/resources/
	cp corpus/dialect.yaml corpus/plain_text.yaml ../wlmarkdown-kotlin/src/test/resources/

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	python3 -m pip install --quiet --upgrade git+https://github.com/botforge-pro/commentcensor.git

format:
	gofmt -w .

comments:
	commentcensor *.go

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
