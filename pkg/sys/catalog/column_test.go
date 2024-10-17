package catalog_test

import (
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sys/catalog"
)

var schema *catalog.Schema
var table *catalog.Table

func init() {
	schema = catalog.InitSchema()
	table = &catalog.Table{
		Namespace: schema.User,
		Name:      "test",
	}
}

func TestColumnToRegister(t *testing.T) {
	column := &catalog.Column{
		Table:    table,
		Name:     "col1",
		DataType: catalog.STRING,
	}
	columnReg, err := (column).ToRegister()
	if err != nil {
		t.Fatal(err)
	}

	if columnReg.Table != column.Table.Namespace.Name+"-"+column.Table.Name {
		t.Fatalf(
			"columnReg.Table expected %s, but got %s",
			column.Table.Namespace.Name+"-"+column.Table.Name,
			columnReg.Table,
		)
	}
	if columnReg.Name != column.Name {
		t.Fatalf(
			"columnReg.Name expected %s, but got %s",
			column.Name,
			columnReg.Name,
		)
	}
	if columnReg.DataType != column.DataType {
		t.Fatalf(
			"columnReg.DataType expected %s, but got %s",
			column.DataType,
			columnReg.DataType,
		)
	}
}

func TestColumnRegisterKey(t *testing.T) {
	column := &catalog.Column{
		Table:    table,
		Name:     "col1",
		DataType: catalog.STRING,
	}
	columnReg, err := column.ToRegister()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := columnReg.Key()
	if err != nil {
		t.Fatal(err)
	}
	expected := column.Table.Namespace.Name + "-" + column.Table.Name + "-" + column.Name
	if actual != expected {
		t.Fatalf(
			"columnReg.Key() expected %s, but got %s",
			expected,
			actual,
		)
	}
}

func TestColumnRegisterValue(t *testing.T) {
	column := &catalog.Column{
		Table:    table,
		Name:     "col1",
		DataType: catalog.STRING,
	}
	columnReg, err := column.ToRegister()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := columnReg.Value()
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(`{"Table":"user-test","Name":"col1","DataType":1}`)
	if string(actual) != string(expected) {
		t.Fatalf(
			"columnReg.Key() expected %s, but got %s",
			expected,
			actual,
		)
	}
}

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
