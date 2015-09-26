package utils

import "strings"

func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func FixIdSyntax(id string) string {
	return strings.Replace(id, "/", ".", 1)
}
