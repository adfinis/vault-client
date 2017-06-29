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

	var MaxKeyLen = 0
	var keys []string
	var has_comments = false

	for key, _ := range secret.Data {
		// Get length of the largest key in order to calculate the
		// "whitespace padded" representation of `show`
		if KeyLen := len(key); KeyLen > MaxKeyLen {
			MaxKeyLen = KeyLen + 4
		}

		if strings.HasSuffix(key, "_comment") {
			// Check whether a secret contains comments
			has_comments = true
		} else {
			keys = append(keys, key)
		}
	}

	// Sort secrets lexicographically
	sort.Strings(keys)

	// Only pad K/V pairs when a secret containts no comments
	kv_output_format := "%-" + fmt.Sprint(MaxKeyLen) + "v %v\n"
	if has_comments || len(secret.Data) == 1 {
		kv_output_format = "%v %v\n"
	}

	output := ""

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
		output += fmt.Sprintf(kv_output_format,
			key+":",
			secret.Data[key])
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
