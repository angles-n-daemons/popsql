package schema_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

func TestGetDataType(t *testing.T) {
	for _, test := range []struct {
		tokenType scanner.TokenType
		datatype  schema.DataType
		error     string
	}{
		{scanner.DATATYPE_BOOLEAN, schema.BOOLEAN, ""},
		{scanner.DATATYPE_STRING, schema.STRING, ""},
		{scanner.DATATYPE_NUMBER, schema.NUMBER, ""},
		{scanner.BANG, 0, "unrecognized data type BANG"},
		{scanner.FROM, 0, "unrecognized data type FROM"},
	} {
		t.Run(fmt.Sprintf("tokenType=%s, dataType=%s", test.tokenType, test.datatype), func(t *testing.T) {
			datatype, err := schema.GetDataType(test.tokenType)
			if test.error != "" && err != nil {
				if test.error != err.Error() {
					t.Fatalf(
						"expected error %s but instead got %s",
						test.error,
						err,
					)
				}
				return
			}
			if err == nil && test.error != "" {
				t.Fatalf("expected error: %s", test.error)
			}
			if err != nil && test.error == "" {
				t.Fatalf("expected no error but received error: %s", err)
			}
			if datatype != test.datatype {
				t.Fatalf(
					"expected datatype %s, but got %s instead",
					test.datatype,
					datatype,
				)
			}
		})
	}
}
