package cli

import (
	"fmt"
	"os"

	"github.com/angles-n-daemons/popsql/pkg/server"
)

var banner = `
      ┏┓┏┓┓ 
┏┓┏┓┏┓┗┓┃┃┃ 
┣┛┗┛┣┛┗┛┗┻┗┛
┛   ┛       
`

func Main(args []string) {
	fmt.Println(banner)
	fmt.Println("version 0.0")
	if len(args) < 2 {
		REPL()
	} else if len(args) == 2 {
		if args[1] == "server" {
			server.Run()
		} else {
			File(args[1])
		}
	} else {
		fmt.Println(fmt.Errorf("too many arguments..."))
		os.Exit(1)
	}
}
