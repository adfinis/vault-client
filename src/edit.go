package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/cli"
)

type EditCommand struct {
	Ui   cli.Ui
	Path string
}

func (c *EditCommand) Run(args []string) int {

	c.Path = args[0]

	_, err := vc.Logical().Read(c.Path)
	if err != nil {
		return 1
	}

	_, err = EditFile()
	if err != nil {
		return 1
	}

	return 0
}

func (c *EditCommand) Help() string {
	return "Edit a secret"
}

func (c *EditCommand) Synopsis() string {
	return "Edit a secret"
}

func EditFile() ([]byte, error) {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	f, err := ioutil.TempFile("", "vaultsecret")
	if err != nil {
		return nil, err
	}

	defer os.Remove(f.Name())

	fmt.Println(f.Name())

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	return content, nil
}
