package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

}
