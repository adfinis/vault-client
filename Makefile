.PHONY: help data dependencies

help: 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

AUTH_SECRET = password

data: ## Write some test data to vault
	for secret in secret0 secret2 secret3; do \
		vault auth --address http://127.0.0.1:8200 $$AUTH_SECRET; \
		vault write --address http://127.0.0.1:8200 secret/$$secret key=$$secret; \
	done

dependencies: ## install go dependencies
	for dep in gopkg.in/yaml.v2 github.com/hashicorp/vault/api github.com/mitchellh/cli; do \
		go get $$dep; \
	done
