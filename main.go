package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var (
		outputFileName string
		typeName       string
		formatName     string
	)
	flag.StringVar(&outputFileName, "output", "", "Output file name")
	flag.StringVar(&typeName, "type", "", "Type name")
	flag.StringVar(&formatName, "format", "", "Output format, default `markdown`")
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
	defer closeWith(outputFile, func(err error) {
		fatal("close output file", err)
	})

	gen, err := newGenerator(inputFileName, execLine,
		withType(typeName), withFormat(formatName))
	if err != nil {
		fatal("Error creating generator:", err)
	}
	if err := gen.generate(outputFile); err != nil {
		fatal("Error generating documentation:", err)
	}
}

func fatal(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}
