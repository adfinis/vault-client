.PHONY: help test build


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs test suite
	go test -v -failfast src/*.go

build: ## Compiles vc
	go build -o vc src/*.go
