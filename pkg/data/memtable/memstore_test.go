package memtable_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/data/memtable"
)

func TestMemstoreBasic(t *testing.T) {
	store := memtable.NewMemstore()
	if err := store.Set("key1", []byte("val1")); err != nil {
		t.Fatal(err)
	}
	if err := store.Set("key2", []byte("val2")); err != nil {
		t.Fatal(err)
	}

	val, err := store.Get("key1")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("val1"), val) {
		t.Fatalf(`expected 'val1' but got '%s'`, val)
	}

	val, err = store.Get("key2")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("val2"), val) {
		t.Fatalf(`expected 'val2' but got '%s'`, val)
	}
}

func assertArraysEqual(t *testing.T, expected [][]byte, actual [][]byte) {
	if len(expected) != len(actual) {
		t.Fatalf("arrays not equal, expected is len %d, while actual is len %d", len(expected), len(actual))
	}
	for i, e := range expected {
		if !bytes.Equal(e, actual[i]) {
			t.Fatalf("arrays not equal, index %d has expected value '%s' and actual value '%s'", i, e, actual[i])
		}
	}
}

func TestMemstoreRanges(t *testing.T) {
	store := memtable.NewMemstore()
	values := [][]byte{}
	// generates { 0 : 0, 2: 1, 4: 2 }
	for i := 0; i < 3; i++ {
		key := strconv.Itoa(i * 2)
		value := []byte(strconv.Itoa(i))
		err := store.Set(key, value)
		if err != nil {
			t.Fatal(err)
		}
		values = append(values, value)
	}

	beforeStart := ""
	atStart := "0"
	afterStart := "1"

	beforeEnd := "3"
	atEnd := "4"
	afterEnd := "5"

	for _, test := range []struct {
		start    string
		end      string
		expected [][]byte
	}{
		{
			start:    beforeStart,
			end:      afterEnd,
			expected: values,
		},
		{
			start:    beforeStart,
			end:      atEnd,
			expected: values[:2],
		},
		{
			start:    beforeStart,
			end:      beforeEnd,
			expected: values[:2],
		},
		{
			start:    atStart,
			end:      afterEnd,
			expected: values,
		},
		{
			start:    atStart,
			end:      atEnd,
			expected: values[:2],
		},
		{
			start:    atStart,
			end:      beforeEnd,
			expected: values[:2],
		},
		{
			start:    afterStart,
			end:      afterEnd,
			expected: values[1:],
		},
		{
			start:    afterStart,
			end:      atEnd,
			expected: values[1:2],
		},
		{
			start:    afterStart,
			end:      beforeEnd,
			expected: values[1:2],
		},
	} {
		t.Run(fmt.Sprintf("start %s, end %s", test.start, test.end), func(t *testing.T) {
			cur, err := store.GetRange(test.start, test.end)
			if err != nil {
				t.Fatal(err)
			}
			result, err := cur.Read(100)
			if err != nil {
				t.Fatal(err)
			}
			assertArraysEqual(t, test.expected, result)
		})
	}
}
