package main

import (
	"fmt"
	"sort"
	"strings"

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

	var maxKeyLength = 0
	var keys []string

	for key, _ := range secret.Data {

		// Get the length of the largest key in order to use the largest offset for the
		// "whitespace padded" representation.
		if len(key) > maxKeyLength {
			maxKeyLength = len(key) + 2
		}

		// Ignore k/v pair that are comments
		if !strings.HasSuffix(key, "_comment") {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var output string
	for _, key := range keys {

		if value, exists := secret.Data[key+"_comment"].(string); exists {

			if multilineComments := strings.Split(value, "\n"); len(multilineComments) > 1 {
				for _, comment := range multilineComments {
					output += "#" + comment + "\n"
				}
			} else {
				output += "#" + value + "\n"
			}
		}

		output += fmt.Sprintf(
			"%-"+fmt.Sprint(maxKeyLength)+"v%v\n",
			key+":",
			strings.TrimSpace(secret.Data[key].(string)),
		)
	}

	c.Ui.Output(output)

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
