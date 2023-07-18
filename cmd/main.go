package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/angles-n-daemons/popsql/internal/backend/pager"
)

func main() {

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Test B-Tree")
		fmt.Println("2. Exit")

		choice, err := readInt()
		if err != nil {
			fmt.Println(err)
			return
		}

		switch choice {
		case 1:
			testBTree()
		case 2:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

func readInt() (int, error) {
	for {
		fmt.Println("Enter your choice: ")
		reader := bufio.NewReader(os.Stdin)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return -1, err
		}
		input = strings.TrimSpace(input)
		fmt.Println()

		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid integer.")
			continue
		}
		return value, nil
	}
}

func testBTree() {
	fmt.Println("Enter the page size: ")

	pageSize, err := readInt()
	if err != nil {
		fmt.Println(err)
		return
	}

	pager := pager.NewMemoryPager(uint16(pageSize))
	fmt.Println(pager)
}
