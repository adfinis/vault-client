package main

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"strings"
)

type KvClient interface {
	Put(key string, value map[string]interface{}) error
	Get(key string) (map[string]interface{}, error)
	Delete(key string) error
	List(key string) ([]string, error)
	Traverse(key string) ([]string, error)
	RawClient() *vault.Client
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

	// TODO: Introduce new error type ErrSecretDoesNotExist
	if sec == nil {
		return nil, fmt.Errorf("Secret does not exist")
	}

	return sec.Data, nil
}

func (c *KvClientV1) Delete(key string) error {
	_, err := c.client.Logical().Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (c *KvClientV1) Traverse(key string) ([]string, error) {
	childs, err := c.List(key)
	if err != nil {
		return nil, err
	}

	if len(childs) == 0 {
		return nil, nil
	}

	var paths []string

	for _, child := range childs {
		// Prefix child with path
		child = fmt.Sprint(key, child)
		if strings.HasSuffix(child, "/") {
			childs, err := c.Traverse(child)
			if err != nil {
				return nil, err
			}
			paths = append(paths, childs...)
		} else {
			paths = append(paths, child)
		}
	}
	return paths, nil
}

func (c *KvClientV1) List(key string) ([]string, error) {

	key = strings.TrimLeft(key, "/")
	if key == "" {
		return kvMounts(c)
	}

	secret, err := c.client.Logical().List(key)
	if err != nil {
		return nil, err
	}

	if secret == nil {

		// Key could be an empty backend
		var tmp string
		if !strings.HasPrefix("/", key) {
			tmp = key + "/"
		}

		mounts, err := kvMounts(c)
		if err != nil {
			return nil, err
		}
		for _, k := range mounts {
			if tmp == k {
				return []string{k}, nil
			}
		}

		// Or it could simply not exist...
		return nil, fmt.Errorf("Secret does not exist")
	}

	var childs []string
	for _, path := range secret.Data {
		childs = strings.Split(strings.Trim(fmt.Sprint(path), "[]"), " ")
	}

	return childs, nil

}

func (c *KvClientV1) RawClient() *vault.Client { return c.client }

// Return a list of kv engines
func kvMounts(c *KvClientV1) ([]string, error) {
	mounts, err := c.client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var backends []string
	for x, i := range mounts {
		if i.Type == "kv" {
			backends = append(backends, x)
		}
	}
	return backends, nil
}
