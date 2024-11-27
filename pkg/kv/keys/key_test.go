package keys_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
)

func TestNewKey(t *testing.T) {
	table := "testTable"
	key := keys.NewKey(table)
	if key.Table != table {
		t.Errorf("expected table %s, got %s", table, key.Table)
	}
	if key.ID != "" {
		t.Errorf("expected empty ID, got %s", key.ID)
	}
}

func TestWithTable(t *testing.T) {
	key := keys.NewKey("testTable")
	newTable := "newTable"
	newKey := key.WithTable(newTable)
	if newKey.Table != newTable {
		t.Errorf("expected table %s, got %s", newTable, newKey.Table)
	}
	if newKey.ID != key.ID {
		t.Errorf("expected ID %s, got %s", key.ID, newKey.ID)
	}
}

func TestWithID(t *testing.T) {
	key := keys.NewKey("testTable")
	newID := "newID"
	newKey := key.WithID(newID)
	if newKey.ID != newID {
		t.Errorf("expected ID %s, got %s", newID, newKey.ID)
	}
	if newKey.Table != key.Table {
		t.Errorf("expected table %s, got %s", key.Table, newKey.Table)
	}
}

func TestWithIDAddition(t *testing.T) {
	key := keys.NewKey("testTable")
	addition := "123"
	newKey := key.WithIDAddition(addition)
	if newKey.ID != key.ID+addition {
		t.Errorf("expected ID %s, got %s", key.ID+addition, newKey.ID)
	}
	if newKey.Table != key.Table {
		t.Errorf("expected table %s, got %s", key.Table, newKey.Table)
	}
}

func TestString(t *testing.T) {
	key := keys.NewKey("testTable")
	if key.String() != "testTable" {
		t.Errorf("expected string %s, got %s", "testTable", key.String())
	}
	key = key.WithID("testID")
	if key.String() != "testTable/testID" {
		t.Errorf("expected string %s, got %s", "testTable/testID", key.String())
	}
}

func TestNext(t *testing.T) {
	tests := []struct {
		table, id, expectedTable, expectedID string
	}{
		{"table", "", "tablf", ""},
		{"table", "id", "table", "ie"},
		{"table", "idz", "table", "idza"},
		{"tablez", "", "tableza", ""},
	}

	for _, test := range tests {
		key := keys.NewKey(test.table).WithID(test.id)
		nextKey := key.Next()
		if nextKey.Table != test.expectedTable {
			t.Errorf("expected table %s, got %s", test.expectedTable, nextKey.Table)
		}
		if nextKey.ID != test.expectedID {
			t.Errorf("expected ID %s, got %s", test.expectedID, nextKey.ID)
		}
	}
}

func TestNextString(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"a", "b"},
		{"z", "za"},
		{"az", "aza"},
		{"zz", "zza"},
	}

	for _, test := range tests {
		result := keys.NextString(test.input)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}
