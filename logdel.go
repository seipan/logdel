package logdel

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "logdel",
	Doc:  "logdel is ......",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	for _, v := range pass.Files {
		RunDeleteLog(pass, v)
	}

	return nil, nil
}

func RunDeleteLog(pass *analysis.Pass, file *ast.File) error {
	fset := token.NewFileSet()

	// file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	// if err != nil {
	// 	log.Fatalln("Error:", err)
	// 	return nil
	// }
	mp, ok := getComment(file)

	for i, _ := range mp {
		log.Println(i)
	}

	file, err := deleteLogfromAST(file, ok)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	f, err := ioutil.TempFile("", "logdelout.go")

	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	writer := bufio.NewWriter(f)

	if err := format.Node(writer, fset, file); err != nil {
		log.Fatalln("Error:", err)
		return err
	}

	writer.Flush()
	f.Close()

	log.Println(file.Name.String())
	log.Println(pass.String())

	if err := os.Rename(f.Name(), file.Name.String()+".go"); err != nil {
		log.Fatalln("Error:", err)
		return err
	}

	return nil
}

func deleteLogfromAST(file *ast.File, importOk bool) (*ast.File, error) {
	file, ok := astutil.Apply(file, func(cur *astutil.Cursor) bool {
		if !importOk {
			found, err := findLogImport(cur)
			if err != nil {
				return true
			}
			if found {
				cur.Delete()
				return true
			}
		}

		found, err := findLogInvocation(cur)
		if err != nil {
			return true
		}
		if found {
			cur.Delete()
			return true
		}

		return true
	}, nil).(*ast.File)

	if !ok {
		return nil, errors.New("failed type conversion from any to *ast.File")
	}

	return file, nil
}

func findLogImport(cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ImportSpec:
		return cr.Index() >= 0 && node.Path.Value == `"log"`, nil
	}
	return false, nil
}

func findLogInvocation(cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ExprStmt:
		switch x := node.X.(type) {
		case *ast.CallExpr:
			return findLogInvocationInCallExpr(x, cr.Index())
		}
	case *ast.AssignStmt:
		for _, r := range node.Rhs {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findLogInvocationInCallExpr(x, cr.Index())
			}
		}
	case *ast.ReturnStmt:
		for _, r := range node.Results {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findLogInvocationInCallExpr(x, cr.Index())
			}
		}
	}
	return false, nil
}

func findLogInvocationInCallExpr(callExpr *ast.CallExpr, idx int) (bool, error) {
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		x2, ok := fun.X.(*ast.Ident)
		if !ok {
			return false, fmt.Errorf("this select-expr's X is not ident: %v", fun.X)
		}

		if idx >= 0 && "log" == x2.Name {
			log.Println(x2.Pos())
			log.Println(x2.Name + "." + fun.Sel.Name)
		}

		return idx >= 0 && "log" == x2.Name, nil

	default:
		return false, nil
	}
}

func getComment(file *ast.File) (map[token.Pos]string, bool) {
	var mp = make(map[token.Pos]string)
	var ok bool

	for _, cg := range file.Comments {
		for _, c := range cg.List {
			pos := c.Pos()
			mp[pos] = c.Text
			if strings.Contains(c.Text, "nocheck:thislog") {
				ok = true
			}
		}
	}

	return mp, ok
}

func getImportObj(pass *analysis.Pass) map[types.Object]bool {
	var mp map[types.Object]bool
	pkgs := pass.Pkg.Imports()
	obj := analysisutil.LookupFromImports(pkgs, "log", "Print")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Println")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Printf")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Fatal")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Fatalln")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Fatalf")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Panicf")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Panic")
	mp[obj] = true
	obj = analysisutil.LookupFromImports(pkgs, "log", "Panicln")
	mp[obj] = true

	return mp

}
