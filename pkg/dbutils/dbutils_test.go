package dbutils

import (
	"testing"
)

func TestTest_Me(t *testing.T) {
	result := Test_Me()
	want := "tests are cool!"
	if result != want {
		t.Errorf("Test_Me returned %s, we want %s", result, want)
	}
}

func TestMongoConnection(t *testing.T) {
	m := InitMongoInstance()
	defer m.close()
	err := m.ping()
	if err != nil {
		t.Errorf("Ping returned an error: %s", err)
	}
}
