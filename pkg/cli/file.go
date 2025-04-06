package cli

import (
	"fmt"
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
		if strings.Trim(query, " \n\t") == "" {
			continue
		}

		fmt.Println(query)
		result, err := db.Query(query, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(Render(result))
	}
}
