package main

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMapName = "notificationsMap"
	commentPrefix  = "//notifier"
	filenameSuffix = "_notifier_gen.go"
)

var (
	packageName      string
	mapName          = defaultMapName
	mapValueTypeName string
	typesInfo        = make([]*TypeData, 0)
)

type TypeData struct {
	Name    string
	Aliases []string
}

func main() {
	errHandler := func(err error) {
		panic(err)
	}

	filename := os.Getenv("GOFILE")
	if len(os.Args) > 1 && strings.HasSuffix(os.Args[len(os.Args)-1], ".go") {
		filename = os.Args[len(os.Args)-1]
	}

	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		errHandler(err)
		return
	}

	err = parse(dir)
	if err != nil {
		errHandler(err)
		return
	}

	fp := filepath.Join(dir, filename[:len(filename)-3]+filenameSuffix)
	file, err := os.Create(fp)
	if err != nil {
		errHandler(err)
		return
	}

	err = tmpl.Execute(file, TemplateParams{
		Package:          packageName,
		MapName:          mapName,
		TypesInfo:        typesInfo,
		MapValueTypeName: mapValueTypeName,
	})

	if err != nil {
		errHandler(err)
		return
	}
}
