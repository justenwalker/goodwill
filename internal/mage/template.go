package mage

import (
	"fmt"
	"io"
	"path"
	"text/template"
)

type ConcordRuntime int

const (
	ConcordRuntimeDefault ConcordRuntime = iota
	ConcordRuntimeV1
	ConcordRuntimeV2
)

const DefaultGoVersion = "1.17"

type ConcordParams struct {
	Dependencies bool
	Runtime      ConcordRuntime
	Version      string
	GoVersion    string
	UseDocker    bool
}

// GenerateConcordYaml generates an example concord.yml file with the given parameters
func GenerateConcordYaml(w io.Writer, params ConcordParams) (err error) {
	if params.Runtime == ConcordRuntimeDefault {
		params.Runtime = ConcordRuntimeV1
	}
	if params.GoVersion == "" {
		params.GoVersion = DefaultGoVersion
	}
	name := fmt.Sprintf("concord-v%d.yml.gotmpl", params.Runtime)
	tfile := path.Join("files", name)
	tpl, err := template.ParseFS(Files, tfile)
	if err != nil {
		return err
	}
	return tpl.Execute(w, params)
}
