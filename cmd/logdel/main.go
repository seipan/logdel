package main

import (
	"github.com/seipan/logdel"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(logdel.Analyzer) }
