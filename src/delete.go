package main

import (
	"fmt"
	"github.com/mitchellh/cli"
)

type DeleteCommand struct {
	Ui cli.Ui
}

func (c *DeleteCommand) Run(args []string) int {

	switch {
	case len(args) > 1:
		c.Ui.Error("The rm command expects at most one argument")
		return 1
	case len(args) == 0:
		c.Ui.Error("The rm command expects an argument")
		return 1
	}

	if len(args) > 1 {
		c.Ui.Error("The rm command expects at most one argument")
		return 1
	}

	path := args[0]

	secret, err := vc.Logical().Read(path)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if secret == nil {
		c.Ui.Output("Secret does not exist")
		return 1
	}

	_, err = vc.Logical().Delete(path)
	if err != nil {
		fmt.Println(err)
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
