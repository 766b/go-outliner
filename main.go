package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var (
	flgVer = flag.Bool("version", false, "Version")

	version = "0.3.0"
)

type Declaration struct {
	Label    string        `json:"label,omitempty"`
	Receiver string        `json:"receiver,omitempty"`
	Type     string        `json:"type,omitempty"`
	File     string        `json:"file,omitempty"`
	Line     int           `json:"line,omitempty"`
	Children []Declaration `json:"children,omitempty"`
}

func main() {
	flag.Parse()
	args := flag.Args()
	if *flgVer {
		DisplayVersion()
		return
	}

	if len(args) < 1 {
		return
	}
	fset := token.NewFileSet()
	asts, _ := parser.ParseDir(fset, args[0], nil, parser.ParseComments)
	var decls = make([]Declaration, 0)
	for _, a := range asts {
		for filename, file := range a.Files {
			for _, decl := range file.Decls {
				switch s := decl.(type) {
				case *ast.GenDecl:
					for _, s := range s.Specs {
						switch spec := s.(type) {
						case *ast.ValueSpec:
							for _, v := range spec.Names {
								decls = append(decls, Declaration{
									Label: v.Name,
									File:  filename,
									Type:  fmt.Sprintf("%s", v.Obj.Kind.String()),
									Line:  fset.Position(v.Pos()).Line,
								})
							}
						case *ast.TypeSpec:
							decls = append(decls, Declaration{
								Label: spec.Name.String(),
								File:  filename,
								Type:  "type",
								Line:  fset.Position(spec.Pos()).Line,
							})
						}
					}
				case *ast.FuncDecl:
					decls = append(decls, Declaration{
						Label:    s.Name.String(),
						Type:     "func",
						File:     filename,
						Receiver: getReceiver(fset, s.Recv),
						Line:     fset.Position(s.Pos()).Line,
					})
				}
			}
		}
	}

	json.NewEncoder(os.Stdout).Encode(decls)
}

func DisplayVersion() {
	fmt.Println("Version:", version)
}

func getReceiver(f *token.FileSet, l *ast.FieldList) string {
	if l == nil {
		return ""
	}

	x := l.List[0].Type
	for {
		switch val := x.(type) {
		case *ast.ArrayType:
			x = val.Elt
		case *ast.StarExpr:
			x = val.X
		case *ast.SelectorExpr:
			return fmt.Sprintf("%s.%s", val.X, val.Sel.Name)
		case *ast.StructType:
			return fmt.Sprintf("%s", val)
		case *ast.InterfaceType:
			return fmt.Sprintf("%s", val)
		case *ast.MapType:
			return fmt.Sprintf("map[%s]%s", val.Key, val.Value)
		default:
			return fmt.Sprintf("%s", val)
		}
	}
}
