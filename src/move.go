package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type MoveCommand struct {
	Ui cli.Ui
}

func (c *MoveCommand) Run(args []string) int {

	if len(args) != 2 {
		c.Ui.Error("The move command expects a source and a destination path")
		return 1
	}

	secret, err := vc.Logical().Read(args[0])
	if err != nil {
		fmt.Println("Unable to find source secret")
		return 1
	}

	if secret == nil {
		c.Ui.Error("Source secret doesn't exist")
		return 1
	}

	_, err = vc.Logical().Write(args[1], secret.Data)
	if err != nil {
		fmt.Println("Unable to write destination secret")
		return 1
	}

	_, err = vc.Logical().Delete(args[0])
	if err != nil {
		fmt.Println("Unable to remove source secret")
		return 1
	}

	return 0
}

func (c *MoveCommand) Help() string {
	return "Move an existing secret to another location"
}

func (c *MoveCommand) Synopsis() string {
	return "Move an existing secret to another location"
}
