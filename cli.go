package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		diff    bool
		write   bool
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprintf(cli.outStream, helpText, Name)
	}

	flags.BoolVar(&diff, "diff", false, "")
	flags.BoolVar(&diff, "d", false, "(Short)")

	flags.BoolVar(&write, "write", false, "")
	flags.BoolVar(&write, "w", false, "(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	files := flags.Args()
	if len(files) == 0 {
		fmt.Fprintf(cli.errStream, "TODO: os.Stdin?\n")
		return ExitCodeError
	}

	file := files[0]
	goFile, err := ParseFile(file)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to parse go file: %s\n", err)
		return ExitCodeError
	}
	Debugf("%#v", goFile)

	testFile, err := TestFile(file)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to get go test file: %s\n", err)
		return ExitCodeError
	}

	goTestFile, err := ParseFile(testFile)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to parse go test file: %s\n", err)
		return ExitCodeError
	}
	Debugf("%#v", goTestFile)

	opts := &DiffOpts{
		Mode: Strict,
	}

	diffFuncs, err := goFile.DiffFuncs(goTestFile, opts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to diff source file and test file: %s\n", err)
		return ExitCodeError
	}
	Debugf("%#v", diffFuncs)

	goTestFile.AddTestFuncs(diffFuncs)

	res, err := goTestFile.Generate()
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to generate source: %s\n", err)
		return ExitCodeError
	}

	if !bytes.Equal(goTestFile.Src, res) {

		if diff {
			data, err := doDiff(goTestFile.Src, res)
			if err != nil {
				fmt.Fprintf(cli.errStream, "Failed to compute diff: %s\n", err)
				return ExitCodeError
			}
			fmt.Fprintf(cli.outStream, "diff %s\n", goTestFile.FileName)
			fmt.Fprintf(cli.outStream, "%s\n", data)
		}

		if write {
			err = ioutil.WriteFile(testFile, res, 0644)
			if err != nil {
				fmt.Fprintf(cli.errStream, "Failed to write resutl to file: %s\n", err)
				return ExitCodeError
			}
		}
	}

	return ExitCodeOK
}

func doDiff(b1, b2 []byte) ([]byte, error) {
	f1, err := ioutil.TempFile("", Name)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", Name)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err := exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if err != nil {
		if len(data) > 0 {
			// diff exits with a non-zero status when the files don't match.
			// Ignore that failure as long as we get output.
			return data, nil
		}

		return nil, err
	}

	return data, err
}

var helpText = `Usage:

  %s [options] PATH ...

Options:

  -diff, -d    Display diffs instead of rewriting files.

  -write, -w   Write result to target file instead of stdout.
               For example, if source file name is 'main.go',
               target file will be 'main_test.go'.

`
