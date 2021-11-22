package dbutils

import (
	"testing"
)

func TestMongoConnection(t *testing.T) {
	m := InitMongoInstance()
	defer m.Close()
	err := m.Ping()
	if err != nil {
		t.Errorf("Ping returned an error: %s", err)
	}
}
