package memtable_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/data/memtable"
)

func TestMemcursorRead(t *testing.T) {
	store := memtable.NewMemstore()
	values := [][]byte{}
	// generates { 0 : a, 2: b, 4: c }
	for i := 0; i < 3; i++ {
		key := strconv.Itoa(i * 2)
		value := []byte{byte('a' + i)}
		err := store.Set(key, value)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(key, value)
		values = append(values, value)
	}

	afterEnd := "5"
	atEnd := "4"
	beforeEnd := "3"
	for _, test := range []struct {
		end      string
		num      int
		isAtEnd  bool
		expected [][]byte
	}{
		{end: afterEnd, num: 3, isAtEnd: true, expected: values},
		{end: atEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{end: beforeEnd, num: 3, isAtEnd: true, expected: values[:2]},
		{end: afterEnd, num: 2, isAtEnd: false, expected: values[:2]},
		{end: atEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{end: beforeEnd, num: 2, isAtEnd: true, expected: values[:2]},
		{end: afterEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{end: atEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{end: beforeEnd, num: 1, isAtEnd: false, expected: values[:1]},
		{end: afterEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{end: atEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
		{end: beforeEnd, num: 0, isAtEnd: false, expected: [][]byte{}},
	} {
		t.Run(fmt.Sprintf("end=%s, num=%d", test.end, test.num), func(t *testing.T) {
			cur, err := store.GetRange("", test.end)
			if err != nil {
				t.Fatal(err)
			}
			result, err := cur.Read(test.num)
			if err != nil {
				t.Fatal(err)
			}
			if !test.isAtEnd == cur.IsAtEnd() {
				t.Fatalf("expected isAtEnd call to be %t, but got %t", test.isAtEnd, cur.IsAtEnd())
			}
			assertArraysEqual(t, test.expected, result)
		})
	}
}

func TestMemcursorMultipleRead(t *testing.T) {
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

	actual := [][]byte{}
	cur, err := store.GetRange("", "9999")
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
