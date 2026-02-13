GO_FILES := $(shell find . -name '*.go')

CURR_TAG := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u '+%Y-%m-%d-%I:%M:%S%p')
COMMIT := $(shell git rev-parse HEAD)
LD_FLAGS := -X main.version=$(CURR_TAG) -X main.date=$(BUILD_TIME) -X main.commit=$(COMMIT)

.PHONY: build

build: wikitable2csv

wikitable2csv: wikitable2csv.go $(GO_FILES)
	go build -ldflags "$(LD_FLAGS)" -o $@ .
