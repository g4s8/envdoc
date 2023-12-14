package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var outputFileName string
	var typeName string
	flag.StringVar(&outputFileName, "output", "", "Output file name")
	flag.StringVar(&typeName, "type", "", "Type name")
	flag.Parse()

	if outputFileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	inputFileName := os.Getenv("GOFILE")
	if inputFileName == "" {
		fatal("No input file specified, this tool should be called by go generate")
	}

	var execLine int
	if e := os.Getenv("GOLINE"); e != "" {
		i, err := strconv.Atoi(e)
		if err != nil {
			fatal("Invalid line number specified, this tool should be called by go generate")
		}
		execLine = i
	} else {
		fatal("No line number specified, this tool should be called by go generate")
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			fatalf("close output file: %v", err)
		}
	}()

	output := newMarkdownOutput(outputFile)
	output.writeHeader()
	defer func() {
		if err := output.Close(); err != nil {
			fatalf("close output: %v", err)
		}
	}()

	insp := newInspector(typeName, output, execLine)
	if err := insp.inspectFile(inputFileName); err != nil {
		fatalf("inspect file: %v", err)
	}

	fmt.Printf("Documentation generated and saved to %s\n", outputFileName)
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
