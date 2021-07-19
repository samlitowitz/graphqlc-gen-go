package golang

import (
	"go/ast"
	"unicode"
)

type GoTypeDefinition interface {
	GenDecl() *ast.GenDecl
	Ident() *ast.Ident
	File() *ast.File
}


func lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// package {name, alias} -> types
