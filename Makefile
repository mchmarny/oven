RELEASE_VERSION    ?=v0.4.1
EMULATOR_IMAGE     ?=ghcr.io/mchmarny/firestore-emulator:v0.3.2
EMULATOR_HOST      ?=localhost
EMULATOR_PORT      ?=8888
EMULATOR_PROJECT   ?=oven

all: help

version: ## Prints the current version
	@echo $(RELEASE_VERSION)
.PHONY: version

tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor
.PHONY: tidy

upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor
.PHONY: upgrade

test: tidy ## Runs unit tests
	go test -short -count=1 -race -covermode=atomic -coverprofile=cover.out ./...
.PHONY: test

bench: tidy ## Runs benchmark tests
	go test -bench=. -count 5 -run=^# -benchtime=10s -benchmem ./...
.PHONY: bench

integration: tidy ## Runs integration tests
	PROJECT_ID=$(EMULATOR_PROJECT) \
	FIRESTORE_EMULATOR_HOST="$(EMULATOR_HOST):$(EMULATOR_PORT)" \
	go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...
.PHONY: integration

examples: tidy ## Runs all examples
	# crud
	PROJECT_ID=$(EMULATOR_PROJECT) \
	FIRESTORE_EMULATOR_HOST="$(EMULATOR_HOST):$(EMULATOR_PORT)" \
	go run examples/crud/main.go

	# query
	PROJECT_ID=$(EMULATOR_PROJECT) \
	FIRESTORE_EMULATOR_HOST="$(EMULATOR_HOST):$(EMULATOR_PORT)" \
	go run examples/query/main.go
.PHONY: examples

cover: test ## Runs unit tests and putputs coverage
	go tool cover -func=cover.out
.PHONY: cover

lint: ## Lints the entire project 
	golangci-lint -c .golangci.yaml run
.PHONY: lint

store: ## Run Firestore emulator image
	tools/fs/run "$(EMULATOR_IMAGE)" "$(EMULATOR_PROJECT)" "$(EMULATOR_HOST)" "$(EMULATOR_PORT)" 
.PHONY: store

storedown: ## Stop previously launched Firestore emulator
	tools/fs/stop
.PHONY: storedown

tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)
.PHONY: tag

tagless: ## Delete the current release tag 
	git tag -d $(RELEASE_VERSION)
	git push --delete origin $(RELEASE_VERSION)
.PHONY: tagless

clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
.PHONY: clean

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help