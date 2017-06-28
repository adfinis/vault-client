package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mitchellh/cli"
)

// Interface to easily add new authentication backends.
type AuthBackend interface {
	Ask() (*http.Request, error)
}

type LDAPAuth struct {
	ui cli.Ui
}

func (l LDAPAuth) Ask() (*http.Request, error) {

	username, err := l.ui.Ask("Username:")
	if err != nil {
		return new(http.Request), err
	}

	password, err := l.ui.AskSecret("Password:")
	if err != nil {
		return new(http.Request), err
	}

	body := []byte(fmt.Sprintf(`{"password":"%s"}`, password))
	url := fmt.Sprintf("%v/v1/auth/%s/login/%s", ComposeUrl(), cfg.AuthMethod, username)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	return req, nil
}
