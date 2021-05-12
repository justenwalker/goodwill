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
	Name    string
	Doc     string
	Context bool
	OutVars bool
}

type ByName []TaskFunction

func (fs ByName) Len() int {
	return len(fs)
}

func (fs ByName) Less(i, j int) bool {
	return fs[i].Name < fs[j].Name
}

func (fs ByName) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

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
	taskFuncs := make(map[string]TaskFunction)
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
				tf, ok := taskFunction(imports, fn)
				if ok {
					taskFuncs[tf.Name] = tf
				}
			}
		}
	}
	var funcs []TaskFunction
	for _, f := range info.DocPkg.Funcs {
		if tf, ok := taskFuncs[f.Name]; ok {
			tf.Doc = oneLine(f.Doc)
			funcs = append(funcs, tf)
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

func taskFunction(imports map[string]string, f *ast.FuncDecl) (TaskFunction, bool) {
	fi := &functionInspector{
		f:       f,
		imports: imports,
	}
	return fi.function()
}

type functionInspector struct {
	f       *ast.FuncDecl
	imports map[string]string
}

func (i *functionInspector) function() (tf TaskFunction, ok bool) {
	params := i.f.Type.Params.List
	switch len(params) {
	case 1:
		if !i.isTaskType(params[0]) {
			return
		}
	case 2:
		if !i.isContextType(params[0]) {
			return
		}
		tf.Context = true
		if !i.isTaskType(params[1]) {
			return
		}
	default:
		return
	}
	returns := i.f.Type.Results.List
	switch len(returns) {
	case 1:
		if !i.isErrType(returns[0]) {
			return
		}
	case 2:
		if !i.isOutVars(returns[0]) {
			return
		}
		tf.OutVars = true
		if !i.isErrType(returns[1]) {
			return
		}
	default:
		return
	}
	tf.Name = ident(i.f.Name)
	return tf, true
}

func (i *functionInspector) isTaskType(field *ast.Field) bool {
	sexpr, ok := field.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	return i.isSelectorType(sexpr.X, "go.justen.tech/goodwill/gw", "Task")
}

func (i *functionInspector) isSelectorType(expr ast.Expr, pkg string, name string) bool {
	selexpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if p := ident(selexpr.X); i.imports[p] != pkg {
		return false
	}
	if ident(selexpr.Sel) != name {
		return false
	}
	return true

}

func (i *functionInspector) isErrType(field *ast.Field) bool {
	return ident(field.Type) == "error"
}

func (i *functionInspector) isOutVars(field *ast.Field) bool {
	mtype, ok := field.Type.(*ast.MapType)
	if !ok {
		return false
	}
	// key is a string
	ktype, ok := mtype.Key.(*ast.Ident)
	if !ok {
		return false
	}
	if ktype.Name != "string" {
		return false
	}
	return i.isSelectorType(mtype.Value, "go.justen.tech/goodwill/gw/value", "Value")
}

func (i *functionInspector) isContextType(field *ast.Field) bool {
	return i.isSelectorType(field.Type, "context", "Context")
}
