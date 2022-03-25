.DEFAULT_GOAL:=help
SHELL:=/bin/bash

BIN ?= swap

.PHONY: help deps build install

# https://suva.sh/posts/well-documented-makefiles/
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

deps:  ## Check dependencies
	go mod tidy

$(BIN):
	go build

build: $(BIN)  ## Build the project

clean: $(BIN)  ## Delete built binary
	rm -rf $(BIN)

install:  ## Install to $GOPATH/bin
	go install

uninstall:  ## Uninstall from $GOPATH/bin
	rm -rf $$GOPATH/bin/$(BIN)

