// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"go.justen.tech/goodwill/internal/parse"
	"text/template"
)

type TemplateData struct {
	Package   string
	Functions []parse.TaskFunction
}

func mainTemplate() *template.Template {
	return template.Must(template.New("main.go").Parse(`// +build goodwill

package {{ .Package }}

import (
	"os"
	"fmt"
	"strings"
	"go.justen.tech/goodwill/gw"
)

func dieOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Usage: %s FUNCTION\n", os.Args[0])
	fmt.Println("Functions:")
	{{- range $func := .Functions }}
	fmt.Println("\t{{ $func.Name }}\t{{ $func.Doc }}")
	{{- end }}
	os.Exit(128)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must provide a function name")
		usage()
	}
	funcName := os.Args[1]
{{- range $func := .Functions }}
	if strings.EqualFold(funcName, "{{ $func.Name }}") {
		dieOnError(gw.Run({{ $func.Name }}))
		return
	}
{{- end }}
	fmt.Fprintln(os.Stderr, "Unknown function name:", funcName)
	usage()
}
`))
}
