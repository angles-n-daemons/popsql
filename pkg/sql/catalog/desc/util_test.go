package desc_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

func TestGetDataType(t *testing.T) {
	for _, test := range []struct {
		tokenType scanner.TokenType
		datatype  desc.DataType
		error     string
	}{
		{scanner.DATATYPE_BOOLEAN, desc.BOOLEAN, ""},
		{scanner.DATATYPE_STRING, desc.STRING, ""},
		{scanner.DATATYPE_NUMBER, desc.NUMBER, ""},
		{scanner.BANG, 0, "unrecognized data type BANG"},
		{scanner.FROM, 0, "unrecognized data type FROM"},
	} {
		t.Run(fmt.Sprintf("tokenType=%s, dataType=%s", test.tokenType, test.datatype), func(t *testing.T) {
			datatype, err := desc.GetDataType(test.tokenType)
			if test.error == "" {
				assert.NoError(t, err)
			} else {
				assert.IsError(t, err, test.error)
			}
			assert.Equal(t, datatype, test.datatype)
		})
	}
}
