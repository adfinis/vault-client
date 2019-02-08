package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestKvClientV1(t *testing.T) {

	var err error
	cfg, vc, err = SetupTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	t.Run("KvClientPut", func(t *testing.T) {

		key := TestBackend + "/insertedsecret"
		value := map[string]interface{}{"password": "test1234"}
		err := kv.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}

		secret, err := kv.GetRawClient().Logical().Read(key)
		if err != nil {
			t.Fatal(err)
		}

		expected := value
		actual := secret.Data
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected:\n%q\n\nactual:\n%q", expected, actual)
		}
	})

	err = TeardownTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
