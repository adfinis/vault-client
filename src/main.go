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
	}

	err = InitializeClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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

	vc.SetToken(cfg.Password)
	vc.Auth()

	return nil
}

func LoadCli() *cli.CLI {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("vault", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"edit": func() (cli.Command, error) {
			return &EditCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"index": func() (cli.Command, error) {
			return &IndexCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"rm": func() (cli.Command, error) {
			return &DeleteCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"insert": func() (cli.Command, error) {
			return &InsertCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"mv": func() (cli.Command, error) {
			return &MoveCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		// TODO: cp
		// TODO: ls
	}

	return c
}
