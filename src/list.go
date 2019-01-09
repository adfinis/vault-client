package main

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"strings"
)

type ListCommand struct {
	Ui cli.Ui
}

func (c *ListCommand) Run(args []string) int {

	var recursiveFlag bool
	var path string
	var err error

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Output(c.Help()) }

	flags.BoolVar(&recursiveFlag, "r", false, "List secrets at path recursively")
	if err := flags.Parse(args); err != nil {
		c.Ui.Output(fmt.Sprintf("%v", err))
		return 1
	}

	args = flags.Args()

	// When no path is specified in the arguments, use "".
	// Otherwise use the last argument as the path.
	switch x := len(args); true {
	case x > 1:
		c.Ui.Output("The list command expects at most one argument")
		return 1
	case x == 0:
		path = ""
	default:
		path = strings.Trim(fmt.Sprint(args[:1]), "[]")
	}

	var paths []string
	if recursiveFlag {
		paths, err = kv.ListRecursively(path)
		if err != nil {
			c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to recursively list path: %q", err)))
			return 1
		}
	} else {
		paths, err = kv.List(path)
		if err != nil {
			c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to list path: %q", err)))
			return 1
		}
	}

	for _, path := range paths {
		c.Ui.Output(path)
	}

	return 0
}

func (c *ListCommand) Help() string {
	return `Usage: vc ls [options] path

  Lists all available secrets at the specified path.

Options:

  -r                             Recursively show all available secrets
`
}

func (c *ListCommand) Synopsis() string {
	return "List all secrets at specified path"
}
