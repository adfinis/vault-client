package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type CopyCommand struct {
	Ui cli.Ui
}

func (c *CopyCommand) Run(args []string) int {

	secret, err := vc.Logical().Read(args[0])
	if err != nil {
		fmt.Println("Unable to find source secret")
		return 1
	}

	if secret == nil {
		fmt.Println("Source secret doesn't exist")
		return 1
	}

	_, err = vc.Logical().Write(args[1], secret.Data)
	if err != nil {
		fmt.Println("Unable to write destination secret")
		return 1
	}

	return 0
}

func (c *CopyCommand) Help() string {
	return "Copy an existing secret to another location"
}

func (c *CopyCommand) Synopsis() string {
	return "Copy an existing secret to another location"
}
