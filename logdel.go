package logdel

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

func Run(filename string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	f, err := ioutil.TempFile("", "logdel.go")
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}
	defer f.Close()

	return nil
}
