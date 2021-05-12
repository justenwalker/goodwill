// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go.justen.tech/goodwill/internal/parse"
	"io"
	"strings"
)

const (
	gwPackage    = "go.justen.tech/goodwill/gw"
	valuePackage = "go.justen.tech/goodwill/gw/value"
)

func renderMain(w io.Writer, data ParsedData) error {
	file := NewFile("main")
	file.HeaderComment("+build goodwill")
	//func dieOnError(err error) {
	//	if err != nil {
	//		fmt.Fprintln(os.Stderr, err)
	//		os.Exit(1)
	//	}
	//}
	file.Func().Id("dieOnError").
		Params(Id("err").Error()).
		Block(
			If(Id("err").Op("!=").Nil()).Block(
				printlnStdErr(Id("err")),
				osExit(1),
			))
	//fmt.Printf("Usage: %s FUNCTION\n", os.Args[0])
	//fmt.Println("Functions:")
	usageCode := []Code{
		Qual("fmt", "Printf").Call(Lit("Usage: %s FUNCTION\n"), Qual("os", "Args").Index(Lit(0))),
		Qual("fmt", "Println").Call(Lit("Functions:")),
	}
	for _, fn := range data.Functions {
		// fmt.Println("\t{{ $func.Name }}\t{{ $func.Doc }}")
		usageCode = append(usageCode, Qual("fmt", "Println").
			Call(Lit(fmt.Sprintf("\t%s\t%s", fn.Name, fn.Doc))))
	}
	// os.Exit(128)
	usageCode = append(usageCode, osExit(128))
	//if len(os.Args) < 2 {
	//	fmt.Fprintln(os.Stderr, "Must provide a function name")
	//	usage()
	//}
	//funcName := os.Args[1]
	mainCode := []Code{
		If(Id("len").Call(Qual("os", "Args")).Op("<").Lit(2)).Block(
			printlnStdErr(Lit("Must provide a function name")),
			Id("usage").Call(),
		),
		Id("funcName").Op(":=").Qual("os", "Args").Index(Lit(1)),
	}
	var funcSwitchCases []Code
	for _, fn := range data.Functions {
		name := strings.ToLower(fn.Name)
		//case "{{ $func.Name }}": {
		//	dieOnError(gw.Run({{ $func.Name }}))
		//	return
		//}
		funcSwitchCases = append(funcSwitchCases, Case(Lit(name)).Block(
			jenRunTask(fn),
			Return()))
	}

	mainCode = append(mainCode,
		//ctx, cancel := context.WithCancel(context.Background())
		List(Id("ctx"), Id("cancel")).Op(":=").Qual("context", "WithCancel").Call(Qual("context", "Background").Call()),
		//defer cancel()
		Defer().Id("cancel").Call(),
		//sigCh := make(chan signal.Signal)
		Id("sigCh").Op(":=").Make(Id("chan").Qual("os", "Signal")),
		//signal.Notify(sigCh, os.Interrupt)
		Qual("os/signal", "Notify").Call(Id("sigCh"), Qual("os", "Interrupt")),
		//go func() {
		Go().Func().Params().Block(
			//	<-sigCh
			Op("<-").Id("sigCh"),
			//	cancel()
			Id("cancel").Call(),
		).Call()) //}()
	// switch  strings.ToLower(funcName) { ... cases ... }
	mainCode = append(mainCode, Switch(Qual("strings", "ToLower").Call(Id("funcName"))).Block(funcSwitchCases...))
	file.Func().Id("usage").Params().Block(usageCode...)
	//fmt.Fprintln(os.Stderr, "Unknown function name:", funcName)
	//usage()
	mainCode = append(mainCode, printlnStdErr(Lit("Unknown function name:"), Id("funcName")),
		Id("usage").Call())
	file.Func().Id("main").Params().Block(mainCode...)
	return file.Render(w)
}

func jenRunTask(fn parse.TaskFunction) Code {
	var block Code
	if fn.Context {
		// Func(ctx,ts)
		block = Id(fn.Name).Call(Id("ctx"), Id("ts"))
	} else {
		// Func(ts)
		block = Id(fn.Name).Call(Id("ts"))
	}
	if fn.OutVars {
		// return Func(...)
		block = Return(block)
	} else {
		// return nil, Func(...)
		block = Return(Nil(), block)
	}
	// dieOnError(gw.Run(gw.TaskRunnerFunc(func(ctx context.Context, ts *gw.Task) (map[string]value.Value,error) {<block>}))
	return Id("dieOnError").Call(Qual(gwPackage, "Run").Call(
		Id("ctx"),
		Qual(gwPackage, "TaskRunnerFunc").Parens(Func().Params(
			Id("ctx").Qual("context", "Context"),
			Id("ts").Add(Op("*")).Qual(gwPackage, "Task"),
		).Params(
			Map(String()).Qual(valuePackage, "Value"),
			Id("error"),
		).Block(block))))
}

func printlnStdErr(code ...Code) Code {
	code = append([]Code{Qual("os", "Stderr")}, code...)
	return Qual("fmt", "Fprintln").Call(code...)
}

func osExit(rc int) Code {
	return Qual("os", "Exit").Call(Lit(rc))
}
