package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}

	switch command := os.Args[1]; command {
	case "help":
		fmt.Println("popsql <command>")
		os.Exit(0)
	default:
		fmt.Printf("unknown command %s", command)
		os.Exit(1)
	}
}
