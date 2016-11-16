package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type IndexCommand struct {
	Ui cli.Ui
}

func (c *IndexCommand) Run(_ []string) int {

	c.Ui.Output("Indexing all available paths in vault")

	_, err := BuildIndex()
	if err != nil {
		return 1
	}

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
			fmt.Println(x)
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
		_, err := WalkPath(backend)
		if err != nil {
			return nil, err
		}
	}
	// TODO: append paths to index
	return index, nil
}

// Recursively walks an accessable path
func WalkPath(startpath string) ([]string, error) {

	var paths []string

	secret, err := vc.Logical().List(startpath)
	if err != nil {
		return nil, err
	}

	// Check whether secret is empty or not
	if secret != nil {
		for _, path := range secret.Data {
			//TODO: Recursively call WalkPath until the entire tree is discovered
			fmt.Println(path)
		}
	}
	//TODO: Return a slice of all paths
	return append(paths, ""), nil

}
