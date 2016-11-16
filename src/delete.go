package main

import (
	"github.com/mitchellh/cli"
)

type DeleteCommand struct {
	Ui cli.Ui
}

func (c *DeleteCommand) Run(args []string) int {

	_, err := vc.Logical().Delete(args[0])
	if err != nil {
		return 1
	}

	return 0
}

func (c *DeleteCommand) Help() string {
	return "Remove an existing secret"
}

func (c *DeleteCommand) Synopsis() string {
	return "Remove an existing secret"
}
