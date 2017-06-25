package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

// Authenticate against vault using the configured method and return a valid token
func GetAuthenticationToken(ui cli.Ui) (string, error) {

	var am AuthBackend

	// Currently supported authentication backends
	switch cfg.AuthBackend {
	case "token":
		token, err := ui.AskSecret("Token:")
		if err != nil {
			return "", fmt.Errorf("Unable to parse input: %q", err)
		}
		return token, nil
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

	var secret vault.Secret
	err = json.Unmarshal(body, &secret)
	if err != nil {
		return "", fmt.Errorf("Unable to parse request body", err)
	}

	return secret.Auth.ClientToken, nil
}

func GetTokenTTL(token string) (time.Time, error) {

	var valid_until time.Time

	// Don't login, just show information about the current token.
	secret, err := vc.Auth().Token().Lookup(cfg.Token)
	if err != nil {
		return valid_until, err
	}

	ttl, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		return valid_until, err
	}

	return time.Unix(time.Now().Unix()+ttl, 0), nil
}
