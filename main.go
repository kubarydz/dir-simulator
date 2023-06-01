package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	inputFilename := flag.String("input", "input.txt", "input file")
	outputFilename := flag.String("output", "output.txt", "output file")
	flag.Parse()

	processCommands(*inputFilename, *outputFilename)
}

func processCommands(inputFilename, outputFilename string) {
	fs := CreateFilesystem()

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		fmt.Printf("cannot open input file %v, error: %v\n", inputFilename, err)
		panic("cannot open file")
	}
	defer inputFile.Close()
	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)

	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("cannot open output file %v, error: %v\n", outputFilename, err)
		panic("cannot open file")
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	for fileScanner.Scan() {
		cmd := fileScanner.Text()
		writer.WriteString(getCommandEcho(cmd) + "\n")
		for _, output := range handleCommand(cmd, fs) {
			writer.WriteString(output + "\n")
		}
	}
	writer.Flush()
}

func readInput(filename string) []string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("cannot open %v, error: %v\n", filename, err)
		panic("cannot open file")
	}
	strs := []string{}
	lastStart := 0
	for i, b := range bytes {
		if b == '\n' {
			strs = append(strs, string(bytes[lastStart:i]))
			lastStart = i + 1
		}
	}

	return strs
}
