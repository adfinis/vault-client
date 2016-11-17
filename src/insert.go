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

	path := args[0]
	data := make(map[string]interface{})

	for _, v := range args[1:] {
		kvpair := strings.Split(v, "=")
		if len(kvpair) < 2 || len(kvpair) > 2 {
			fmt.Println("Invalid key/value arguments")
			return 1
		}
		data[kvpair[0]] = kvpair[1]
	}

	_, err := vc.Logical().Write(path, data)
	if err != nil {
		fmt.Println("Unable to write secret")
		return 1
	}

	return 0
}

func (c *InsertCommand) Help() string {
	return "Remove an existing secret"
}

func (c *InsertCommand) Synopsis() string {
	return "Remove an existing secret"
}
