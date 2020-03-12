package main

import "text/template"

type TemplateParams struct {
	MapName          string
	TypesInfo        []*TypeData
	Package          string
	MapValueTypeName string
}

var (
	tmpl = template.Must(template.New("map").Parse(`// Code generated by notifier; DO NOT EDIT.
{{ $mvtn := .MapValueTypeName }}
package {{ .Package }}

func init() {
	{{ .MapName }} = map[string]func () {{ $mvtn }}{
		{{ range $i, $el := .TypesInfo -}}
		{{ range $i2, $inner := $el.Aliases -}}
		"{{ $inner }}": func () {{ $mvtn }} { return new({{ $el.Name }})},
		{{ end }}
		{{- end }}
	}
}
`))
)
