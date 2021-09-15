.PHONY: *

ORG_AND_REPO=salesforce/ci-results-to-slack
MOUNT_DIR=/go/src/github.com/${ORG_AND_REPO}
BUILD_CONTAINER=golang:1.17-bullseye
BINARY_NAME=ci-result-to-slack
ci: build test

build: clean
	@echo "Building..."
	CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)/main.go

test: clean
	@echo "Running unit tests..."
	go test -race -coverprofile=coverage.out ./...

lint:
	@echo "Check format..."
	$(eval NEED_TO_FORMAT := $(shell go fmt ./...))
	@test -z "$(NEED_TO_FORMAT)" || (echo "Need to format the following (done for you already locally): $(NEED_TO_FORMAT)" && exit 1)
	@echo "Running linter..."
	golangci-lint run

clean:
	@echo "Cleaning..."
	rm -rf ./bin/$(BINARY_NAME)
	rm -rf coverage.out

local-docker-test: ## Build and run unit tests in docker container like CI without building the container
	docker run --rm=true -v `pwd`:$(MOUNT_DIR) $(BUILD_CONTAINER) bash -c 'cd $(MOUNT_DIR) && make ci'

local-docker-build: ## Build the container image
	docker build --no-cache -t $(ORG_AND_REPO) .
