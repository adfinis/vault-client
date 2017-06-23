package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/cli"
)

type VaultAuthenticationResponse struct {
	LeaseID       string  `json:"lease_id"`
	Renewable     bool    `json:"renewable"`
	LeaseDuration int     `json:"lease_duration"`
	Data          *string `json:"data"`
	Auth          struct {
		ClientToken string   `json:"client_token"`
		Policies    []string `json:"policies"`
		Metadata    struct {
			Username string `json:"username`
		}
		LeaseDuration int  `json:"lease_duration"`
		Renewable     bool `json:"renewable"`
	}
}

// Authenticate against vault using the configured method and return a valid token
func GetAuthenticationToken(ui cli.Ui) (string, error) {

	var am AuthBackend

	// Currently supported authentication backends
	switch cfg.AuthBackend {
	case "token":
		return ui.AskSecret("Token:")
	case "ldap":
		am = LDAPAuth{ui}
	}

	// Collect information such as username, password, ...
	req, err := am.Ask()
	if err != nil {
		return "", fmt.Errorf("Unable to parse input: %q", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("Unable retrieve authentication token from vault (status code %v)", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to parse request body", err)
	}

	var data VaultAuthenticationResponse
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", fmt.Errorf("Unable to parse request body", err)
	}

	return data.Auth.ClientToken, nil
}

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
