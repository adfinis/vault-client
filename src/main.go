package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

var vc *vault.Client
var cfg Config

func main() {

	_, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = InitializeClient()

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

func InitializeClient() (*vault.Client, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.VerifyTLS},
	}

	config := vault.Config{
		Address:    ComposeUrl(),
		HttpClient: &http.Client{Transport: tr},
		Timeout:    3 * time.Second,
	}

	var err error

	vc, err = vault.NewClient(&config)
	if err != nil {
		return vc, err
	}

	vc.SetToken(cfg.Token)
	vc.Auth()
	return vc, nil
}

func LoadCli() *cli.CLI {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("vc", "1.1.4")
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
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				Ui: ui,
			}, nil
		},
	}

	return c
}
