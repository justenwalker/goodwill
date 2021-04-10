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
	return Id("dieOnError").Call(Qual(parse.GwPackage, "Run").Call(Id(fn.Name)))
}

func printlnStdErr(code ...Code) Code {
	code = append([]Code{Qual("os", "Stderr")}, code...)
	return Qual("fmt", "Fprintln").Call(code...)
}

func osExit(rc int) Code {
	return Qual("os", "Exit").Call(Lit(rc))
}
