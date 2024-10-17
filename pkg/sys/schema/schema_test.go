package schema_test

import "testing"

// The goals of this package are:
// typedefs for schema, table, column
// serialization, deserialization of schema
// key values for table
// information about system tables
// all saveable records which include tables, columns require:
// - a key fun
// - a value func
// - a FromBytes
// - a ToBytes
// the latter two aren't really required to be written
// as they will be just pass through for json

func TestGetTable(t *testing.T) {

}

func TestTableFunctionality(t *testing.T) {

}

func TestColumnFunctionality(t *testing.T) {

}

func TestShouldIhavepointers(t *testing.T) {

}

func TestCantOverrideExistingTable(t *testing.T) {

}

func TestLoadSchema(t *testing.T) {

}
