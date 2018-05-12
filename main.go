package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Running CI Tricks...")
	if err := RunTricks(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
