package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

type EditCommand struct {
	Ui cli.Ui
}

func (c *EditCommand) Run(args []string) int {

	if len(args) != 1 {
		c.Ui.Output("The edit command expects one argument")
		return 1
	}

	path := args[0]

	secret, err := vc.Logical().Read(path)
	if err != nil {
		return 1
	}

	file, err := ioutil.TempFile("", "vaultsecret")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to create temporary secret file %q", err))
		return 1
	}
	defer os.Remove(file.Name())

	if secret == nil {
		// If the secret does not exist, it will not have any data. In that case initialize
		// it to avoid a nil pointer exception
		secret = &vault.Secret{Data: make(map[string]interface{})}
	}

	WriteSecretToFile(file, secret.Data)

	err = EditFile(file.Name())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to edit secret file %q", err))
		return 1
	}

	data, err := ParseSecret(file.Name())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Secret has not changed %q", err))
		return 1
	}

	if len(data) == 0 {
		// Delete the secret if no key/value pairs are left
		_, err = vc.Logical().Delete(path)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to delete empty secret"))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Secret was deleted because no K/V pairs were associated with it."))
	} else {
		_, err = vc.Logical().Write(path, data)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to save secret %q", err))
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

// Parses k/v pairs and comments from a secret file
func ParseSecret(path string) (map[string]interface{}, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	var data = make(map[string]interface{})
	var comment string

	for scanner.Scan() {

		line := scanner.Text()

		if line != "" {
			if strings.HasPrefix(line, "#") {
				// If line is a comment, store it's value for later composition into it's
				// own k/v pair
				if comment != "" {
					// If a comment is alreay set, then assume that the comment spans
					// across multiple lines
					comment += "\n" + strings.TrimPrefix(line, "#")
				} else {
					comment = strings.TrimPrefix(line, "#")
				}

			} else {
				// If a line is a k/v pair then split it up
				kv_pair := strings.Split(line, ": ")

				// Check whether a related comment to this secret was stored
				if comment != "" {
					data[kv_pair[0]+"_comment"] = comment
					comment = ""
				}

				if len(kv_pair) != 2 {
					return nil, fmt.Errorf("Unable to parse key/value pair: %q", line)
				}
				data[kv_pair[0]] = kv_pair[1]
			}
		}
	}

	return data, nil
}

func WriteSecretToFile(file *os.File, kv_pairs map[string]interface{}) {

	// Sort secrets lexicographically
	var keys []string
	for key := range kv_pairs {
		// Ignore comments
		if !strings.HasSuffix(key, "_comment") {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	for _, key := range keys {
		// Write comment right before the related k/v pair
		if value, exists := kv_pairs[key+"_comment"].(string); exists {

			if multilineComments := strings.Split(value, "\n"); len(multilineComments) > 1 {
				for _, comment := range multilineComments {
					file.WriteString("#" + comment + "\n")
				}
			} else {
				file.WriteString("#" + value + "\n")
			}
		}
		file.WriteString(key + ": " + kv_pairs[key].(string) + "\n")
	}
}
