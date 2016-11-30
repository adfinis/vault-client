.PHONY: help test-all dependencies

help: 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-all: ## Runs all tests
	go test src/*.go

build: ## Compiles the program
	go build -o vc src/*.go

dependencies: ## install go dependencies
	for dep in gopkg.in/yaml.v2 github.com/hashicorp/vault/api github.com/mitchellh/cli; do \
		go get $$dep; \
	done
