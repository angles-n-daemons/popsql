package main

import (
	"os"

	"github.com/angles-n-daemons/popsql/pkg/cli"
)

func main() {
	cli.REPL(os.Args)
}
