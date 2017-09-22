package yaag

import (
	"go/token"
	"path/filepath"
	"os"
	"go/parser"
	"fmt"
	"go/ast"
	"bufio"
	"strings"
	"github.com/betacraft/yaag/yaag/models"
)

func ParseAnnotations(rootDir string) error {

	dirs := []string{}

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	for _, d := range dirs {
		fs := token.NewFileSet()
		packages, err := parser.ParseDir(fs, d, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing file %v", err)
			return err
		}

		for _, p := range packages {
			for _, f := range p.Files {
				ast.Inspect(f, func(node ast.Node) bool {
					switch t := node.(type) {
					case *ast.FuncDecl:
						if t.Type.Params.NumFields() == 2 {
							fields := t.Type.Params.List
							if checkFieldType(fields[0], "http.ResponseWriter") &&
								checkFieldType(fields[1], "http.Request") {
								fmt.Printf("processing function %s \n", t.Name)
								parseComment(t.Doc.Text())
								return false
							}
						}
					}
					return true
				})
			}
		}

	}
	return nil
}

func parseComment(c string) {
	m := models.MetaData{}
	s := bufio.NewScanner(strings.NewReader(c))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "@Description") {
			m.Description = strings.TrimSpace(strings.TrimPrefix(line, "@Description"))
		} else if strings.HasPrefix(line, "@Path") {
			m.Path = strings.TrimSpace(strings.TrimPrefix(line, "@Path"))
		}
	}
	if m.Path != "" {
		spec.MetaData = append(spec.MetaData, m)
	}
}

func checkFieldType(f *ast.Field, t string) bool {
	ts := strings.Split(t, ".")
	if len(ts) != 2 {
		return false
	}
	var sexp *ast.SelectorExpr
	switch e := f.Type.(type) {
	case *ast.StarExpr:
		switch s := e.X.(type) {
		case *ast.SelectorExpr:
			sexp = s
			break
		}
		break
	case *ast.SelectorExpr:
		sexp = e
		break
	}
	n1 := sexp.X.(*ast.Ident).Name
	n2 := sexp.Sel.Name
	if ts[0] == n1 && ts[1] == n2 {
		return true
	}
	return false
}
