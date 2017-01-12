package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
)

type EditCommand struct {
	Ui cli.Ui
}

func (c *EditCommand) Run(args []string) int {

	switch {
	case len(args) > 1:
		c.Ui.Output("The edit command expects at most one argument")
		return 1
	case len(args) == 0:
		c.Ui.Output("The edit command expects an argument")
		return 1
	}

	path := args[0]

	secret, err := vc.Logical().Read(path)
	if err != nil {
		return 1
	}

	data := make(map[string]interface{})

	if secret == nil {
		answer, err := c.Ui.Ask("Secret doesn't exist. Would you like to create it? [Yn]")
		if err != nil {
			return 1
		}

		if answer := strings.ToLower(answer); answer == "n" {
			return 0
		}
		data["key"] = "value"
	} else {
		data = secret.Data
	}

	editedData, err := ProcessSecret(data)
	if err != nil {
		c.Ui.Output(fmt.Sprintf("%v\nSecret has not changed.", err))
		return 1
	}

	_, err = vc.Logical().Write(path, editedData)
	if err != nil {
		return 1
	}

	return 0
}

func (c *EditCommand) Help() string {
	return `Usage: vc edit path

  This command edits a secret at a certain path with your editor of choice
  (set through $EDITOR). If no editor is specified vi will be used as fallback.
`
}

func (c *EditCommand) Synopsis() string {
	return "Edit a secret at specified path"
}

// Processes a secret by unmarshaling and writting it into a tempfile.
// After the file was edit it will reread the tempfile marhsal the data and clean up.
func ProcessSecret(data map[string]interface{}) (map[string]interface{}, error) {

	parsedData := make(map[string]interface{})

	f, err := ioutil.TempFile("", "vaultsecret")
	if err != nil {
		return nil, err
	}

	defer os.Remove(f.Name())

	ymldata, err := yaml.Marshal(&data)
	_, err = f.Write(ymldata)
	if err != nil {
		return nil, err
	}

	editedData, err := EditFile(f.Name())
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(editedData, parsedData)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse yaml tempfile: %q", err)
	}

	err = ValidateData(parsedData)
	if err != nil {
		return nil, err
	}

	return parsedData, nil

}

// Edit a file with the editor specified in $EDITOR or vi as fallback
func EditFile(path string) ([]byte, error) {

	var cmdstring []string

	editor := os.Getenv("EDITOR")
	if editor == "" {
		cmdstring = append(cmdstring, "vi")
	} else {
		cmdstring = strings.Split(editor, " ")
	}

	cmdstring = append(cmdstring, path)
	_ = cmdstring

	cmd := exec.Command(cmdstring[0], cmdstring[1:len(cmdstring)]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Check if data keys of a secrets contain only valid characters
func ValidateData(data map[string]interface{}) error {

	allowedCharacters := "^[A-Za-z0-9-_]*$"

	for k := range data {
		matched, err := regexp.MatchString(allowedCharacters, k)
		if err != nil {
			return fmt.Errorf("Unable to validate secret keys: %q", err)
		}

		if !matched {
			return fmt.Errorf("Invalid characters in key %q", k)
		}

	}
	return nil
}
