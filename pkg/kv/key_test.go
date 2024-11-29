package kv_test

import (
	"testing"
	"unicode/utf8"

	"github.com/angles-n-daemons/popsql/pkg/kv"
)

func TestNewKey(t *testing.T) {
	table := "testTable"
	key := kv.NewKey(table)
	if key.Table != table {
		t.Errorf("expected table %s, got %s", table, key.Table)
	}
	if key.ID != "" {
		t.Errorf("expected empty ID, got %s", key.ID)
	}
}

func TestWithTable(t *testing.T) {
	key := kv.NewKey("testTable")
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
	key := kv.NewKey("testTable")
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
	key := kv.NewKey("testTable").WithID("123")
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
	key := kv.NewKey("testTable")
	if key.String() != "testTable/" {
		t.Errorf("expected string %s, got %s", "testTable", key.String())
	}
	key = key.WithID("testID")
	if key.String() != "testTable/testID" {
		t.Errorf("expected string %s, got %s", "testTable/testID", key.String())
	}

	// special case check for the end string
	if key.WithID(string(kv.END_ID)).String() != "testTable/<END>" {
		t.Errorf("expected string %s, got %s", "testTable/<END>", key.String())
	}
}

func TestNext(t *testing.T) {
	mr := string(utf8.MaxRune)
	end := string(kv.END_ID)
	tests := []struct {
		table, id, expectedTable, expectedID string
	}{
		{"table", end, "table", end},
		{"table", "", "table", string(kv.END_ID)},
		{"table", "id", "table", "ie"},
		{"table", "id" + mr, "table", "id" + mr + " "},
	}

	for _, test := range tests {
		key := kv.NewKey(test.table).WithID(test.id)
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
	mr := string(utf8.MaxRune)
	tests := []struct {
		input, expected string
	}{
		{"", " "},
		{" ", "!"},
		{"a", "b"},
		{mr, mr + " "},
		{" " + mr, " " + mr + " "},
		{mr + mr, mr + mr + " "},
	}

	for _, test := range tests {
		result := kv.NextString(test.input)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}
