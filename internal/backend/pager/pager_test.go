package pager

import (
	"errors"
	"testing"

	"github.com/matryer/is"
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
	is := is.New(t)

	tests := []struct {
		name         string
		pageNumber   int
		expected     []byte
		fsError      error
		returnsError bool
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
			returnsError: true,
		},
		{
			name:         "read negative page",
			pageNumber:   -2,
			returnsError: true,
		},
		{
			name:         "read past freelist",
			pageNumber:   4,
			returnsError: true,
		},
		{
			name:         "read page at freelist start",
			pageNumber:   3,
			returnsError: true,
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

			output, err := pager.Get(test.pageNumber)
			if test.returnsError {
				is.True(err != nil)
			} else {
				is.NoErr(err)
			}

			is.Equal(
				string(test.expected),
				string(output),
			)
		})
	}
}

func TestPagerSet(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name         string
		pageNumber   int
		content      []byte
		expected     []byte
		fsError      error
		returnsError bool
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
			returnsError: true,
		},
		{
			name:         "write negative page",
			pageNumber:   -2,
			returnsError: true,
		},
		{
			name:         "write past freelist",
			pageNumber:   4,
			returnsError: true,
		},
		{
			name:         "write page at freelist start",
			pageNumber:   3,
			returnsError: true,
		},
		{
			name:         "write with content larger than page size",
			pageNumber:   0,
			returnsError: true,
			content:      []byte("abc"),
		},
		{
			name:         "write with content smaller than page size",
			pageNumber:   0,
			returnsError: true,
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

			err := pager.Set(test.pageNumber, test.content)
			if test.returnsError {
				is.True(err != nil)
			} else {
				is.NoErr(err)
				is.Equal(
					string(test.expected),
					string(m.content),
				)
			}

		})
	}
}
