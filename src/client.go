package main

import (
	vault "github.com/hashicorp/vault/api"
)

// A vault api client that transparently supports version 1 & 2 key/value engines.
type KvClient struct {
	Client *vault.Client
}


func NewKvClient(cfg *vault.Config, token string) (*KvClient, error) {
	c, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)
	c.Auth()

	return &KvClient{Client: c}, nil
}


func (c *KvClient) Read(path string) (*vault.Secret, error) {
	return c.Client.Logical().Read(path)
}

func (c *KvClient) Delete(path string) (*vault.Secret, error) {
	return c.Client.Logical().Delete(path)
}

func (c *KvClient) Write(path string, data map[string]interface{}) (*vault.Secret, error) {
	return c.Client.Logical().Write(path, data)
}

func (c *KvClient) List(path string) (*vault.Secret, error) {
	return c.Client.Logical().List(path)
}
