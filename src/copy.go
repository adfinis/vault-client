package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type CopyCommand struct {
	Ui cli.Ui
}

func (c *CopyCommand) Run(args []string) int {

	if len(args) != 2 {
		c.Ui.Error("The copy command expects a source and a destination path")
		return 1
	}

	secret, err := vc.Read(args[0])
	if err != nil {
		c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to find source secret: %q", err)))
		return 1
	}

	if secret == nil {
		c.Ui.Error("Source secret doesn't exist")
		return 1
	}

	_, err = vc.Write(args[1], secret.Data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to write destination secret: %q", err))
		return 1
	}

	return 0
}

func (c *CopyCommand) Help() string {
	return `Usage: vc cp source dest

  Copies an existing secret to a new destination path.
  The source secret will be preserved.
`
}

func (c *CopyCommand) Synopsis() string {
	return "Copy an existing secret to another location"
}
