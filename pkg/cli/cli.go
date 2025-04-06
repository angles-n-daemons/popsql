package cli

import (
	"fmt"
	"os"
)

func Main(args []string) {
	if len(args) < 2 {
		REPL()
	} else if len(args) == 2 {
		if args[1] == "server" {
			Server()
		} else {
			File(args[1])
		}
	} else {
		fmt.Errorf("too many arguments...")
		os.Exit(1)
	}
}
