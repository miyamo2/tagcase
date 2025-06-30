package internal

import (
	"flag"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

func runner(config *Config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		visitor_ := visitor[vetCallback](config, func(diagnostics []analysis.Diagnostic) {
			for _, diagnostic := range diagnostics {
				pass.Report(diagnostic)
			}
		})
		for _, file := range pass.Files {
			astutil.Apply(file, nil, visitor_)
		}
		return nil, nil
	}
}

func Analyzer(config *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:       "tagcase",
		Doc:        "check that struct tags are in the correct case",
		Run:        runner(config),
		ResultType: nil,
		Flags:      flag.FlagSet{},
	}
}
