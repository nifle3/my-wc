package main

import (
	"fmt"
	"os"
	"strconv"
)

type Arg struct {
	FileName string
	Mode     []string
}

func main() {
	arg, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.OpenFile(arg.FileName, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	result, err := ParseMode(arg.Mode, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func ParseArgs(args []string) (Arg, error) {
	if len(args) < 1 {
		return Arg{}, fmt.Errorf("usage: [mode] <file>")
	}

	result := Arg{
		FileName: args[len(args)-1],
		Mode:     args[:len(args)-1],
	}

	return result, nil
}

func ParseMode(args []string, file *os.File) (string, error) {
	funcs := map[string]func(file *os.File) (string, error){
		"-c": NumberOfBytes,
		"-l": NumberOfLines,
		"-w": NumberOfWords,
	}

	result := ""

	for _, mode := range args {
		f, ok := funcs[mode]
		if !ok {
			return "", fmt.Errorf("unknown mode: %s", mode)
		}

		semiResult, err := f(file)
		if err != nil {
			return result, err
		}

		result += semiResult
	}

	return result, nil
}

func NumberOfBytes(file *os.File) (string, error) {
	fileStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(fileStat.Size(), 10), nil
}

func NumberOfLines(file *os.File) (string, error) {
	result := make([]byte, 0)
	_, err := file.Read(result)
	if err != nil {
		return "", err
	}

	var numberOfLines int64 = 0
	for _, val := range result {
		if val == '\n' {
			numberOfLines++
		}
	}
	return strconv.FormatInt(numberOfLines, 10), nil
}

func NumberOfWords(file *os.File) (string, error) {
	return "", nil
}
