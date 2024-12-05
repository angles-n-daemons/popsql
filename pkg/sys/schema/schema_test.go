package schema_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

func SchemaFromBytes(t *testing.T) {

}

func TestAddTable(t *testing.T) {
	sc := schema.NewSchema()
	sc.AddTable(testTableFromArgs("tt", nil, nil))
}

func TestAddExistingTable(t *testing.T) {

}

func TestGetTable(t *testing.T) {

}
func TestGetMissingTable(t *testing.T) {

}

func TestDropTable(t *testing.T) {

}
func TestDropMissingTable(t *testing.T) {

}
func TestDropRootTable(t *testing.T) {

}
