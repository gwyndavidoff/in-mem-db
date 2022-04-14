package main

import (
	"bufio"
	"fmt"
	"os"

	"in-mem-db/db"
)

func main() {
	d := db.Init()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Command: ")
		scanner.Scan()
		input := scanner.Text()
		d.Handle(input)
	}
}
