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
	all            bool
	envPrefix      string
	noStyles       bool
	fieldNames     bool
	debug          bool
}

func (cfg *appConfig) parseFlags(f *flag.FlagSet) error {
	f.StringVar(&cfg.outputFileName, "output", "", "Output file name")
	f.StringVar(&cfg.typeName, "type", "", "Type name")
	f.StringVar(&cfg.formatName, "format", "", "Output format, default `markdown`")
	f.BoolVar(&cfg.all, "all", false, "Generate documentation for all types in the file")
	f.StringVar(&cfg.envPrefix, "env-prefix", "", "Environment variable prefix")
	f.BoolVar(&cfg.noStyles, "no-styles", false, "Disable styles in html output")
	f.BoolVar(&cfg.fieldNames, "field-names", false, "Use field names if tag is not specified")
	f.BoolVar(&cfg.debug, "debug", false, "Enable debug mode")

	if err := f.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parsing CLI args: %w", err)
	}

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
	cfg, err := getConfig()
	if err != nil {
		fatal(err)
	}
	if cfg.debug {
		debugLogs = true
	}
	if err := run(&cfg); err != nil {
		fatal("Generate error:", err)
	}
}

type getConfigOpt func(f *flag.FlagSet)

var getConfigSilent = func(f *flag.FlagSet) {
	f.Usage = func() {}
	nopWriter := os.NewFile(0, os.DevNull)
	f.SetOutput(nopWriter)
}

func getConfig(opts ...getConfigOpt) (appConfig, error) {
	var cfg appConfig
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	for _, opt := range opts {
		opt(flagSet)
	}
	if err := cfg.parseFlags(flagSet); err != nil {
		flagSet.Usage()
		return cfg, fmt.Errorf("invalid CLI args: %w", err)
	}

	if err := cfg.parseEnv(); err != nil {
		return cfg, fmt.Errorf("invalid environment: %w", err)
	}
	return cfg, nil
}

func run(cfg *appConfig) (err error) {
	outputFile, err := os.Create(cfg.outputFileName)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer closeWith(outputFile, func(closeErr error) {
		if closeErr != nil {
			err = fmt.Errorf("closing output file: %w", err)
		}
	})

	var opts []generatorOption
	if cfg.all {
		opts = append(opts, withAll())
	} else if cfg.typeName != "" {
		opts = append(opts, withType(cfg.typeName))
	}
	if cfg.formatName != "" {
		opts = append(opts, withFormat(cfg.formatName))
	}
	if cfg.envPrefix != "" {
		opts = append(opts, withPrefix(cfg.envPrefix))
	}
	if cfg.noStyles {
		opts = append(opts, withNoStyles())
	}
	if cfg.fieldNames {
		opts = append(opts, withFieldNames())
	}
	gen, err := newGenerator(cfg.inputFileName, cfg.execLine, opts...)
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
