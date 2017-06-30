.PHONY: help test build install install-deps
.DEFAULT_GOAL := help

GO_DEPENDENCIES := gopkg.in/yaml.v2 github.com/hashicorp/vault/api github.com/mitchellh/cli github.com/fatih/color

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs all tests
	GOPATH=$$(pwd) go test src/*.go

install-deps:  ## Installs go dependencies
	for dep in $(GO_DEPENDENCIES); do \
		GOPATH=$$(pwd) go get $$dep; \
	done

build: install-deps  ## Compiles the program
	GOPATH=$$(pwd) go build -o vc src/*.go

install: build  ## Install vault-client
	install -Dm755 vc $(DESTDIR)$(bindir)/vc
	install -Dm644 sample/vc-completion.bash $(DESTDIR)$(datarootdir)/bash-completion/completion/vc
