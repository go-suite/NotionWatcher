## Project details
MODULE  := github.com/gennesseaux/NotionWatcher

## go
GO           	  ?= go
GO_VERSION        ?= $(shell $(GO) version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))

## os
OS_ARCH = $(word 4, $(shell $(GO) version))
OS = $(word 1,$(subst /, ,$(OS_ARCH)))
ARCH = $(word 2,$(subst /, ,$(OS_ARCH)))

## git
GIT_COMMIT     = $(shell git rev-parse HEAD)
GIT_SHA        = $(shell git rev-parse --short HEAD)
GIT_TAG        = $(shell git describe --tags --abbrev=0 2>/dev/null || true)
GIT_TAG_COMMIT = $(shell git rev-list --abbrev-commit --tags --max-count=1)
GIT_DIRTY      = $(shell test "`git diff --name-only | wc -l`" != "0" && echo +dirty)

##
BUILD_DATE := $(shell git log -1 --format=%cd --date=format:"%d/%m/%Y %H:%M:%S")

## Version number
VERSION := $(GIT_TAG:v%=%)
ifeq ($(VERSION), )
	VERSION := 0.0.0
endif
ifneq ($(GIT_SHA), $(GIT_TAG_COMMIT))
	VERSION := $(VERSION)-$(GIT_SHA)
endif
VERSION := $(VERSION)$(GIT_DIRTY)

## project folders
BASE_BUILD_FOLDER   := build
BUILD_FOLDER        := ${BASE_BUILD_FOLDER}

## Custom build flags
LDFLAGS := -w -s
LDFLAGS += -X "${MODULE}/setup/version.Version=${VERSION}"
LDFLAGS += -X "${MODULE}/setup/version.BuildDate=${BUILD_DATE}"
LDFLAGS += -linkmode external -extldflags "-static"
LDFLAGS += $(EXT_LDFLAGS)

##
.DEFAULT_GOAL := all

.PHONY: all
all: clean generate vet fmt lint test tidy gosec build

.PHONY: clean
clean:
	$(call print-target)
	@go clean
	@rm -rf ${BASE_BUILD_FOLDER} 2>/dev/null || true
	@rm -rf ${BASE_RELEASE_FOLDER} 2>/dev/null || true
	@rm -f coverage.*

.PHONY: generate
generate:
	$(call print-target)
	@go generate ./...

.PHONY: vet
vet:
	$(call print-target)
	@go vet ./...

.PHONY: fmt
fmt:
	$(call print-target)
	@go fmt ./...

.PHONY: lint
lint:
	$(call print-target)
	@golangci-lint run

.PHONY: test
test:
	$(call print-target)
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: tidy
tidy:
	$(call print-target)
	@go mod tidy

.PHONY: gosec
gosec:
	$(call print-target)
	@gosec ./...

.PHONY: build
build:
	$(call print-target)
	CGO_ENABLED=1; go build -a --ldflags '$(LDFLAGS)' -o ${BUILD_FOLDER}/ .

.PHONY: build-docker
build-docker:
	$(call print-target)
	docker build --no-cache --tag notionwatcher .


.PHONY: version
version: ## Get current version
	$(call print-target)
	@echo "GIT_COMMIT: ${GIT_COMMIT}"
	@echo "GIT_SHA: ${GIT_SHA}"
	@echo "GIT_TAG: ${GIT_TAG}"
	@echo "GIT_TAG_COMMIT: ${GIT_TAG_COMMIT}"
	@echo "GIT_DIRTY: ${GIT_DIRTY}"
	@echo "Version: ${VERSION}"
	@echo "BuildDate: ${BUILD_DATE}"
	@echo "GoVersion: ${GO_VERSION_NUMBER}"

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
