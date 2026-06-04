MODULE  := github.com/sratabix/cem/v3
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X $(MODULE)/cmd.version=$(VERSION)
BIN     := bin

.PHONY: all build test vet staticcheck vuln lint clean

all: lint test build

build:
	go build -ldflags='$(LDFLAGS)' -o $(BIN)/cem .

test:
	go test -race -coverprofile=coverage.out ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

vuln:
	govulncheck ./...

lint: vet staticcheck vuln

clean:
	rm -rf $(BIN) coverage.out

cross:
	GOOS=darwin  GOARCH=arm64 go build -ldflags='$(LDFLAGS)' -o $(BIN)/cem-darwin-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags='$(LDFLAGS)' -o $(BIN)/cem-darwin-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags='$(LDFLAGS)' -o $(BIN)/cem-linux-arm64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags='$(LDFLAGS)' -o $(BIN)/cem-linux-amd64 .
