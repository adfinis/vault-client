package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/cli"
)

type IndexCommand struct {
	Ui cli.Ui
}

func (c *IndexCommand) Run(_ []string) int {

	c.Ui.Output("Indexing all available paths in vault")

	index, err := BuildIndex()
	if err != nil {
		return 1
	}

	r := strings.NewReader(strings.Join(index[:], "\n"))

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return 1
	}

	ioutil.WriteFile(cfg.IndexFile, data, 0644)

	return 0
}

func (c *IndexCommand) Help() string {
	return "Index all available paths of accessable generic backends"
}

func (c *IndexCommand) Synopsis() string {
	return "Index all available paths of accessable generic backends"
}

// Returns all accessable generic backends
func GetGenericBackends() ([]string, error) {

	mounts, err := vc.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var backends []string

	for x, i := range mounts {
		if i.Type == "generic" {
			backends = append(backends, x)
		}
	}

	return backends, nil
}

// Build a list of all available paths
func BuildIndex() ([]string, error) {

	backends, err := GetGenericBackends()
	if err != nil {
		return nil, err
	}

	var index []string

	for _, backend := range backends {
		paths, err := WalkPath(backend)
		if err != nil {
			return nil, err
		}
		index = append(index, paths...)
	}

	for _, v := range index {
		fmt.Println(v)
	}

	return index, nil
}

// Recursively walks an accessable path
func WalkPath(startpath string) ([]string, error) {

	var paths []string

	secret, err := vc.Logical().List(startpath)
	if err != nil {
		return nil, err
	}

	for _, path := range secret.Data {

		// expecting "[secret0 secret1 secret2...]"
		secrets := strings.Split(strings.Trim(fmt.Sprint(path), "[]"), " ")

		for _, v := range secrets {

			path_to_secret := fmt.Sprint(startpath, v)

			if !strings.HasSuffix(v, "/") {
				paths = append(paths, path_to_secret)
			} else {

				child_paths, err := WalkPath(path_to_secret)
				if err != nil {
					return nil, err
				}
				paths = append(paths, child_paths...)

			}
		}
	}

	return paths, nil
}
