package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:generate ./bin/gotests -godoc

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// defaultExcludes are default directory where walkFunc does not walk
var defaultExcludes = []string{".git", "Godep", "vendor"}

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		diff              bool
		write             bool
		list              bool
		includeUnexported bool
		reverse           bool
		version           bool

		doc bool
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

	flags.BoolVar(&list, "list", false, "")
	flags.BoolVar(&list, "l", false, "(Short)")

	flags.BoolVar(&reverse, "reverse", false, "")
	flags.BoolVar(&reverse, "r", false, "(Short)")

	flags.BoolVar(&includeUnexported, "include-unexported", false, "")
	flags.BoolVar(&includeUnexported, "i", false, "")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit.")

	// This flag is only for developer to generate godoc via go generate.
	flags.BoolVar(&doc, "godoc", false, "")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Generate godoc (only for developer)
	if doc {
		if err := godoc("doc.go"); err != nil {
			fmt.Fprintf(cli.errStream, "Failed to generate godoc: %s", err)
			return ExitCodeError
		}
		return ExitCodeOK
	}

	paths := flags.Args()
	if len(paths) == 0 {
		fmt.Fprintf(cli.errStream, "Invalid arguments. You must provide PATHs\n")
		return ExitCodeError
	}

	// opts are option struct for processGenerate()
	opts := &generateOpts{
		diffOpts: &diffOpts{
			Mode:              Strict,
			IncludeUnexported: includeUnexported,
		},
		diff:    diff,
		write:   write,
		list:    list,
		reverse: reverse,
	}

	// By default, statusCode is ExitCodeOK and Run() returns it.
	// It is updated only when processGogenerate returns non-ExitCodeOK.
	exitCode := ExitCodeOK

	for _, path := range paths {
		switch fi, err := os.Stat(path); {
		case err != nil:
			// Output the error and proceeds next (but Change status code).
			fmt.Fprintf(cli.errStream, "Failed to get file info: %s", err)
			exitCode = ExitCodeError

		case fi.IsDir():
			// walkFn is function for filepath.Walk. It walks through .go files
			// (but non _test.go) and executes processGenerate() to each file.
			// If error happens while processing, it display it to errStream
			// and continues processing. This is same as gofmt does.
			//
			// It updates exitCode only when processGenerates returns not non-zero code.
			walkFn := func(srcPath string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Ignore if it's directory
				if fi.IsDir() {
					return nil
				}

				// Ignore .git and vendoring directory.
				for _, exclude := range defaultExcludes {
					if strings.Contains(srcPath, exclude) {
						return nil
					}
				}

				// Ignore non .go file and _test.go file.
				if !isGoFile(fi) {
					return nil
				}

				Debugf("Walk to %q", srcPath)
				status := cli.processGenerate(srcPath, opts)
				if status != ExitCodeOK {
					exitCode = status
				}

				return nil
			}

			// Start walking.
			if err := filepath.Walk(path, walkFn); err != nil {
				fmt.Fprintf(cli.errStream, "Failed to walk: %s\n", err)
				exitCode = ExitCodeError
			}
		default:
			status := cli.processGenerate(path, opts)
			if status != ExitCodeOK {
				exitCode = status
			}
		}
	}

	return exitCode
}

type generateOpts struct {
	diffOpts *diffOpts

	diff  bool
	write bool
	list  bool

	reverse bool
}

