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

		key := TestBackend + "/putsecret"
		value := map[string]interface{}{"password": "test1234"}
		err := kv.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}

		secret, err := kv.RawClient().Logical().Read(key)
		if err != nil {
			t.Fatal(err)
		}

		expected := value
		actual := secret.Data
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected:\n%q\n\nactual:\n%q", expected, actual)
		}
	})

	t.Run("KvClientGet", func(t *testing.T) {

		key := TestBackend + "/getsecret"
		value := map[string]interface{}{"password": "test1234"}

		_, err := kv.RawClient().Logical().Write(key, value)
		if err != nil {
			t.Fatal(err)
		}

		kvPairs, err := kv.Get(key)
		if err != nil {
			t.Fatal(err)
		}

		expected := value
		actual := kvPairs
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected:\n%q\n\nactual:\n%q", expected, actual)
		}
	})

	t.Run("KvClientDelete", func(t *testing.T) {

		key := TestBackend + "/deletesecret"
		value := map[string]interface{}{"password": "test1234"}
		_, err := kv.RawClient().Logical().Write(key, value)
		if err != nil {
			t.Fatal(err)
		}

		if kv.Delete(key) != nil {
			t.Fatal()
		}

		secret, err := kv.RawClient().Logical().Read(key)
		if secret != nil {
			t.Fatal()
		}
	})

	t.Run("KvClientList", func(t *testing.T) {

		keys := []string{
			TestBackend + "/list/listsecret1",
			TestBackend + "/list/listsecret2",
			TestBackend + "/list/nested/listsecret3",
		}
		value := map[string]interface{}{"password": "test1234"}
		for _, key := range keys {
			_, err := kv.RawClient().Logical().Write(key, value)
			if err != nil {
				t.Fatal(err)
			}
		}

		expected := []string{"listsecret1", "listsecret2", "nested/"}

		actual, err := kv.List(TestBackend + "/list/")
		if err != nil {
			t.Fatal()
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected:\n%q\n\nactual:\n%q", expected, actual)
		}
	})

	err = TeardownTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
