package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMapName = "errorsMap"
	commentPrefix  = "eh:"
	filenameSuffix = "_errors_mapper_gen.go"
)

var (
	packageName string
	errorsInfo  = make([]*ErrorData, 0)
	createMap   *bool
	mapName     *string
	logCount    *bool
	logList     *bool
)

type ErrorData struct {
	Name string
	Text string
}

func main() {
	errHandler := func(err error) {
		panic(err)
	}

	logList = flag.Bool("log-list", true, "should log list of generated errors or not")
	logCount = flag.Bool("log-count", true, "should count generated errors")
	createMap = flag.Bool("create-map", true, "should generate map or not")
	mapName = flag.String("map-name", defaultMapName, "name of errors map")
	flag.Parse()

	filename := os.Getenv("GOFILE")
	if len(os.Args) > 1 && strings.HasSuffix(os.Args[len(os.Args)-1], ".go") {
		filename = os.Args[len(os.Args)-1]
	}

	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		errHandler(err)
		return
	}

	err = parse(filepath.Join(dir, filename))
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
		Package:    packageName,
		CreateMap:  *createMap,
		MapName:    *mapName,
		ErrorsInfo: errorsInfo,
	})

	if err != nil {
		errHandler(err)
		return
	}

	if *logCount {
		fmt.Printf("amount of generated errors: %d. ", len(errorsInfo))
	}

	if *logList {
		fmt.Print("generated errors: \n")
		for _, info := range errorsInfo {
			if *logList {
				fmt.Printf("\t%s: %s\n", info.Name, info.Text)
			}
		}
	} else {
		fmt.Println()
	}
}
