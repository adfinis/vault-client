package main

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

var vc *vault.Client
var cfg Config

func main() {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = InitializeClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	c := LoadCli()

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}

func InitializeClient(cfg Config) error {

	vcfg := vault.Config{
		Address: fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port),
	}

	var err error

	vc, err = vault.NewClient(&vcfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	vc.SetToken(cfg.Token)
	vc.Auth()

	return nil
}

func LoadCli() *cli.CLI {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("vc", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"edit": func() (cli.Command, error) {
			return &EditCommand{
				Ui: ui,
			}, nil
		},
		"rm": func() (cli.Command, error) {
			return &DeleteCommand{
				Ui: ui,
			}, nil
		},
		"insert": func() (cli.Command, error) {
			return &InsertCommand{
				Ui: ui,
			}, nil

		},
		"mv": func() (cli.Command, error) {
			return &MoveCommand{
				Ui: ui,
			}, nil
		},
		"cp": func() (cli.Command, error) {
			return &CopyCommand{
				Ui: ui,
			}, nil
		},
		"show": func() (cli.Command, error) {
			return &ShowCommand{
				Ui: ui,
			}, nil
		},
		"ls": func() (cli.Command, error) {
			return &ListCommand{
				Ui: ui,
			}, nil
		},
	}

	return c
}
