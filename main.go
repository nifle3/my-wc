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

	if len(os.Args) == 2 {
	} else if len(os.Args) == 3 {
		if os.Args[1] != "-c" {
			fmt.Printf("Usage: %s [-c] <file>\n", os.Args[0])
			os.Exit(1)
		}

		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()
		statFile, err := file.Stat()
		fmt.Println(statFile.Size())

	} else {
		fmt.Printf("Usage: %s [-command] <file>\n", os.Args[0])
		os.Exit(1)
	}
}
