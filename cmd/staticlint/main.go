package main

//Создайте свой multichecker, состоящий из:
//
//стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes;
//всех анализаторов класса SA пакета staticcheck.io;
//не менее одного анализатора остальных классов пакета staticcheck.io;
//двух или более любых публичных анализаторов на ваш выбор.
//
//Напишите и добавьте в multichecker собственный анализатор, запрещающий использовать прямой вызов os.Exit в
//функции main пакета main. При необходимости перепишите код своего проекта так, чтобы он удовлетворял данному анализатору.
//Поместите анализатор в директорию cmd/staticlint вашего проекта. Добавьте документацию в формате godoc,
//подробно опишите в ней механизм запуска multichecker, а также каждый анализатор и его назначение.
//Исходный код вашего проекта должен проходить статический анализ созданного multichecker.

// LAUNCH:
// go build -o /tmp/mcheck url-shortener/cmd/staticlint #gosetup
// /tmp/mcheck ./...

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

var ErrNoExitAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "check for direct usage of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.CallExpr: // os.Exit
				switch ce := x.Fun.(type) {
				case *ast.SelectorExpr:
					if ce.Sel.Name == "Exit" && fmt.Sprintf("%s", ce.X) == "os" {
						pass.Reportf(ce.Pos(), "using os.Exit in the main() function is undesirable")
					}
				}
			}
			return true
		})
	}

	return nil, nil
}

func main() {
	var mychecks = []*analysis.Analyzer{
		// check for direct usage of os.Exit
		ErrNoExitAnalyzer,
		// check consistency of Printf format strings and arguments
		printf.Analyzer,
		// check that struct field tags conform to reflect.StructTag.Get
		structtag.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	multichecker.Main(
		mychecks...,
	)
}
