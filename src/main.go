package main

import (
	"bufio"
	"fmt"
	"os"

	"net/http"
	_ "net/http/pprof"

	"in-mem-db/src/db"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	d := db.Init()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Command: ")
		scanner.Scan()
		input := scanner.Text()
		d.Handle(input)
	}
}
