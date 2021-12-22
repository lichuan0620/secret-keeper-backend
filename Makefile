# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make                - default to 'build' target
#   make lint           - code analysis
#   make test           - run unit test (or plus integration test)
#   make build          - alias to build-local target
#   make container      - build containers
#   make push           - push containers
#   make clean          - clean up targets

# Module name.
NAME := secret-keeper

# Container registries.
REGISTRY ?= cr-cn-beijing.volces.com/sailor-moon

# Project output directory.
OUTPUT_DIR := ./bin

# Module version, you might want to change this if not building on a git tag.
VERSION ?= $(shell git describe --tags --always --dirty)

# Default golang flags used in build and test
# -count: run each test and benchmark 1 times. Set this flag to disable test cache
export GOFLAGS ?= -count=1

#
# Define all targets. At least the following commands are required:
#

# All targets.
.PHONY: test build container push clean

test:
	@go test -race -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

build:
	@go build -v -o $(OUTPUT_DIR)/$(NAME) .;

container:
	@docker build -t $(REGISTRY)/$(NAME):$(VERSION) .

push: container
	@docker push $(REGISTRY)/$(NAME):$(VERSION);

clean:
	@rm -vrf ${OUTPUT_DIR}