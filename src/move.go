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

	srcPath := args[0]
	destPath := args[1]

	kvPairs, err := kv.Get(srcPath)
	if err != nil {
		c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to find source secret: %q", err)))
		return 1
	}

	if kv.Put(destPath, kvPairs) != nil {
		fmt.Println("Unable to write destination secret")
		return 1
	}

	if kv.Delete(srcPath) != nil {
		fmt.Println("Unable to remove source secret")
		return 1
	}

	return 0
}

func (c *MoveCommand) Help() string {
	return `Usage: vc mv source dest

  Moves an existing secret to a new destination path.
  The source secret will be removed.
`
}

func (c *MoveCommand) Synopsis() string {
	return "Move an existing secret to another location"
}
