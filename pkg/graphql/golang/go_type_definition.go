package golang

import (
	"unicode"
)

type GoTypeDefinition interface {
	UnqualifiedName() string
	Definition() string
}


func lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}
