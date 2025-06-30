package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/miyamo2/tagcase/internal"
	"golang.org/x/tools/go/analysis"
)

type plugin struct {
	config *internal.Config
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		internal.Analyzer(p.config),
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func new(input any) (register.LinterPlugin, error) {
	config, err := register.DecodeSettings[internal.Config](input)
	if err != nil {
		return nil, err
	}
	return &plugin{
		config: &config,
	}, nil
}

// compatibility check
var _ register.LinterPlugin = (*plugin)(nil)

func init() {
	register.Plugin("tagcase", new)
}
