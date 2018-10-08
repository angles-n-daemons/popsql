package processor_test

import (
	"errors"
	"fmt"
	"testing"

	core "github.com/angles-n-daemons/popsql/internal/core"
	"github.com/angles-n-daemons/popsql/internal/core/vm"
)

func TestQuery(t *testing.T) {
	tests := []*testCompiler{
		{query: "1", errors: false},
		{query: "2", errors: true},
	}

	for _, tc := range tests {
		t.Run(
			fmt.Sprintf("query=%s errors=%t", tc.query, tc.errors),
			func(t *testing.T) {
				p := core.NewProcessor(tc)
				_, err := p.Query(tc.query)

				if (err != nil) != tc.errors {
					t.Fatalf("expected to error %t got %v", tc.errors, err)
				}

				if tc.calledQuery != tc.query {
					t.Fatalf("expected to be called with %s got %s", tc.query, tc.calledQuery)
				}
			},
		)
	}
}

type testCompiler struct {
	calledQuery string
	query       string
	errors      bool
}

func (t *testCompiler) Compile(s string) ([]vm.Instruction, error) {
	t.calledQuery = s
	if t.errors {
		return nil, errors.New("testCompiler.Compile: forced error")
	}
	return nil, nil
}
