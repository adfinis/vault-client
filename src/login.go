package main

import (
	"fmt"
	"github.com/mitchellh/cli"
)

type LoginCommand struct {
	Ui cli.Ui
}

func (c *LoginCommand) Run(args []string) int {

	if len(args) > 0 {
		c.Ui.Error("The login command does not expect arguments")
		return 1
	}

	var err error

	token, err := GetAuthenticationToken(c.Ui)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to retrieve token: %q", err))
		return 1
	}

	err = UpdateConfigToken(token)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to retrieve token: %q", err))
		return 1
	}

	return 0
}

func (c *LoginCommand) Help() string {
	return `Usage: vc login 

  Authenticates against Vault thorugh your prefered method 
  e.g (username/password, ldap) and stores the retrieved
  Token in your ~/.vaultrc
`
}

func (c *LoginCommand) Synopsis() string {
	return "Authenticate against Vault using your prefered method"
}
