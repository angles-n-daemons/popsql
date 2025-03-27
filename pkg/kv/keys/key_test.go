package keys_test

import (
	"testing"
	"unicode/utf8"

	"github.com/angles-n-daemons/popsql/pkg/kv/keys"
)

func TestNew(t *testing.T) {
	table := "testTable"
	key := keys.New(table)
	if key.Table != table {
		t.Errorf("expected table %s, got %s", table, key.Table)
	}
	if key.ID != "" {
		t.Errorf("expected empty ID, got %s", key.ID)
	}
}

func TestWithTable(t *testing.T) {
	key := keys.New("testTable")
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
	key := keys.New("testTable")
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
	key := keys.New("testTable").WithID("123")
	addition := "123"
	newKey := key.WithIDAddition(addition)
	if newKey.ID != key.ID+addition {
		t.Errorf("expected ID %s, got %s", key.ID+addition, newKey.ID)
	}
	if newKey.Table != key.Table {
		t.Errorf("expected table %s, got %s", key.Table, newKey.Table)
	}
}

func TestEncode(t *testing.T) {
	key := keys.New("testTable")
	if key.Encode() != "testTable/" {
		t.Errorf("expected string %s, got %s", "testTable", key.Encode())
	}
	key = key.WithID("testID")
	if key.Encode() != "testTable/testID" {
		t.Errorf("expected string %s, got %s", "testTable/testID", key.Encode())
	}

	// special case check for the end string
	if key.WithID(string(keys.END_ID)).Encode() != "testTable?" {
		t.Errorf("expected string %s, got %s", "testTable?", key.Encode())
	}
}

func TestNext(t *testing.T) {
	mr := string(utf8.MaxRune)
	end := string(keys.END_ID)
	tests := []struct {
		table, id, expectedTable, expectedID string
	}{
		{"table", end, "table", end},
		{"table", "", "table", string(keys.END_ID)},
		{"table", "id", "table", "ie"},
		{"table", "id" + mr, "table", "id" + mr + " "},
	}

	for _, test := range tests {
		key := keys.New(test.table).WithID(test.id)
		nextKey := key.Next()
		if nextKey.Table != test.expectedTable {
			t.Errorf("expected table %s, got %s", test.expectedTable, nextKey.Table)
		}
		if nextKey.ID != test.expectedID {
			t.Errorf("expected ID %s, got %s", test.expectedID, nextKey.ID)
		}
	}
}

func TestNextEncode(t *testing.T) {
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
		result := keys.NextString(test.input)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}
