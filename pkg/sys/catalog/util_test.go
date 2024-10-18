package catalog_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/catalog"
)

func TestGetDataType(t *testing.T) {
	for _, test := range []struct {
		tokenType scanner.TokenType
		datatype  catalog.DataType
		error     string
	}{
		{scanner.DATATYPE_BOOLEAN, catalog.BOOLEAN, ""},
		{scanner.DATATYPE_STRING, catalog.STRING, ""},
		{scanner.DATATYPE_NUMBER, catalog.NUMBER, ""},
		{scanner.BANG, 0, "unrecognized data type BANG"},
		{scanner.FROM, 0, "unrecognized data type FROM"},
	} {
		t.Run(fmt.Sprintf("tokenType=%s, dataType=%s", test.tokenType, test.datatype), func(t *testing.T) {
			datatype, err := catalog.GetDataType(test.tokenType)
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
