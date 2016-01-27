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

func parse(filename string, rd io.Reader) (*GoFile, error) {

	// Store src as []byte for doDiff
	srcBytes, err := ioutil.ReadAll(rd)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("readall: %s", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, srcBytes, 0)
	if err != nil {
		return nil, err
	}
	DebugAst(fset, f)

	var funcs []string
	var methods []*Method
	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			Debugf("FuncDecl: %#v", x.Name)
			// receiver (methods) or nil (functions)
			if x.Recv == nil {
				funcs = append(funcs, x.Name.Name)
				return true
			}

			fields := x.Recv.List
			if len(fields) != 1 {
				// Is this happend ..?
				return true
			}

			field := fields[0]
			t := field.Type
			var recvName string
			switch x2 := t.(type) {
			case *ast.StarExpr:
				switch x3 := x2.X.(type) {
				case *ast.Ident:
					recvName = x3.Name
				}
			case *ast.Ident:
				recvName = x2.Name
			default:
				// Should not reach here...
				return false
			}

			methods = append(methods, &Method{
				RecvName: recvName,
				Name:     x.Name.Name,
			})
		}
		return true
	})

	Debugf("Funcs: %#v", methods)
	Debugf("Methods: %#v", methods)

	return &GoFile{
		PackageName: f.Name.Name,
		FileName:    filename,
		SrcBytes:    srcBytes,
		Funcs:       funcs,
		Methods:     methods,
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

	return parse(fi.Name(), f)
}
