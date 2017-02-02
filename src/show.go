package main

import (
	"fmt"
	"sort"

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
	MaxKeyLen := 0
	for k, _ := range secret.Data {
		if KeyLen := len(k); KeyLen > MaxKeyLen {
			MaxKeyLen = KeyLen
		}
	}

	// Add an additional X whitespaces between "key:" and "value"
	MaxKeyLen += 4

	// Sort secrets lexicographically
	var keys []string
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		c.Ui.Output(fmt.Sprintf("%-"+fmt.Sprint(MaxKeyLen)+"v %v",
			k+":",           // Secret identifier
			secret.Data[k])) // Secret value
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
