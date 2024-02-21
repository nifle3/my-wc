package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	countWords = "-w"
	countLines = "-l"
	countBytes = "-c"
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

	result += arg.FileName

	fmt.Println(result)
}

func ParseArgs(args []string) (Arg, error) {
	if len(args) < 1 {
		return Arg{}, fmt.Errorf("usage: [mode] <file>")
	}

	var mods []string

	if len(args) == 1 && !strings.HasPrefix(args[0], "-") {
		mods = []string{countWords, countBytes, countLines}
	} else {
		mods = args[:len(args)-1]
	}

	result := Arg{
		FileName: args[len(args)-1],
		Mode:     mods,
	}

	return result, nil
}

func ParseMode(args []string, file *os.File) (string, error) {
	funcs := map[string]func(file *os.File) (string, error){
		countBytes: NumberOfBytes,
		countLines: NumberOfLines,
		countWords: NumberOfWords,
	}

	result := ""

	for _, mode := range args {
		f, ok := funcs[mode]
		if !ok {
			return "", fmt.Errorf("unknown mode: %s", mode)
		}

		semiResult, err := f(file)
		if err != nil {
			return "", err
		}

		result += semiResult + " "
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
	reader := bufio.NewReader(file)
	numberOfLines := 0

	buffer := make([]byte, 1024)

	for {
		c, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		numberOfLines += bytes.Count(buffer[:c], []byte{'\n'})
	}

	return strconv.FormatInt(int64(numberOfLines+1), 10), nil
}

func NumberOfWords(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	numberOfWords := 0

	for {
		c, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		wordInFile := strings.FieldsFunc(string(buffer[:c]), func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})

		numberOfWords += len(wordInFile)
	}

	return strconv.FormatInt(int64(numberOfWords), 10), nil
}
