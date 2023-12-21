package main

import (
	"flag"
	"os"
	"path"
	"testing"
)

func TestConfig(t *testing.T) {
	var cfg appConfig
	os.Args = []string{"cmd", "-output", "test.md", "-type", "test", "-format", "markdown"}
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := cfg.parseFlags(flagSet); err != nil {
		t.Fatal("Invalid CLI args:", err)
	}

	if cfg.outputFileName != "test.md" {
		t.Fatal("Invalid output file name")
	}
	if cfg.typeName != "test" {
		t.Fatal("Invalid type name")
	}
	if cfg.formatName != "markdown" {
		t.Fatal("Invalid format name")
	}

	t.Setenv("GOFILE", "test.go")
	t.Setenv("GOLINE", "42")
	if err := cfg.parseEnv(); err != nil {
		t.Fatal("Invalid environment:", err)
	}
	if cfg.inputFileName != "test.go" {
		t.Fatalf("Invalid input file name: `%s`", cfg.inputFileName)
	}
	if cfg.execLine != 42 {
		t.Fatalf("Invalid line number: `%d`", cfg.execLine)
	}
}

func TestInvalidConfig(t *testing.T) {
	t.Run("outputFileName", func(t *testing.T) {
		var cfg appConfig
		os.Args = []string{"cmd", "-type", "test", "-format", "markdown"}
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		err := cfg.parseFlags(flagSet)
		if err == nil {
			t.Fatal("Invalid CLI args:", err)
		}
		t.Logf("Got error as expected: %v", err)
	})

	t.Run("inputFileName", func(t *testing.T) {
		var cfg appConfig
		err := cfg.parseEnv()
		t.Setenv("GOFILE", "")
		if err == nil {
			t.Fatal("Invalid environment:", err)
		}
		t.Logf("Got error as expected: %v", err)
	})
	t.Run("noExecLine", func(t *testing.T) {
		var cfg appConfig
		t.Setenv("GOFILE", "test.go")
		err := cfg.parseEnv()
		if err == nil {
			t.Fatal("Invalid environment:", err)
		}
		t.Logf("Got error as expected: %v", err)
	})
	t.Run("execLine", func(t *testing.T) {
		var cfg appConfig
		t.Setenv("GOFILE", "test.go")
		t.Setenv("GOLINE", "a")
		err := cfg.parseEnv()
		if err == nil {
			t.Fatal("Invalid environment:", err)
		}
		t.Logf("Got error as expected: %v", err)
	})
}

func TestMainRun(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		inputFile := path.Join(t.TempDir(), "example.go")
		if err := copyTestFile(path.Join("testdata", "example_type.go"), inputFile); err != nil {
			t.Fatal("copy test file", err)
		}
		outputFile := path.Join(t.TempDir(), "example.md")
		config := appConfig{
			typeName:       "Type1",
			formatName:     "markdown",
			outputFileName: outputFile,
			inputFileName:  inputFile,
			execLine:       0,
		}
		if err := run(&config); err != nil {
			t.Fatal("run", err)
		}
	})
}
