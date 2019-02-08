package main

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
)

type KvClient interface {
	Put(key string, value map[string]interface{}) error
	Get(key string) (map[string]interface{}, error)
	Delete(key string) error
	List(key string) ([]string, error)
	GetRawClient() *vault.Client
}

type KvClientV1 struct {
	client *vault.Client
}

func NewKvClientV1(cfg *vault.Config, token string) (*KvClientV1, error) {
	c, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)
	c.Auth()

	return &KvClientV1{client: c}, nil
}

func (c *KvClientV1) Put(key string, value map[string]interface{}) error {
	_, err := c.client.Logical().Write(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *KvClientV1) Get(key string) (map[string]interface{}, error) {
	sec, err := c.client.Logical().Read(key)
	if err != nil {
		return nil, err
	}

	if sec == nil {
		return nil, fmt.Errorf("Secret does not exist")
	}

	return sec.Data, nil
}

func (c *KvClientV1) Delete(key string) error { return nil }

func (c *KvClientV1) List(key string) ([]string, error) { return nil, nil }

func (c *KvClientV1) GetRawClient() *vault.Client { return c.client }
