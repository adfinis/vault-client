package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
	"time"
)

type EditCommand struct {
	Ui cli.Ui
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

	lock, err := LockSecret(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to create the lock: %v\n", err))
		return 1
	}

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

	} else {
		data = secret.Data
	}

	file, err := ioutil.TempFile("", "vaultsecret")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to create tempfile: %v\n", err))
		return 1
	}

	defer os.Remove(file.Name())

	WriteSecretToFile(data, file)

	err = EditFileWithEditor(file.Name())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to edit file with editor: %v\n", err))
		return 1
	}

	editedData, err := ParseSecretFromFile(file)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Secret has not changed: %v\n.", err))
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

	// Release lock so other people can edited the secret
	err = lock.Unlock()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("%v\nUnable to remove the lock", err))
		return 1
	}

	return 0
}

// Return a Lock to the secret that the Users wants to edit
func LockSecret(path string) (*consul.Lock, error) {

	var lock *consul.Lock

	user, err := user.Current()
	if err != nil {
		return lock, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return lock, err
	}

	lockKey := cfg.Consul.LockKVRoot + path
	lockValue := []byte(fmt.Sprintf("%v@%v", user.Name, hostname))

	lockOpts := &consul.LockOptions{
		Key:          lockKey,
		Value:        lockValue,
		LockWaitTime: 2 * time.Second,
		LockTryOnce:  true,
		SessionOpts: &consul.SessionEntry{
			LockDelay: 1,
			TTL:       "600s",
		},
	}

	lock, err = cc.LockOpts(lockOpts)
	if err != nil {
		return lock, err
	}

	stpCh, err := lock.Lock(nil)
	if err != nil {
		return lock, err
	}

	if stpCh == nil {
		kv := cc.KV()
		kvpair, _, err := kv.Get(lockKey, nil)
		if err != nil {
			return lock, err
		}

		return lock, fmt.Errorf("Secret is already locked by %v", string(kvpair.Value))
	}

	return lock, nil

}

// Write k/v pairs of secret to a file
func WriteSecretToFile(data map[string]interface{}, file *os.File) {

	// Sort secrets lexicographically
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Write secrets to tempfile in sorted order
	for _, k := range keys {
		file.WriteString(k + ": " + data[k].(string) + "\n")
	}

}

// Edit a file with the editor specified in $EDITOR. If $EDITOR is not defined vi will be used as
// fallback.
func EditFileWithEditor(path string) error {

	var cmdstr []string

	editor := os.Getenv("EDITOR")
	if editor == "" {
		cmdstr = append(cmdstr, "vi")
	} else {
		// If $EDITOR has arguments (e.g. "emacs -nw") split them up
		cmdstr = strings.Split(editor, " ")
	}

	cmdstr = append(cmdstr, path)

	cmd := exec.Command(cmdstr[0], cmdstr[1:len(cmdstr)]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Parse secrets from file
func ParseSecretFromFile(file *os.File) (map[string]interface{}, error) {

	data := make(map[string]interface{})

	// Parse the file from the beginning
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		kv_pair := strings.Split(line, ": ")

		if len(kv_pair) == 2 {
			data[kv_pair[0]] = kv_pair[1]
		} else {
			return nil, fmt.Errorf("Unable to parse key/value pair: %q", line)
		}
	}

	return data, nil
}
