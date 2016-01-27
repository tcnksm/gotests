package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

var (
	defaultExpectTestFuncTmpl       = "Test{{ title .Name }}"
	defaultExpectTestFuncMethodTmpl = "Test{{ title .RecvName }}_{{ title .Name }}"
)

// defaultIgnoreFuncs is default function name to be ignored in Parse
var defaultIgnoreFuncs = []string{"init"}

var funcMap = template.FuncMap{
	"title": strings.Title,
}

// Mode is diff mode to change function diff behavior
type Mode int

const (
	Strict Mode = iota
)

type diffOpts struct {
	Mode Mode

	IgnoreFuncs []string

	// IncludeUnexported include unexported function
	// for test generating target. By default it's false.
	IncludeUnexported bool

	ExpectTestFuncTmpl       string
	ExpectTestFuncMethodTmpl string
}

// GoFile is .go source file
type GoFile struct {
	PackageName string
	FileName    string
	SrcBytes    []byte
	Funcs       []string
	Methods     []*Method

	FSet    *token.FileSet
	AstFile *ast.File
}

type Method struct {
	RecvName string
	Name     string
}

func NewGoFile(filename, pkgName string) (*GoFile, error) {
	code := fmt.Sprintf("package %s", pkgName)
	rd := bytes.NewReader([]byte(code))

	goFile, err := parse(filename, rd)
	if err != nil {
		return nil, err
	}

	// Src should be emtpy (it's used for diff)
	goFile.SrcBytes = []byte{}

	return goFile, nil
}

func (o *diffOpts) init() {
	if o.ExpectTestFuncTmpl == "" {
		o.ExpectTestFuncTmpl = defaultExpectTestFuncTmpl
	}

	if o.ExpectTestFuncMethodTmpl == "" {
		o.ExpectTestFuncMethodTmpl = defaultExpectTestFuncMethodTmpl
	}

	if len(o.IgnoreFuncs) == 0 {
		o.IgnoreFuncs = defaultIgnoreFuncs
	}
}

func (gf *GoFile) Generate() ([]byte, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, gf.FSet, gf.AstFile); err != nil {
		return nil, err
	}

	return imports.Process(gf.FileName, buf.Bytes(), nil)
}

func (gf *GoFile) addFuncTestFuncs(funcs []string, funcTmpl string) error {
	for _, fun := range funcs {

		tmpl, err := template.New("testFunc").Funcs(funcMap).Parse(funcTmpl)
		if err != nil {
			return err
		}

		tmplData := struct {
			Name string
		}{
			Name: fun,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, tmplData); err != nil {
			return err
		}

		testFunDecl := NewTestFuncDecl(buf.String())
		gf.AstFile.Decls = append(gf.AstFile.Decls, testFunDecl)
	}

	return nil
}

func (gf *GoFile) addMethodTestFuncs(methods []*Method, funcTmpl string) error {
	for _, method := range methods {

		tmpl, err := template.New("testFunc").Funcs(funcMap).Parse(funcTmpl)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, method); err != nil {
			return err
		}

		testFunDecl := NewTestFuncDecl(buf.String())
		gf.AstFile.Decls = append(gf.AstFile.Decls, testFunDecl)
	}

	return nil
}

func (goFile *GoFile) diffFuncs(goTestFile *GoFile, opts *diffOpts) ([]string, error) {
	opts.init()

	var diff []string
	for _, fun := range goFile.Funcs {

		if contains(opts.IgnoreFuncs, fun) {
			continue
		}

		if !opts.IncludeUnexported && isUnExported(fun) {
			continue
		}

		// exist indicate expected test function is exist on goTestFile
		// function list.
		exist := false

		tmpl, err := template.New("testFunc").Funcs(funcMap).Parse(opts.ExpectTestFuncTmpl)
		if err != nil {
			return diff, err
		}

		tmplData := struct {
			Name string
		}{
			Name: fun,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, tmplData); err != nil {
			return diff, err
		}

		Debugf("Expect TestFunc Name: %s", buf.String())
		expectTestFun := buf.String()

		for _, testFun := range goTestFile.Funcs {
			switch mode := opts.Mode; mode {
			case Strict:
				if expectTestFun == testFun {
					exist = true
				}
			default:
				// Should not reach here...
				return diff, fmt.Errorf("unknown diff mode is provided: %d", mode)
			}
		}

		if !exist {
			diff = append(diff, fun)
		}
	}

	return diff, nil
}

func (goFile *GoFile) diffMethods(goTestFile *GoFile, opts *diffOpts) ([]*Method, error) {
	opts.init()

	var diff []*Method
	for _, method := range goFile.Methods {

		if !opts.IncludeUnexported && isUnExported(method.Name) {
			continue
		}

		// exist indicate expected test function is exist on goTestFile
		// function list.
		exist := false

		tmpl, err := template.New("testFunc").Funcs(funcMap).Parse(opts.ExpectTestFuncMethodTmpl)
		if err != nil {
			return diff, err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, method); err != nil {
			return diff, err
		}

		Debugf("Expect TestFunc Name: %s", buf.String())
		expectTestFun := buf.String()

		for _, testFun := range goTestFile.Funcs {
			switch mode := opts.Mode; mode {
			case Strict:
				if expectTestFun == testFun {
					exist = true
				}
			default:
				// Should not reach here...
				return diff, fmt.Errorf("unknown diff mode is provided: %d", mode)
			}
		}

		if !exist {
			diff = append(diff, method)
		}
	}

	return diff, nil
}

var reLower = regexp.MustCompile("^[a-z]+")

// isUnexported checks the given function is unxported (
// Check name start with lower case).
func isUnExported(name string) bool {
	return reLower.Match([]byte(name))
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// TestFilePath returns go test file of given source file.
func TestFilePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", nil
	}

	if !strings.HasSuffix(path, ".go") {
		return "", fmt.Errorf("%s is not go file", path)
	}

	if strings.HasSuffix(path, "_test.go") {
		return "", fmt.Errorf("%s is go test file", path)
	}

	return strings.Replace(path, ".go", "_test.go", -1), nil
}

func SrcFilePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", nil
	}

	if !strings.HasSuffix(path, "_test.go") {
		return "", fmt.Errorf("%s is not go test file", path)
	}

	return strings.Replace(path, "_test.go", ".go", -1), nil
}
