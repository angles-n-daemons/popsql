package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/engine"
)

var banner = `
      ┏┓┏┓┓ 
┏┓┏┓┏┓┗┓┃┃┃ 
┣┛┗┛┣┛┗┛┗┻┗┛
┛   ┛       
`

func REPL(args []string) {
	fmt.Println(banner)
	fmt.Println("version 0.0")
	db := engine.NewEngine(engine.Options{})
	loop(db)
}

func loop(db *engine.Engine) {
	reader := bufio.NewReader(os.Stdin)
	query := ""
	for {
		if query == "" {
			fmt.Print("ppql> ")
		} else {
			fmt.Print("....  ")
		}
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %s\n", err)
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
