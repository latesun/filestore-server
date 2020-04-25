EXECUTABLES := app

GO ?= go
GOFMT ?= gofmt "-s"

GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

.PHONY: all
all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	$(GO) test -race -cover -coverprofile=cover.out ./...
	$(GO) tool cover -func=cover.out

.PHONY: build
build: $(EXECUTABLES)

$(EXECUTABLES):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -v -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -o bin/$@ cmd/main.go

.PHONY: clean
clean:
	-rm -rf 'bin/*'

version:
	@echo $(VERSION)
