package utils

import (
	"io/ioutil"
	"strings"
)

func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func FixIdSyntax(id string) string {
	return strings.Replace(id, "/", ".", 1)
}

func ReadFile(pathToDoc string) ([]byte, error) {
	return ioutil.ReadFile(pathToDoc)
}
