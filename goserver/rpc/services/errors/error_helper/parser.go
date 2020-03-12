package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func parse(fp string) error {
	astF, err := parser.ParseFile(token.NewFileSet(), fp, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	packageName = astF.Name.Name
	for _, d := range astF.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if !*createMap && *mapName == defaultMapName && decl.Tok == token.VAR {
			processVar(decl)
		}

		if decl.Tok == token.CONST {
			processConst(decl)
		}
	}

	return nil
}

func processConst(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if spec.Comment != nil {
			comment := spec.Comment.Text()
			if strings.HasPrefix(comment, commentPrefix) {
				errorsInfo = append(errorsInfo, &ErrorData{
					Name: spec.Names[0].Name,
					Text: strings.TrimPrefix(comment[:len(comment)-1], commentPrefix),
				})
			}
		}
	}
}

func processVar(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if !isErrorMap(spec.Type) {
			continue
		}

		mapName = &spec.Names[0].Name
		return
	}
}

func isErrorMap(value ast.Expr) bool {
	mapType, ok := value.(*ast.MapType)
	if !ok {
		return false
	}

	key, ok1 := mapType.Key.(*ast.Ident)
	val, ok2 := mapType.Value.(*ast.Ident)

	return ok1 && ok2 && key.Name == "int" && val.Name == "string"
}
