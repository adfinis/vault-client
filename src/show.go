package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type ShowCommand struct {
	Ui cli.Ui
}

func (c *ShowCommand) Run(args []string) int {

	switch {
	case len(args) > 1:
		c.Ui.Output("The show command expects at most one argument")
		return 1
	case len(args) == 0:
		c.Ui.Output("The show command expects an argument")
		return 1
	}

	path := args[0]

	secret, err := vc.Logical().Read(path)
	if err != nil {
		c.Ui.Info(fmt.Sprintf("The was an error while retrieving the secret: %q", err))
		return 1
	}

	if secret == nil {
		c.Ui.Info("Secret does not exist")
		return 1
	}

	for k, v := range secret.Data {
		fmt.Printf("%v: %v\n", k, v)
	}
	return 0
}

func (c *ShowCommand) Help() string {
	return `Usage: vc show path

  Prints a secret specified by it's path to stdout.
`
}

func (c *ShowCommand) Synopsis() string {
	return "Show an existing secret"
}
