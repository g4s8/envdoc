package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type appConfig struct {
	typeName       string
	formatName     string
	outputFileName string
	inputFileName  string
	execLine       int
}

func (cfg *appConfig) parseFlags(f *flag.FlagSet) error {
	f.StringVar(&cfg.outputFileName, "output", "", "Output file name")
	f.StringVar(&cfg.typeName, "type", "", "Type name")
	f.StringVar(&cfg.formatName, "format", "", "Output format, default `markdown`")
	f.Parse(os.Args[1:])

	if cfg.outputFileName == "" {
		return fmt.Errorf("output file name is required")
	}
	return nil
}

func (cfg *appConfig) parseEnv() error {
	inputFileName := os.Getenv("GOFILE")
	if inputFileName == "" {
		return fmt.Errorf("No input file specified, this tool should be called by go generate")
	}
	cfg.inputFileName = inputFileName

	if e := os.Getenv("GOLINE"); e != "" {
		i, err := strconv.Atoi(e)
		if err != nil {
			return fmt.Errorf("Invalid line number specified, this tool should be called by go generate")
		}
		cfg.execLine = i
	} else {
		return fmt.Errorf("No line number specified, this tool should be called by go generate")
	}
	return nil
}

func main() {
	var cfg appConfig
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := cfg.parseFlags(flagSet); err != nil {
		flagSet.Usage()
		fatal("Invalid CLI args:", err)
	}

	if err := cfg.parseEnv(); err != nil {
		fatal("Invalid environment:", err)
	}

	if err := run(&cfg); err != nil {
		fatal("Generate error:", err)
	}
}

func run(cfg *appConfig) (err error) {
	outputFile, err := os.Create(cfg.outputFileName)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer closeWith(outputFile, func(err error) {
		if err != nil {
			err = fmt.Errorf("closing output file: %w", err)
		}
	})

	gen, err := newGenerator(cfg.inputFileName, cfg.execLine,
		withType(cfg.typeName), withFormat(cfg.formatName))
	if err != nil {
		return fmt.Errorf("creating generator: %w", err)
	}
	if err := gen.generate(outputFile); err != nil {
		return fmt.Errorf("generating documentation: %w", err)
	}
	return nil
}

func fatal(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}
