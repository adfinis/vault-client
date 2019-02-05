package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestEdit(t *testing.T) {

	_, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	t.Run("ParseInvalidSecretFile", func(t *testing.T) {

		test_files := map[string]error{
			"invalid_secret_multiple_delimiters.txt": ErrMultipleDelimiters,
			"invalid_secret_missing_delimiter.txt":   ErrMissingDelimiter,
			"invalid_secret_duplicated_key.txt":      ErrDuplicateKey,
		}

		for file, expected := range test_files {

			_, actual := ParseSecret("../sample/tests/secrets/" + file)

			if actual != expected {
				t.Fatalf("\nexpected:\t%v\nto include:\t%v", actual, expected)
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
				t.Fatalf("\nexpected:\t%s\nto include:\t%s", actual, expected)
				t.Fatalf("\nexpected:\n%q\n\nto include:\n%q", actual, expected)
			}
		}
	})
}
