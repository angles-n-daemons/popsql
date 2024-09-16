package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/engine"
)

var banner = `
                            __
   ___  ___  ___  ___ ___ _/ /
  / _ \/ _ \/ _ \(_-</ _ '/ / 
 / .__/\___/ .__/___/\_, /_/  
/_/       /_/         /_/     
`

var db = engine.Engine{DataDir: ""}

func REPL() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(banner)
	fmt.Println("version 0.0")
	query := ""
	for {
		if query == "" {
			fmt.Print("ppql> ")
		} else {
			fmt.Print("....  ")
		}
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading input: %s\n", err)
			os.Exit(1)
		}
		parts := strings.Split(text, ";")
		query += parts[0]

		// semicolon sent
		if len(parts) > 1 {
			err = db.Query(query)
			if err != nil {
				fmt.Println(err)
			}
			query = ""
		}
	}
}
