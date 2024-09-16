package main

import (
	"fmt"
	"os"

	"github.com/angles-n-daemons/popsql/pkg/cli"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "repl")
	}

	switch command := os.Args[1]; command {
	case "repl":
		cli.REPL()
	case "help":
		fmt.Println("popsql <command>")
		os.Exit(0)
	default:
		fmt.Printf("unknown command %s", command)
		os.Exit(1)
	}
}
