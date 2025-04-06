package cli

import (
	"os"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/db"
)

func File(filename string) {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	db := db.GetEngine()
	queries := strings.Split(string(b), ";")

	for _, query := range queries {
		db.Query(query, nil)
	}
}
