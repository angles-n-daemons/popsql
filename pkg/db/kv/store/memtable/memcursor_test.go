package memtable_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/kv/store/memtable"
)

func TestMemcursorRead(t *testing.T) {
	store := memtable.NewStore()
	values := [][]byte{}
	// generates { 2 : a, 4: b, 6: c }
	for i := 1; i < 4; i++ {
		key := strconv.Itoa(i * 2)
		value := []byte{byte('a' + i)}
		err := store.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}
		values = append(values, value)
	}

	afterStart := "3"
	atStart := "2"
	beforeStart := "1"
	afterEnd := "7"
	atEnd := "6"
	beforeEnd := "5"
	for _, test := range []struct {
		start    string
		end      string
		num      int
		isAtEnd  bool
		expected [][]byte
	}{
		{start: beforeStart, end: afterEnd, num: 3, isAtEnd: true, expected: values},
		{start: beforeStart, end: atEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{start: beforeStart, end: beforeEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{start: beforeStart, end: afterEnd, num: 2, isAtEnd: false, expected: values[:2]},
		{start: beforeStart, end: atEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{start: beforeStart, end: beforeEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{start: beforeStart, end: afterEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: beforeStart, end: atEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: beforeStart, end: beforeEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: beforeStart, end: afterEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: beforeStart, end: atEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: beforeStart, end: beforeEnd, num: 0, isAtEnd: false, expected: [][]byte{}},

		{start: atStart, end: afterEnd, num: 3, isAtEnd: true, expected: values},
		{start: atStart, end: atEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{start: atStart, end: beforeEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{start: atStart, end: afterEnd, num: 2, isAtEnd: false, expected: values[:2]},
		{start: atStart, end: atEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{start: atStart, end: beforeEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{start: atStart, end: afterEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: atStart, end: atEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: atStart, end: beforeEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{start: atStart, end: afterEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: atStart, end: atEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: atStart, end: beforeEnd, num: 0, isAtEnd: false, expected: [][]byte{}},

		{start: afterStart, end: afterEnd, num: 3, isAtEnd: true, expected: values[1:]},
		{start: afterStart, end: atEnd, num: 3, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: beforeEnd, num: 3, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: afterEnd, num: 2, isAtEnd: true, expected: values[1:]},
		{start: afterStart, end: atEnd, num: 2, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: beforeEnd, num: 2, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: afterEnd, num: 1, isAtEnd: false, expected: values[1:2]},
		{start: afterStart, end: atEnd, num: 1, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: beforeEnd, num: 1, isAtEnd: true, expected: values[1:2]},
		{start: afterStart, end: afterEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: afterStart, end: atEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{start: afterStart, end: beforeEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
	} {
		t.Run(fmt.Sprintf("start %s, end=%s, num=%d", test.start, test.end, test.num), func(t *testing.T) {
			cur, err := store.Scan(test.start, test.end)
			if err != nil {
				t.Fatal(err)
			}
			result, err := cur.Read(test.num)
			if err != nil {
				t.Fatal(err)
			}
			if test.isAtEnd != cur.IsAtEnd() {
				t.Fatalf("expected isAtEnd call to be %t, but got %t", test.isAtEnd, cur.IsAtEnd())
			}
			assertArraysEqual(t, test.expected, result)
		})
	}
}

func TestMemcursorMultipleRead(t *testing.T) {
	store := memtable.NewStore()
	values := [][]byte{}
	// generates { 0 : 0, 2: 1, 4: 2 }
	for i := 0; i < 3; i++ {
		key := strconv.Itoa(i * 2)
		value := []byte{byte('a' + i)}
		err := store.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}
		values = append(values, value)
	}

	actual := [][]byte{}
	cur, err := store.Scan("", "9999")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		result, err := cur.Read(1)
		if err != nil {
			t.Fatal(err)
		}
		if i < 2 && cur.IsAtEnd() {
			t.Fatalf("didnt expect to be at the end at i=%d", i)
		}
		if i >= 2 && !cur.IsAtEnd() {
			t.Fatalf("expected to be at end at i=%d", i)
		}

		for _, elem := range result {
			actual = append(actual, elem)
		}
	}
	assertArraysEqual(t, values, actual)
}
