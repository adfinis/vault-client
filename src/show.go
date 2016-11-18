package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type ShowCommand struct {
	Ui cli.Ui
}

func (c *ShowCommand) Run(args []string) int {

	// use the last argument as path
	path := strings.Join(args[len(args)-1:], "")

	secret, err := vc.Logical().Read(path)
	if err != nil {
		c.Ui.Info("The was an error while retrieving the secret")
		return 1
	}

	if secret == nil {
		c.Ui.Info("Secret does not exist")
		return 1
	}

	for k, v := range secret.Data {
		fmt.Printf("%v = %v\n", k, v)
	}
	return 0
}

func (c *ShowCommand) Help() string {
	return "Copy an existing secret to another location"
}

func (c *ShowCommand) Synopsis() string {
	return "Copy an existing secret to another location"
}
