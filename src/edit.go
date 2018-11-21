package main

import (
	"bufio"
	"errors"
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

var (
	ErrDuplicateKey       = errors.New("duplicate key found in secret")
	ErrMultipleDelimiters = errors.New("multiple \": \" delimiters found in secret")
	ErrMissingDelimiter   = errors.New("\": \" delimiter missing secret")
)

func (c *EditCommand) Run(args []string) int {

	if len(args) != 1 {
		c.Ui.Output("The edit command expects one argument")
		return 1
	}

	path := args[0]

	secret, err := kv.Read(path)
	if err != nil {
		c.Ui.Error(CheckError(err, err.Error()))
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

	var data map[string]interface{}
	secretIsValid := false

	// Re-open the text editor if the parsing of the resulting secret fails
	for(secretIsValid == false) {

		secretIsValid = true

		err = EditFile(file.Name())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Unable to edit secret file: %q", err))
			return 1
		}

		data, err = ParseSecret(file.Name())
		switch err {
		case ErrDuplicateKey, ErrMultipleDelimiters, ErrMissingDelimiter:
			secretIsValid = false
			// Let the user read the error in his shell before re-opening his editor to
			// correct the mistake
			_, _, _ = bufio.NewReader(os.Stdin).ReadLine()
		default:
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Secret has not changed %q", err))
				return 1
			}
		}
	}

	if len(data) == 0 {
		// Delete the secret if no key/value pairs are left
		_, err = kv.Delete(path)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to delete empty secret"))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Secret was deleted because no K/V pairs were associated with it."))
	} else {

		_, err = kv.Write(path, data)
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

// Validates and parses key/value pairs and comments from the temporary secret file
//
// Lines starting with "#" will be recognized as comments. Following lines that also start with "#"
// will be appended to the first.
//
// Lines not starting with "#" will be recognized as secrets. If the key identifier of a secret
// multiple times the user will get a chance to reedit the secret
//
func ParseSecret(path string) (map[string]interface{}, error) {

	var data = make(map[string]interface{})
	var comment string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		if line != "" {

			if strings.HasPrefix(line, "#") {

				// If a comment is alreay set, then assume that the comment spans
				// across multiple lines
				if comment != "" {
					comment += "\n" + strings.TrimPrefix(line, "#")
				} else {
					comment = strings.TrimPrefix(line, "#")
				}

			} else {

				kvPair := strings.Split(line, ": ")
				if len(kvPair) < 2 {
					fmt.Fprintf(os.Stderr, "Unable to parse key/value pair %q. Make sure that there is at least one \": \" delimiter in it ", line)
					return data, ErrMissingDelimiter

				} else if len(kvPair) > 2 {
					fmt.Fprintf(os.Stderr, "Unable to parse key/value pair %q. Make sure that there is only one \": \" delimiter in it.", line)
					return data, ErrMultipleDelimiters
				}


				key, value := kvPair[0], kvPair[1]

				// Check whether the previous lines have been parsed as comment. If
				// thats case then compose a key/value pair with a unique identifier
				// by adding a suffix.
				if comment != "" {
					data[key+"_comment"] = comment
					comment = ""
				}

				// Check that key is not used multiple times
				if _, already_used := data[key]; already_used {
					fmt.Fprintf(os.Stderr, "Secret identifier %q is used multiple times. Please make sure that the key only is used once.", key)
					return data, ErrDuplicateKey

				} else {
					data[key] = value
				}

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
