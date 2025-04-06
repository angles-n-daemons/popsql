package db_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/cli"
	"github.com/angles-n-daemons/popsql/pkg/db"
)

func TestQuery(t *testing.T) {
	db := db.GetEngine()
	result, err := db.Query("select * from __tables__", []any{})
	if err != nil {
		t.Fatal(err)
	}
	cli.Render(result)
}
