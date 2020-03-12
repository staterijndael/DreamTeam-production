package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var spaces = regexp.MustCompile(`\s`)

func parse(dirPath string) error {
	packages, err := parser.ParseDir(token.NewFileSet(), dirPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for k, p := range packages {
		packageName = k
		for _, f := range p.Files {
			parseFile(f)
		}

		break
	}

	return nil
}

func parseFile(f *ast.File) {
	for _, d := range f.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok == token.VAR {
			parseVar(genDecl)
		} else {
			parseType(genDecl)
		}

	}
}

func parseVar(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		_, ok = spec.Type.(*ast.MapType)
		if !ok {
			continue
		}

		for _, v := range spec.Comment.List {
			if strings.HasPrefix(v.Text, commentPrefix) {
				mapName = spec.Names[0].Name
				mapValueTypeName = spec.
					Names[0].
					Obj.
					Decl.(*ast.ValueSpec).
					Type.(*ast.MapType).
					Value.(*ast.FuncType).
					Results.
					List[0].
					Type.(*ast.Ident).
					Name
				return
			}
		}
	}
}

func parseType(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok || !hasPrefix(typeSpec.Comment) {
			continue
		}

		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		parseStruct(typeSpec)
	}
}

func hasPrefix(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}

	for _, v := range cg.List {
		if strings.HasPrefix(v.Text, commentPrefix) {
			return true
		}
	}

	return false
}

func parseStruct(decl *ast.TypeSpec) {
	data := &TypeData{
		Name:    decl.Name.Obj.Name,
		Aliases: make([]string, 1),
	}

	data.Aliases[0] = strings.ToLower(data.Name)
	typesInfo = append(typesInfo, data)
	parseComment(data, decl)
}

func parseComment(data *TypeData, decl *ast.TypeSpec) {
	for _, group := range decl.Comment.List {
		if !strings.HasPrefix(group.Text, commentPrefix+":") {
			continue
		}

		withoutPrefix := strings.Replace(group.Text, commentPrefix+":", "", -1)
		for _, alias := range spaces.Split(withoutPrefix, -1) {
			if replaced := spaces.ReplaceAllString(alias, ""); len(replaced) > 0 {
				data.Aliases = append(data.Aliases, strings.ToLower(replaced))
			}
		}
	}
}
