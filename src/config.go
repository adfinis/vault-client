package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Token     string `yaml:"token"`
	TLS       bool   `yaml:"tls"`
	VerifyTLS bool   `yaml:"verify_tls"`
	// Name of authentication method
	AuthMethod string `yaml:"auth_method"`
	// Type of the authentication backend
	AuthBackend string `yaml:"auth_backend"`
	Path        string
}

func LoadConfig() (Config, error) {

	cfg = Config{
		Host:        "127.0.0.1",
		Port:        8200,
		Token:       "password",
		TLS:         true,
		VerifyTLS:   true,
		AuthMethod:  "token",
		AuthBackend: "token",
	}

	var err error

	cfg.Path, err = GetConfigPath()
	if err != nil {
		return cfg, err
	}

	file, err := os.Stat(cfg.Path)
	if err != nil {
		return cfg, err
	}

	config_file_permissions := file.Mode().String()

	// Check that the config file is only readable by the user.
	// And not by his group or others (-rwx------)
	if !strings.HasSuffix(config_file_permissions, "------") {
		return cfg, fmt.Errorf("Your ~/.vaultrc is accessible for others (chmod 700 ~/.vaultrc)")
	}

	content, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func ComposeUrl() string {

	protocol := "http"
	if cfg.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%v://%v:%v", protocol, cfg.Host, cfg.Port)

}

// Update the token in the configuration file
func UpdateConfigToken(token string) error {

	// Reauthenticate against Vault and update in-memory config
	vc.SetToken(token)
	vc.Auth()
	cfg.Token = token

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	token_found := false

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "token:") {
			lines[i] = "token: " + token
			token_found = true
		}
	}

	if !token_found {
		lines = append(lines, "token: "+token)

	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0600)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigPath() (string, error) {

	path := os.Getenv("VAULT_CLIENT_CONFIG")

	if path != "" {
		return path, nil
	} else {

		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		return usr.HomeDir + "/.vaultrc", nil
	}
}
