// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package parse

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

type PkgInfo struct {
	AstPkg  *ast.Package
	DocPkg  *doc.Package
	Imports map[string]string
}

type TaskFunction struct {
	Name string
	Doc  string
}

const GwPackage = "go.justen.tech/goodwill/gw"
const TaskType = "Task"
const ErrorType = "error"

func Package(path string, files []string) (*PkgInfo, error) {
	fset := token.NewFileSet()
	var filter func(file os.FileInfo) bool
	if len(files) > 0 {
		fm := make(map[string]bool, len(files))
		for _, f := range files {
			fm[f] = true
		}

		filter = func(f os.FileInfo) bool {
			return fm[f.Name()]
		}
	}
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("go parse error in %s: %w", path, err)
	}
	switch len(pkgs) {
	case 1:
		var pkg *ast.Package
		for _, pkg = range pkgs {
		}
		return &PkgInfo{
			AstPkg: pkg,
			DocPkg: doc.New(pkg, "./", 0),
		}, nil
	case 0:
		return nil, fmt.Errorf("no importable packages found in %s", path)
	default:
		var names []string
		for name := range pkgs {
			names = append(names, name)
		}
		return nil, fmt.Errorf("multiple packages found in %s: %v", path, strings.Join(names, ", "))
	}
}

func TaskFunctions(path string, info *PkgInfo) ([]TaskFunction, error) {
	var funcs []TaskFunction
	funcNames := make(map[string]struct{})
	imports := make(map[string]string)
	for _, file := range info.AstPkg.Files {
		for _, is := range file.Imports {
			importPath, err := strconv.Unquote(is.Path.Value)
			if err != nil {
				return nil, fmt.Errorf("could not unquote %s: %w", is.Path.Value, err)
			}
			if is.Name != nil {
				imports[is.Name.Name] = importPath
				continue
			}
			pkg, err := build.Import(importPath, path, 0)
			if err != nil {
				return nil, fmt.Errorf("could not import %s: %w", is.Path.Value, err)
			}
			imports[pkg.Name] = importPath
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				name, ok := flowFunctionName(imports, fn)
				if ok {
					funcNames[name] = struct{}{}
				}
			}
		}
	}
	for _, f := range info.DocPkg.Funcs {
		if _, ok := funcNames[f.Name]; ok {
			funcs = append(funcs, TaskFunction{
				Name: f.Name,
				Doc:  oneLine(f.Doc),
			})
		}
	}
	return funcs, nil
}

func oneLine(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", " "))
}

func ident(node ast.Node) string {
	id, ok := node.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

func flowFunctionName(imports map[string]string, f *ast.FuncDecl) (string, bool) {
	params := f.Type.Params
	if len(params.List) != 1 {
		return "", false
	}
	p1 := params.List[0]
	sexpr, ok := p1.Type.(*ast.StarExpr)
	if !ok {
		return "", false
	}
	selexpr, ok := sexpr.X.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	if pkg := ident(selexpr.X); imports[pkg] != GwPackage {
		return "", false
	}
	if ident(selexpr.Sel) != TaskType {
		return "", false
	}
	// return
	returns := f.Type.Results
	if len(returns.List) != 1 {
		return "", false
	}
	r1 := returns.List[0]
	if ident(r1.Type) != ErrorType {
		return "", false
	}
	return ident(f.Name), true
}
