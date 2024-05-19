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
	countChars = "-m"
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

	if arg.FileName != "" {
		reader, err := os.OpenFile(arg.FileName, os.O_RDONLY, os.FileMode(0644))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Stdin = reader
	}

	result, err := ParseMode(arg.Mode)
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
	fileName := args[len(args)-1]
	if len(args) == 1 && !strings.HasPrefix(args[0], "-") {
		mods = []string{countWords, countBytes, countLines}
	} else if strings.HasPrefix(args[len(args)-1], "-") {
		fileName = ""
		mods = args[:len(args)-1]
	} else {
		mods = args[:len(args)-1]
	}

	result := Arg{
		FileName: fileName,
		Mode:     mods,
	}

	return result, nil
}

func ParseMode(args []string) (string, error) {
	funcs := map[string]func() (string, error){
		countBytes: NumberOfBytes,
		countLines: NumberOfLines,
		countWords: NumberOfWords,
		countChars: NumberOfChars,
	}

	result := ""

	for _, mode := range args {
		f, ok := funcs[mode]
		if !ok {
			return "", fmt.Errorf("unknown mode: %s", mode)
		}

		semiResult, err := f()
		if err != nil {
			return "", err
		}

		result += semiResult + " "
	}

	return result, nil
}

func NumberOfBytes() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(stat.Size(), 10), nil
}

func NumberOfLines() (string, error) {
	bufferReader := bufio.NewReader(os.Stdin)
	numberOfLines := 0

	for {
		_, err := bufferReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		numberOfLines += 1
	}

	return strconv.FormatInt(int64(numberOfLines), 10), nil
}

func NumberOfWords() (string, error) {
	bufferReader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, 1024)
	numberOfWords := 0

	for {
		c, err := bufferReader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		wordInFile := bytes.FieldsFunc(buffer[:c], func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})

		numberOfWords += len(wordInFile)
	}

	return strconv.FormatInt(int64(numberOfWords), 10), nil
}

func NumberOfChars() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)
	result := 0
	for scanner.Scan() {
		result++
	}

	return strconv.FormatInt(int64(result), 10), nil
}
