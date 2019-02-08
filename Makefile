.PHONY: help test build install install-deps
.DEFAULT_GOAL := help

PKGNAME=vault-client
DESCRIPTION="A command-line interface to HashiCorp's Vault "
VERSION=1.1.4

INSTALL		:= install

# Common prefix for installation directories.
# NOTE: This directory must exist when you start the install.
prefix = /usr/local
datarootdir = $(prefix)/share
datadir = $(datarootdir)
exec_prefix = $(prefix)
# Where to put the executable for the command 'gcc'.
bindir = $(exec_prefix)/bin
# Where to put the directories used by the compiler.
libexecdir = $(exec_prefix)/libexec
# Where to put the Info files.
infodir = $(datarootdir)/info

GO_DEPENDENCIES := gopkg.in/yaml.v2 github.com/hashicorp/vault/api github.com/mitchellh/cli github.com/fatih/color

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs all tests
	GOPATH=$$(pwd)/vendor go test src/*.go

install-deps:  ## Installs go dependencies
	for dep in $(GO_DEPENDENCIES); do \
		GOPATH=$$(pwd)/vendor go get -v $$dep; \
	done

build: install-deps  ## Compiles the program
	GOPATH=$$(pwd)/vendor go build -o vc src/*.go

install: build  ## Install vault-client
	$(INSTALL) -Dm755 vc $(DESTDIR)$(bindir)/vc
	$(INSTALL) -Dm644 sample/vc-completion.bash $(DESTDIR)$(datarootdir)/bash-completion/completions/vc
	$(INSTALL) -Dm644 sample/vc-completion.zsh $(DESTDIR)$(datarootdir)/zsh/site-functions/_vc

artifacts: build deb
