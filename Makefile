.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/ddcat

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run -E gofmt,misspell
