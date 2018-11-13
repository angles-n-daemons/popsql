package pager

import "testing"

type mockFile struct {
	index   int
	content []bytes
	err     error
}

func (m mockFile) Read(int) ([]byte, error) {
	return m.content, m.err
}
func (m mockFile) Write(int, []content) error {
	return m.err

}
func (m mockFile) Size() (int, error) {
	return m.size, m.err
}

const (
	testPageSize = 2
	alphabet     = "abcdefghijklmnopqrstuvwxyz"
)

func TestPagerGet(t *testing.T) {
	tests := []struct {
		pageNumber int
		expected   []byte
		fsErr      error
		errors     bool
	}{
		{0, []byte("cd"), nil, false},
	}

	pager := Pager{
		pageSize: 2,
	}

	for _, test := tests {
		// set the fs on the pager
		pager.fs = mockFile{
			content: []byte(alphabet),
			err: test.fsErr,
		}

		output, err := pager.Get(test.pageNumber)
	}
}
