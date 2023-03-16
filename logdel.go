package logdel

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func Run(filename string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	_, ok := getComment(file)

	file, err = deleteLogfromAST(file, ok)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	f, err := ioutil.TempFile("", "logdelout.go")
	fpath := f.Name()
	fmt.Println(fpath)

	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}

	writer := bufio.NewWriter(f)

	if err := format.Node(writer, fset, file); err != nil {
		log.Fatalln("Error:", err)
		return err
	}
	//log.Println(f)
	writer.Flush()
	f.Close()
	
	if err := os.Rename(f.Name(), filename); err != nil {
		log.Fatalln("Error:", err)
		return err
	}

	return nil
}

func deleteLogfromAST(file *ast.File, importOk bool) (*ast.File, error) {
	file, ok := astutil.Apply(file, func(cur *astutil.Cursor) bool {
		if !importOk {
			found, err := findLogImportInImportSpec(cur)
			if err != nil {
				return true
			}
			if found {
				cur.Delete()
				return true
			}
		}

		// if c.Node belongs to ExprStmt, remove callExpr for dl
		found, err := findLogInvocationInStmt(cur)
		if err != nil {
			return true
		}
		if found {
			cur.Delete()
			return true
		}

		// if return false, traversing is stopped immediately
		return true
	}, nil).(*ast.File)

	if !ok {
		return nil, errors.New("failed type conversion from any to *ast.File")
	}

	return file, nil
}

func findLogImportInImportSpec(cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ImportSpec:
		return cr.Index() >= 0 && node.Path.Value == "log", nil
	}
	return false, nil
}

func findLogInvocationInStmt(cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ExprStmt:
		switch x := node.X.(type) {
		case *ast.CallExpr:
			return findDlInvocationInCallExpr(x, cr.Index())
		}
	case *ast.AssignStmt:
		for _, r := range node.Rhs {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findDlInvocationInCallExpr(x, cr.Index())
			}
		}
	case *ast.ReturnStmt:
		for _, r := range node.Results {
			switch x := r.(type) {
			case *ast.CallExpr:
				return findDlInvocationInCallExpr(x, cr.Index())
			}
		}
	}
	return false, nil
}

func findDlInvocationInCallExpr(callExpr *ast.CallExpr, idx int) (bool, error) {
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		x2, ok := fun.X.(*ast.Ident)
		if !ok {
			return false, fmt.Errorf("x2 is not *ast.Ident: %v", fun.X)
		}

		// check node is in a slice
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
