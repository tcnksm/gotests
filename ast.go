package main

import (
	"go/ast"
	"go/token"
	"log"
	"os"
)

// DebugAst runs ast.Print only when debug mode is enabled.
func DebugAst(fset *token.FileSet, x interface{}) {
	if os.Getenv(EnvDebug) != "" {
		log.Println("[DEBUG] Ast")
		ast.Print(fset, x)
	}
}

// NewTestFuncDecl creates a new FuncDecl for starndard testing
// without position.
func NewTestFuncDecl(name string) *ast.FuncDecl {

	ident := ast.NewIdent(name)
	ident.Obj = ast.NewObj(ast.Fun, name)

	identVarT := ast.NewIdent("t")
	identVarT.Obj = ast.NewObj(ast.Var, "t")

	// params are params for func
	params := &ast.FieldList{
		List: []*ast.Field{
			{
				// t
				Names: []*ast.Ident{
					identVarT,
				},

				// *testing.T
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("testing"),
						Sel: ast.NewIdent("T"),
					},
				},
			},
		},
	}

	funcType := &ast.FuncType{
		Params: params,
	}

	return &ast.FuncDecl{
		Name: ident,
		Type: funcType,
		Body: &ast.BlockStmt{},
	}
}
