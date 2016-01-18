package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Parse(filename string, rd io.Reader) (*GoFile, error) {

	// Store src as []byte for doDiff
	src, err := ioutil.ReadAll(rd)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("readall: %s", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}
	DebugAst(fset, f)

	var funcs []string
	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			Debugf("FuncDecl: %#v", x.Name)
			funcs = append(funcs, x.Name.Name)
		}
		return true
	})

	return &GoFile{
		PackageName: f.Name.Name,
		FileName:    filename,
		Src:         src,
		Funcs:       funcs,
		FSet:        fset,
		AstFile:     f,
	}, nil
}

func ParseFile(path string) (*GoFile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(path, ".go") {
		return nil, fmt.Errorf("%s is not go file", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return Parse(fi.Name(), f)
}
