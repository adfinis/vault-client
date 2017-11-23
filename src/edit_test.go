package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestEdit(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	t.Run("ParseInvalidSecretFile", func(t *testing.T) {

		test_files := map[string]string{
			"invalid_secret1.txt": "Unable to parse key/value pair \"valid_key: invalid: _value\". Make sure that there is only/at least one \":\" in it",
			"invalid_secret2.txt": "Unable to parse key/value pair \"invalid_line\". Make sure that there is only/at least one \":\" in it",
		}

		for file, expected := range test_files {

			_, actual := ParseSecret("../sample/tests/secrets/" + file)

			if !strings.Contains(actual.Error(), expected) {
				t.Fatalf("\nexpected:\t%s\nto include:\t%s", actual, expected)
			}
		}
	})

	t.Run("ParseValidSecretFileWithComments", func(t *testing.T) {

		test_files := map[string]map[string]interface{}{
			"valid_secret.txt": {
				"valid_key1":         "valid_value",
				"valid_key1_comment": " This is a valid comment",
				"valid_key2":         "valid_value",
				"valid_key2_comment": " Multiline\n comment",
			},
		}

		for file, expected := range test_files {

			actual, err := ParseSecret("../sample/tests/secrets/" + file)
			if err != nil {
				t.Fatalf("Unable to parse secret file: %q", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("\nexpected:\n%q\n\nto include:\n%q", expected, actual)
			}
		}
	})
}
