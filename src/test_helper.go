package main

import (
	vault "github.com/hashicorp/vault/api"
)

const TestBackend = "/test"

func SetupTestEnvironment() (Config, *vault.Client, error) {

	cfg, err := LoadConfig()
	if err != nil {
		return cfg, nil, err
	}

	client, err := InitializeClient()
	if err != nil {
		return cfg, client, err
	}

	mountConfig := vault.MountInput{
		Type:        "kv",
		Description: "vault-client integration tests",
	}
	err = vc.Sys().Mount(TestBackend, &mountConfig)
	if err != nil {
		return cfg, client, err
	}
	return cfg, client, nil
}

func TeardownTestEnvironment() error {
	return vc.Sys().Unmount(TestBackend)
}
