package noosexit

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "direct call to os.Exit in package main is not allowed",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		if strings.Contains(filename, "/.cache/") || strings.Contains(filename, "\\.cache\\") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if ident, ok := selector.X.(*ast.Ident); ok &&
				ident.Name == "os" && selector.Sel.Name == "Exit" {

				obj := pass.TypesInfo.Uses[ident]
				if pkgName, ok := obj.(*types.PkgName); ok && pkgName.Imported().Path() == "os" {
					pass.Reportf(callExpr.Pos(), "direct call to os.Exit in package main is not allowed")
				}
			}
			return true
		})
	}
	return nil, nil
}
