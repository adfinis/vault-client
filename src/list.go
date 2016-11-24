package main

import (
	"fmt"

	"github.com/mitchellh/cli"
	"strings"
)

type ListCommand struct {
	Ui cli.Ui
}

func (c *ListCommand) Run(args []string) int {

	if len(args) > 1 {
		c.Ui.Output("The list command expects at most one argument")
		return 1
	}

	var path string
	var childs []string
	var err error

	// When no path is specified in the arguments, use ""
	if len(args) == 0 {
		path = ""
	} else {
		path = args[0]
	}

	// Top level directories are mounts
	if path == "/" || path == "" {
		childs, err = GetGenericBackends()
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to get a list generic backends: %q", err))
			return 1
		}

	} else {
		childs, err = ListPath(path)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to get a list of secrets: %q", err))
			return 1
		}

	}

	for _, child := range childs {
		c.Ui.Output(child)
	}

	return 0
}

func (c *ListCommand) Help() string {
	return "List all available secrets and subdirectory of path"
}

func (c *ListCommand) Synopsis() string {
	return "List all available secrets and subdirectory of path"
}

func ListPath(path string) ([]string, error) {

	var childs []string

	secret, err := vc.Logical().List(path)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, nil
	}

	for _, path := range secret.Data {

		// expecting "[secret0 secret1 secret2...]"
		childs = strings.Split(strings.Trim(fmt.Sprint(path), "[]"), " ")
	}

	return childs, nil
}
