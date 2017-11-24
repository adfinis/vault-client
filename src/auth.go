package main

import (
	"crypto/tls"
	"strings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

// Authenticate against Vault using the configured method/auth backend and
// return a valid token
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.VerifyTLS},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve authentication token from vault %q", err)
	} else if res.StatusCode != 200 {
		return "", fmt.Errorf("Unable to retrieve authentication token from vault (status code %d)", res.StatusCode)
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

// Returns the the time when a token will expire
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


// Check whether a generic error occured because:
//   1. The token expired
//   2. Vault is down
// Otherwise return an alternative text that got passed by the caller
func CheckError(err error, alternate_text string) string {

	switch true {
	case strings.HasSuffix(err.Error(), "getsockopt: connection refused"):
		return "Unable to connect to Vault"
	case strings.HasPrefix(err.Error(), "Error making API request."):
		return "Your token has expired. Please reauthenticate."
	case strings.HasSuffix(err.Error(), "request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)"):
		return "Connection to Vault has timed out"
	default:
		return fmt.Sprintf("Unkown error occured: %q", err.Error())
	}
}
