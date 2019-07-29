VERSION=$(shell git tag -l --points-at HEAD)
GIT_COMMIT=$(shell git log --format="%H" -n 1)
DATE=$(shell date -u)

# Project name
PROJECT = tusd

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: clean all fmt vet lint build test install static deps docker
.DEFAULT: default

all: clean build fmt lint test vet install

build:
	@echo "+ $@"
	go build -ldflags="-X github.com/tus/tusd/cmd/tusd/cli.VersionName=$(VERSION) -X github.com/tus/tusd/cmd/tusd/cli.GitCommit=$(GIT_COMMIT) -X 'github.com/tus/tusd/cmd/tusd/cli.BuildDate=$(DATE)'" -o "bin/tusd" ./cmd/tusd/main.go

static:
	@echo "+ $@"
	CGO_ENABLED=1 go build -tags "$(BUILDTAGS) cgo static_build" -ldflags "-w -extldflags -static" -o reg .

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

clean:
	@echo "+ $@"
	@rm -rf reg ./bin

install:
	@echo "+ $@"
	@go install .

docker:
	@docker build . -t $(PROJECT)
