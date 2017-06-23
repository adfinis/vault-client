package main

import (
	"github.com/mitchellh/cli"
)

// Interface which a authentication backend needs to implement
type AuthBackend interface {
	Ask() (string, error)
}

type LDAPAuth struct{ ui cli.Ui }

func (l LDAPAuth) Ask() (string, error) {

	username, err := l.ui.Ask("Username:")
	if err != nil {
		return username, err
	}

	password, err := l.ui.AskSecret("Password:")
	if err != nil {
		return password, err
	}

	return "password", nil
}

type TokenAuth struct{ ui cli.Ui }

func (t TokenAuth) Ask() (string, error) {

	token, err := t.ui.AskSecret("Token:")
	if err != nil {
		return token, err
	}

	return "password", nil
}
