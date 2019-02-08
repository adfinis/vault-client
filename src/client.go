package main

import (
	vault "github.com/hashicorp/vault/api"
)

type KvClient interface {
	Put(key string, value map[string]interface{}) error
	Get(key string) (map[string]interface{}, error)
	Delete(key string) error
	List(key string) ([]string, error)
}

type KvClientV1 struct {
	Client *vault.Client
}

func NewKvClientV1(cfg *vault.Config, token string) (*KvClientV1, error) {
	c, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)
	c.Auth()

	return &KvClientV1{Client: c}, nil
}

func (c *KvClientV1) Put(key string, value map[string]interface{}) error {
	_, err := c.Client.Logical().Write(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *KvClientV1) Get(key string) (map[string]interface{}, error) { return nil, nil }

func (c *KvClientV1) Delete(key string) error { return nil }

func (c *KvClientV1) List(key string) ([]string, error) { return nil, nil }
