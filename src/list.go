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
	var paths []string
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

	if recursiveFlag {
		paths, err = RecursivelyListSecrets(path)
		if err != nil {
			c.Ui.Error(CheckError(err, fmt.Sprintf("Unable to recursively list path: %q", err)))
			return 1
		}
	} else {
		paths, err = ListSecrets(path)
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

// Build a list of all available paths
func RecursivelyListSecrets(path string) ([]string, error) {

	// Could be secrets or backends...
	var items []string

	secrets, err := ListSecrets(path)
	if err != nil {
		return nil, err
	}

	// Return if path holds no secrets
	if len(secrets) == 0 {
		return nil, nil
	}

	for _, secret := range secrets {

		// Prefix secret with it's path
		secret = fmt.Sprint(path, secret)

		// Recurse if path is a directory
		if strings.HasSuffix(secret, "/") {
			childs, err := RecursivelyListSecrets(secret)
			if err != nil {
				return nil, err
			}
			items = append(items, childs...)
		} else {
			items = append(items, secret)
		}

	}

	return items, err

}

func ListSecrets(path string) ([]string, error) {

	// Could be secrets or backends...
	var items []string
	var err error

	// Vault can have multiple backends mounted at different paths (e.g "/customer1", "/customer2"...).
	// `vc` only cares about generic backends.
	if path == "/" || path == "" {

		items, err = ListKvBackends()
		if err != nil {
			return nil, err
		}

	} else {

		secret, err := kv.List(path)
		if err != nil {
			return nil, err
		}

		if secret == nil {
			return nil, nil
		}

		for _, path := range secret.Data {
			// expecting "[secret0 secret1 secret2...]"
			//
			// if the name both exists as directory and as file
			// e.g. "/secret/" and "/secret" it will print an empty line
			items = strings.Split(strings.Trim(fmt.Sprint(path), "[]"), " ")
		}
	}

	return items, nil
}

// Returns the paths to all of all kv backends.
func ListKvBackends() ([]string, error) {

	mounts, err := kv.Client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var backends []string

	for x, i := range mounts {
		// With the 0.8.3 release of vault the "generic" backend was renamed to "kv". For backwards
		// compatibility consider both of them
		if i.Type == "kv" || i.Type == "generic" {
			backends = append(backends, x)
		}
	}

	return backends, nil
}
