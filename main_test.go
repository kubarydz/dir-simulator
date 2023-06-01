package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestProcessCommands(t *testing.T) {
	tests := []struct {
		name                   string
		inputFilename          string
		outputFilename         string
		expectedOutputFilename string
	}{
		{
			name:                   "test dir, mkdir, up and cd",
			inputFilename:          "resources/test_input1.txt",
			outputFilename:         "resources/output1.txt",
			expectedOutputFilename: "resources/test_output1.txt",
		},
		{
			name:                   "test tree and mv",
			inputFilename:          "resources/test_input2.txt",
			outputFilename:         "resources/output2.txt",
			expectedOutputFilename: "resources/test_output2.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Remove(tt.outputFilename)
			})
			processCommands(tt.inputFilename, tt.outputFilename)

			if !deepCompare(tt.expectedOutputFilename, tt.outputFilename) {
				t.Fatal("output doesn't match expected file")
			}
		})
	}
}

const chunckSize = 64000

// compare files line by line
// failfast if the size differs
func deepCompare(file1, file2 string) bool {
	f1s, err := os.Stat(file1)
	if err != nil {
		log.Fatal(err)
	}
	f2s, err := os.Stat(file2)
	if err != nil {
		log.Fatal(err)
	}

	if f1s.Size() != f2s.Size() {
		return false
	}

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, chunckSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunckSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF && err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
