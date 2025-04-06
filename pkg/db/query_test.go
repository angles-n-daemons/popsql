package db_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db"
)

func TestQuery(t *testing.T) {
	db := db.GetEngine()
	db.Query("select * from users", []any{})
}
