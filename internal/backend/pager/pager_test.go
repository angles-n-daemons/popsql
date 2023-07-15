package pager

import (
	"bytes"
	"errors"
	"testing"
)

type mockFile struct {
	index   int
	content []byte
	err     error
}

func (m *mockFile) Read(i int, length int) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.content[i : i+length], nil
}

func (m *mockFile) Write(i int, content []byte) error {
	if m.err != nil {
		return m.err
	}

	m.content = []byte(
		string(m.content[:i]) + string(content) + string(m.content[i+len(content):]),
	)
	return m.err

}

const (
	testPageSize = 2
	numbers      = "0123456789"
)

func TestPagerGet(t *testing.T) {
	tests := []struct {
		name         string
		pageNumber   int
		expected     []byte
		fsError      error
		errors bool
	}{
		{
			name:       "read zeroth page",
			pageNumber: 0,
			expected:   []byte("23"),
		},
		{
			name:       "read second page",
			pageNumber: 2,
			expected:   []byte("67"),
		},

		// errors
		{
			name:         "read second page with fs error",
			pageNumber:   2,
			fsError:      errors.New("error"),
			errors: true,
		},
		{
			name:         "read negative page",
			pageNumber:   -2,
			errors: true,
		},
		{
			name:         "read past freelist",
			pageNumber:   4,
			errors: true,
		},
		{
			name:         "read page at freelist start",
			pageNumber:   3,
			errors: true,
		},
	}

	pager := Pager{
		pageSize:      2,
		freelistIndex: 3,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set the fs on the pager
			pager.fs = &mockFile{
				content: []byte(numbers),
				err:     test.fsError,
			}

			output, err := pager.GetPage(test.pageNumber)

			if test.errors && err == nil {
				t.Errorf("Expected error in test '%s' but received none", test.name)
			} else if !test.errors && err != nil{
				t.Errorf("Expected no error in test '%s' but got %v", test.name, err)
			}

			if !bytes.Equal(test.expected, output) {
				t.Errorf("Test '%s' failed, expected: \n%s\nbut got\n%s", test.name, test.expected, output)
			}
		})
	}
}

func TestPagerSet(t *testing.T) {
	tests := []struct {
		name         string
		pageNumber   int
		content      []byte
		expected     []byte
		fsError      error
		errors bool
	}{
		{
			name:       "write zeroth page",
			pageNumber: 0,
			content:    []byte("ab"),
			expected:   []byte("01ab456789"),
		},
		{
			name:       "write second page",
			pageNumber: 2,
			content:    []byte("cd"),
			expected:   []byte("012345cd89"),
		},

		// errors
		{
			name:         "write second page with fs error",
			pageNumber:   2,
			fsError:      errors.New("error"),
			content:      []byte("cd"),
			errors: true,
		},
		{
			name:         "write negative page",
			pageNumber:   -2,
			errors: true,
		},
		{
			name:         "write past freelist",
			pageNumber:   4,
			errors: true,
		},
		{
			name:         "write page at freelist start",
			pageNumber:   3,
			errors: true,
		},
		{
			name:         "write with content larger than page size",
			pageNumber:   0,
			errors: true,
			content:      []byte("abc"),
		},
		{
			name:         "write with content smaller than page size",
			pageNumber:   0,
			errors: true,
			content:      []byte("a"),
		},
	}

	pager := Pager{
		pageSize:      2,
		freelistIndex: 3,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set the fs on the pager
			m := &mockFile{
				content: []byte(numbers),
				err:     test.fsError,
			}
			pager.fs = m

			err := pager.SetPage(test.pageNumber, test.content)
			if test.errors {
				if err == nil {
				t.Errorf("Expected error in test '%s' but received none", test.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error in test '%s' but got %v", test.name, err)
				}
				if !bytes.Equal(test.expected, m.content) {
					t.Errorf("Test '%s' failed, expected: \n%s\nbut got\n%s", test.name, test.expected, m.content)
				}
			}
		})
	}
}
