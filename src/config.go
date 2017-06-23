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
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Token       string `yaml:"token"`
	TLS         bool   `yaml:"tls"`
	VerifyTLS   bool   `yaml:"verify_tls"`
	AuthMethod  string `yaml:"auth_method"`
	AuthBackend string `yaml:"auth_backend"`
}

func LoadConfig() error {

	cfg = Config{
		Host:        "127.0.0.1",
		Port:        8200,
		Token:       "password",
		TLS:         true,
		VerifyTLS:   true,
		AuthMethod:  "token",
		AuthBackend: "token",
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	path := usr.HomeDir + "/.vaultrc"

	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	config_file_permissions := file.Mode().String()

	// Check that the config file is only readable by the user.
	// And not by his group or others (-rwx------)
	if !strings.HasSuffix(config_file_permissions, "------") {
		return fmt.Errorf("Your ~/.vaultrc is accessible for others (chmod 700 ~/.vaultrc)")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func ComposeUrl() string {

	protocol := "http"
	if cfg.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%v://%v:%v", protocol, cfg.Host, cfg.Port)

}
