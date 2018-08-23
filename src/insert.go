package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type InsertCommand struct {
	Ui cli.Ui
}

func (c *InsertCommand) Run(args []string) int {

	if len(args) < 2 {
		c.Ui.Error("The insert command expects at least a path and a k/v pair (key=value)")
		return 1
	}

	path := args[0]
	data := make(map[string]interface{})

	for _, v := range args[1:] {
		kvpair := strings.SplitN(v, "=", 2)
		if len(kvpair) < 2 || len(kvpair) > 2 {
			c.Ui.Error(fmt.Sprintf("Invalid key/value arguments: %q", v))
			return 1
		}
		data[kvpair[0]] = kvpair[1]
	}

	_, err := vc.Logical().Write(path, data)
	if err != nil {
		c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to write secret: %q", err)))
		return 1
	}

	return 0
}

func (c *InsertCommand) Help() string {
	return `Usage: vc insert key1=value1 key2=value2...

  Inserts a new secret at the specified path with a set of key/value pairs.
`
}

func (c *InsertCommand) Synopsis() string {
	return "Insert an new secret"
}
