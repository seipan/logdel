package logdel_test

import (
	"testing"

	"github.com/seipan/logdel"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, logdel.Analyzer, "a")
}
