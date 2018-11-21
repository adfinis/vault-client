package main

import (
	"fmt"

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

func (c *KvClient) Put(path string, data map[string]interface{}) (*vault.Secret, error) {

	mountPath, v2, err := isKVv2(path, kv.Client)
	if err != nil {
		return nil, fmt.Errorf("Unable to determin kv engine version: %q", err)
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		if err != nil {
			return nil, fmt.Errorf("Unable to patch path with prefix: %q", err)
		}
	}

	if v2 {
		data["data"] = data
	}

	sec, err:= c.Client.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}
	return sec, nil
}


func (c *KvClient) Get(path string) (map[string]interface{}, error) {
	mountPath, v2, err := isKVv2(path, kv.Client)
	if err != nil {
		return nil, fmt.Errorf("Unable to determin kv engine version: %q", err)
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		if err != nil {
			return nil, fmt.Errorf("Unable to patch path with prefix: %q", err)
		}
	}

	sec, err:= c.Client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if sec == nil {
		return make(map[string]interface{}), nil
	}

	if v2 {
		return sec.Data["data"].(map[string]interface{}), nil
	}
	return sec.Data, nil
}


func (c *KvClient) Delete(path string) (*vault.Secret, error) {
	return c.Client.Logical().Delete(path)
}

func (c *KvClient) List(path string) (*vault.Secret, error) {
	return c.Client.Logical().List(path)
}
