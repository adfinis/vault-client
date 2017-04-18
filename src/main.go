package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

var cc *consul.Client
var vc *vault.Client
var cfg Config

func main() {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = InitializeConsulClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = InitializeVaultClient()
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

// Initalizes a globally accessible Consul HTTP API client
func InitializeConsulClient() error {

	var err error

	var protocol string

	if cfg.Consul.TLS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	ccfg := consul.DefaultConfig()

	// Compose Consul HTTP URL
	ccfg.Address = fmt.Sprintf("%v://%v:%v", protocol, cfg.Consul.Host, cfg.Consul.Port)

	cc, err = consul.NewClient(ccfg)
	if err != nil {
		return err
	}

	return nil
}

// Initalizes a globally accessible Vault HTTP API client
func InitializeVaultClient() error {

	var protocol string

	if cfg.Vault.TLS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	tr := &http.Transport{}

	if !cfg.Vault.VerifyTLS {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	vcfg := vault.Config{
		Address:    fmt.Sprintf("%v://%v:%v", protocol, cfg.Vault.Host, cfg.Vault.Port),
		HttpClient: &http.Client{Transport: tr},
	}

	var err error

	vc, err = vault.NewClient(&vcfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	vc.SetToken(cfg.Vault.Token)
	vc.Auth()

	return nil
}

func LoadCli() *cli.CLI {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("vc", "1.0")
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
