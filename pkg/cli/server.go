package cli

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/api/wire"
)

func Server() {
	srv := wire.NewServer()
	err := srv.ListenAndServe(":5432")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
