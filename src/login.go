package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
)

type LoginCommand struct {
	Ui cli.Ui
}

func (c *LoginCommand) Run(args []string) int {

	var statusFlag bool
	flags := flag.NewFlagSet("login", flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Output(c.Help()) }

	flags.BoolVar(&statusFlag, "s", false, "Show the status of your current token")
	if err := flags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("%v", err))
		return 1
	}

	args = flags.Args()

	if len(args) > 0 {
		c.Ui.Error("The login command does not expect arguments")
		return 1
	}

	if !statusFlag {

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

		c.Ui.Output(fmt.Sprintf("Automatically stored the retrieved token in %q", cfg.Path))

	}

	ttl, err := GetTokenTTL(cfg.Token)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to retrieve ttl: %q", err))
		return 1
	}
	c.Ui.Output(fmt.Sprintf("Your token will expire on %v", ttl.Format("02/01/2006 15:04:05")))

	return 0
}

func (c *LoginCommand) Help() string {
	return `Usage: vc login 

  Authenticates against Vault through your prefered method 
  e.g (username/password, ldap) and stores the retrieved
  Token in your ~/.vaultrc

Options:

  -s                             Shows when your current token will expire
`
}

func (c *LoginCommand) Synopsis() string {
	return "Authenticate against Vault using your prefered method"
}
