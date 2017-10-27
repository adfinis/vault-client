.PHONY: help test build install install-deps
.DEFAULT_GOAL := help

PKGNAME=vault-client
DESCRIPTION="A command-line interface to HashiCorp's Vault "
VERSION=1.1.2

INSTALL := install
DESTDIR := /usr/local
datarootdir := $(DESTDIR)/share
bindir := $(DESTDIR)/bin

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
	$(INSTALL) -Dm755 vc $(bindir)/vc
	$(INSTALL) -Dm644 sample/vc-completion.bash $(datarootdir)/bash-completion/completions/vc

deb:  ## Create .deb package
	mkdir -p build/usr/bin build/etc/bash_completion.d
	install -Dm755 vc build/usr/bin/vc
	install -Dm644 sample/vc-completion.bash build/etc/bash_completion.d/vc
	fpm \
	  -s dir \
	  -t deb \
	  -n $(PKGNAME) \
	  -v $(VERSION) \
	  -d $(DESCRIPTION) \
	  -d bash-completion \
	  -C build \
	  .

artifacts: build deb
