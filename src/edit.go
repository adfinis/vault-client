package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
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
		c.Ui.Error(fmt.Sprintf("%v\nSecret has not changed.", err))
		return 1
	}

	if len(editedData) == 0 {
		// Delete the secret if no key/value pairs are left
		_, err = vc.Logical().Delete(path)
		if err != nil {
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Secret was deleted because no K/V pairs were associated with it."))
	} else {
		_, err = vc.Logical().Write(path, editedData)
		if err != nil {
			return 1
		}
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

	f, err := ioutil.TempFile("", "vaultsecret")
	if err != nil {
		return nil, err
	}

	defer os.Remove(f.Name())

	// Sort secrets lexicographically
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Write secrets to tempfile
	for _, k := range keys {
		f.WriteString(k + ": " + data[k].(string) + "\n")
	}
	f.Close()

	// Edit temporary file
	err = EditFile(f.Name())
	if err != nil {
		return nil, err
	}

	// Parse secret
	parsedData := make(map[string]interface{})
	editedFile, err := os.Open(f.Name())
	scanner := bufio.NewScanner(editedFile)

	for scanner.Scan() {
		line := scanner.Text()
		kv_pair := strings.Split(line, ": ")
		if len(kv_pair) == 2 {
			parsedData[kv_pair[0]] = kv_pair[1]
		} else {
			return nil, fmt.Errorf("Unable to parse key/value pair: %q", line)
		}
	}

	return parsedData, nil
}

// Edit a file with the editor specified in $EDITOR or vi as fallback
func EditFile(path string) error {

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
		return err
	}

	return nil
}
