package main

import (
	"io/ioutil"
	"os/user"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Token string `yaml:"token"`
}

func LoadConfig() error {

	// Configuration default values
	cfg = Config{
		Host:  "127.0.0.1",
		Port:  8200,
		Token: "password",
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(usr.HomeDir + "/.vaultrc")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	return nil
}
