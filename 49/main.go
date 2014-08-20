package main

import (
	"go/ast"
	"go/doc"
	"os"
	"text/template"

	"github.com/shurcooL/go/gists/gist5504644"
	"github.com/shurcooL/go/gists/gist5639599"
)

func findType(dpkg *doc.Package, name string) *doc.Type {
	for _, t := range dpkg.Types {
		if t.Name == name {
			return t
		}
	}
	return nil
}

var t = template.Must(template.New("repo").Parse(`package whatever

import (
	"fmt"
)

type debug{{.InterfaceName}} struct {
	real {{.InterfaceName}}
}
{{range .Methods}}
func (this *debug{{$.InterfaceName}}) {{.Name}}{{.Something}} {
	fmt.Println("{{$.InterfaceName}}.{{.Name}}")
	this.real.{{.Name}}()
	return
}
{{end}}`))

func main() {
	//goon.Dump(t.Doc)
	//gist5639599.PrintlnAstBare(t.Decl)
	/*err = ast.Fprint(os.Stdout, token.NewFileSet(), interfaceType, nil)
	if err != nil {
		panic(err)
	}*/

	/*for _, m := range methods.List {
		fmt.Print(m.Names[0].Name + "		")
		gist5639599.PrintlnAstBare(m.Type)
	}*/

	x, err := newGen("github.com/russross/blackfriday", "Renderer")
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, x)
	if err != nil {
		panic(err)
	}
}

func newGen(importPath, interfaceName string) (*gen, error) {
	dpkg, err := gist5504644.GetDocPackage(gist5504644.BuildPackageFromImportPath(importPath))
	if err != nil {
		return nil, err
	}

	t := findType(dpkg, interfaceName)

	methods := t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods

	return &gen{
		dpkg:          dpkg,
		methods:       methods,
		InterfaceName: interfaceName,
	}, nil
}

type gen struct {
	dpkg    *doc.Package
	methods *ast.FieldList

	InterfaceName string
}

func (this *gen) Methods() <-chan Method {
	out := make(chan Method)
	go func() {
		for _, m := range this.methods.List {
			out <- Method{m.Names[0].Name, gist5639599.SprintAstBare(m.Type)}
		}
		close(out)
	}()
	return out
}

type Method struct {
	Name      string
	something string
}

func (this Method) Something() string {
	return this.something[4:]
}
