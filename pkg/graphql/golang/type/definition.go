package _type

import (
	"unicode"
)

type Definition interface {
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