func (cli *CLI) processGenerate(srcPath string, opts *generateOpts) int {
	var testPath string
	if opts.reverse {
		var err error
		testPath = srcPath
		srcPath, err = SrcFilePath(testPath)
		if err != nil {
			fmt.Errorf("Failed to get src file path: %s\n", err)
			return ExitCodeError
		}
	} else {
		var err error
		testPath, err = TestFilePath(srcPath)
		if err != nil {
			fmt.Errorf("Failed to get go test file path: %s\n", err)
			return ExitCodeError
		}
	}

	// Run actual gotests to path
	goTestFile, err := goTestGenerate(srcPath, testPath, opts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to generate: %s\n", err)
		return ExitCodeError
	}

	// Genreate results as a []byte
	resBytes, err := goTestFile.Generate()
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to generate result from ast: %s\n", err)
		return ExitCodeError
	}

	// Handle diff/write only when there is diff between result and original code.
	if !bytes.Equal(goTestFile.SrcBytes, resBytes) {

		if opts.list {

			path, err := fmtPath(goTestFile.FileName)
			if err != nil {
				fmt.Fprintf(cli.errStream, "Failed to format path: %s\n", err)
				return ExitCodeError
			}
			fmt.Fprintf(cli.outStream, "%s\n", path)
		}

		if opts.diff {
			data, err := doDiff(goTestFile.SrcBytes, resBytes)
			if err != nil {
				fmt.Fprintf(cli.errStream, "Failed to compute diff: %s\n", err)
				return ExitCodeError
			}
			fmt.Fprintf(cli.outStream, "diff %s\n", goTestFile.FileName)
			fmt.Fprintf(cli.outStream, "%s\n", data)
		}

		if opts.write {
			err = ioutil.WriteFile(testPath, resBytes, 0644)
			if err != nil {
				fmt.Fprintf(cli.errStream, "Failed to write resutl to file: %s\n", err)
				return ExitCodeError
			}
		}
	}

	if !opts.list && !opts.diff && !opts.write {
		_, err := cli.outStream.Write(resBytes)
		if err != nil {
			fmt.Fprintf(cli.errStream, "Failed to write resutl: %s\n", err)
			return ExitCodeError
		}
	}

	return ExitCodeOK
}

func goTestGenerate(srcPath, testPath string, opts *generateOpts) (*GoFile, error) {
	goFile, err := ParseFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go file: %s", err)
	}
	Debugf("%#v", goFile)

	var goTestFile *GoFile
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		// If test file is not exist, create new one with the same pacakge
		// declare with the source.
		var err error
		goTestFile, err = NewGoFile(testPath, goFile.PackageName)
		if err != nil {
			return nil, fmt.Errorf("failed to create new test file: %s", err)
		}
	} else {
		// If test file is exist, just parse it.
		var err error
		goTestFile, err = ParseFile(testPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse go test file: %n", err)
		}
	}
	Debugf("goTestFile: %#v", goTestFile)

	diffFuncs, err := goFile.diffFuncs(goTestFile, opts.diffOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to diff source file and test file: %s", err)
	}
	Debugf("Diff Funcs: %#v", diffFuncs)

	funcTmpl := defaultExpectTestFuncTmpl
	if err := goTestFile.addFuncTestFuncs(diffFuncs, funcTmpl); err != nil {
		return nil, fmt.Errorf("failed to add func test funcs: %s", err)
	}

	diffMethods, err := goFile.diffMethods(goTestFile, opts.diffOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to diff source file and test file: %s", err)
	}
	Debugf("Diff Methods: %#v", diffMethods)

	funcTmpl = defaultExpectTestFuncMethodTmpl
	if err := goTestFile.addMethodTestFuncs(diffMethods, funcTmpl); err != nil {
		return nil, fmt.Errorf("failed to add method test funcs: %s", err)
	}

	return goTestFile, nil
}

func fmtPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Rel(currentPath, path)
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

// isGoFile returns true if file is go file and it's not test file.
func isGoFile(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	name := fi.Name()
	if strings.HasPrefix(name, ".") {
		return false
	}

	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

// godoc generates doc.go file for godoc to prevent from writing
// the same documentation twice. If any, return error.
func godoc(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	tmpl, err := template.New("godoc").Parse(godocTmpl)
	if err != nil {
		return err
	}

	data := struct {
		Content string
	}{
		Content: helpText,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	return nil
}

var godocTmpl = `// DON"T EDIT THIS FILE
// THIS IS GENERATED VIA GO GENERATE

/*
{{ .Content }}
*/
package main
`

var helpText = `gotests is tool to generate Go test functions from
the given source code like gofmt. 

https://github.com/tcnksm/gotests

Usage:

  gotests [options] PATH ...

Options:

  -diff, -d    Display diffs instead of rewriting files.

  -write, -w   Write result to target file instead of stdout.
               For example, if source file name is 'main.go',
               target file will be 'main_test.go'.

  -list, -l    List test files to be updated/generated.

  -i           Include unexport function/method for generating target.
`
