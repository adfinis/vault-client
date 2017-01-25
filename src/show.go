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
		c.Ui.Error(fmt.Sprintf("There was an error while retrieving the secret: %q", err))
		return 1
	}

	if secret == nil {
		c.Ui.Error("Secret does not exist")
		return 1
	}

	// Get length of the largest key in order to calculate the
	// "whitespace padded" representation of `show`
	max_key_len := 0
	for k, _ := range secret.Data {
		if key_len := len(k); key_len > max_key_len {
			max_key_len = key_len
		}
	}

	// Add an additional X whitespaces between "key:" and "value"
	max_key_len += 4

	for k, v := range secret.Data {
		c.Ui.Output(fmt.Sprintf("%-"+fmt.Sprint(max_key_len)+"v %v", fmt.Sprint(k, ":"), v))
	}

	return 0
}

func (c *ShowCommand) Help() string {
	return `Usage: vc show path

  Prints a secret specified by its path to stdout.
`
}

func (c *ShowCommand) Synopsis() string {
	return "Show an existing secret"
}
