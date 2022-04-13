package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gwyndavidoff/in-mem-db/db"
)

// type DB struct {
// 	database map[string]string
// 	counts   map[string]int
// }

func main() {
	db := &db.DB{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Command: ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Split(input, " ")
		if strings.EqualFold(args[0], "END") {
			break
		}
		for _, r := range args {
			fmt.Println(r)
		}
		fmt.Println(input)
	}
}
