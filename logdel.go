package logdel

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"strconv"
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

	files := pass.Fset.File(pass.Files[0].Pos())

	cmp, ok := getComment(pass, file)
	omp := getImportObj(pass)

	file, err := deleteLogfromAST(pass, file, ok, cmp, omp)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	f, err := os.CreateTemp("", "logdelout.go")

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

	if err := os.Rename(f.Name(), files.Name()); err != nil {
		log.Fatalln("Error:", err)
		return err
	}

	return nil
}

func deleteLogfromAST(pass *analysis.Pass, file *ast.File, importOk bool, cmp map[string]string, omp map[types.Object]bool) (*ast.File, error) {
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

		found, err := findLogInvocation(pass, cur, cmp, omp)
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

func findLogInvocation(pass *analysis.Pass, cr *astutil.Cursor, cmp map[string]string, omp map[types.Object]bool) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ExprStmt:
		switch x := node.X.(type) {
		case *ast.CallExpr:
			return findLogInvocationInCallExpr(pass, x, cr.Index(), cmp, omp)
		}
	case *ast.AssignStmt:
		for _, r := range node.Rhs {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findLogInvocationInCallExpr(pass, x, cr.Index(), cmp, omp)
			}
		}
	case *ast.ReturnStmt:
		for _, r := range node.Results {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findLogInvocationInCallExpr(pass, x, cr.Index(), cmp, omp)
			}
		}
	}
	return false, nil
}

func findLogInvocationInCallExpr(pass *analysis.Pass, callExpr *ast.CallExpr, idx int, cmp map[string]string, omp map[types.Object]bool) (bool, error) {
	types := pass.TypesInfo

	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		x2, ok := fun.X.(*ast.Ident)
		if !ok {
			return false, fmt.Errorf("this select-expr's X is not ident: %v", fun.X)
		}
		_, ok = omp[types.ObjectOf(fun.Sel)]

		if idx >= 0 && ok {
			pos := pass.Fset.Position(x2.Pos())
			c, ok := cmp[pos.Filename+"_"+strconv.Itoa(pos.Line)]
			if ok {
				if strings.Contains(c, "nocheck:thislog") {
					return false, nil
				}
			}
			log.Println(pos.Filename + " Line" + strconv.Itoa(pos.Line) + "  " + x2.Name + "." + fun.Sel.Name + " id deleted")
		}

		return idx >= 0 && ok, nil

	default:
		return false, nil
	}
}

func getComment(pass *analysis.Pass, file *ast.File) (map[string]string, bool) {
	var mp = make(map[string]string)
	var ok bool

	for _, cg := range file.Comments {
		for _, c := range cg.List {
			pos := pass.Fset.Position(c.Pos())
			mp[pos.Filename+"_"+strconv.Itoa(pos.Line)] = c.Text
			if strings.Contains(c.Text, "nocheck:thislog") {
				ok = true
			}
		}
	}

	return mp, ok
}

func getImportObj(pass *analysis.Pass) map[types.Object]bool {
	var mp = make(map[types.Object]bool)
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
