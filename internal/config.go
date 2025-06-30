package internal

import (
	"encoding"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ettle/strcase"
	"github.com/goccy/go-yaml"
)

// Config represents the configuration for the application.
type Config struct {
	Tags             map[string]Rule `yaml:"tags"       json:"tags"`
	Initialism       Initialism      `yaml:"initialism" json:"initialism"`
	caser            *strcase.Caser  `yaml:"-"          json:"-"`
	initialCaserOnce sync.Once
}

// Initialism defines the configuration for handling initialism in tags.
type Initialism struct {
	Enable  []string `yaml:"enable"  json:"enable"`
	Disable []string `yaml:"disable" json:"disable"`
}

// Rule defines a transformation rule for a specific target.
type Rule struct {
	// Case indicates the transformation case to apply
	Case Case `yaml:"case" json:"case"`

	// Delimit defines how to handle delimiters in the tag
	Delimit *Delimit `yaml:"delimit,omitempty" json:"delimit,omitempty"`
}

type Delimit struct {
	// Index of the name to be checked in the tag
	Index int `yaml:"index" json:"index"`

	// Delimiter in the tag
	Delimiter string `yaml:"delimiter" json:"delimiter"`
}

// compatibility checks to ensure that Case implements the encoding.TextUnmarshaler/encoding.TextMarshaler interfaces
var (
	_ encoding.TextUnmarshaler = (*Case)(nil)
	_ encoding.TextMarshaler   = (*Case)(nil)
)

type Case int

func (c Case) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (c *Case) UnmarshalText(text []byte) error {
	v, err := CaseFromString(strings.TrimSpace(string(text)))
	if err != nil {
		return err
	}
	*c = v
	return nil
}

func CaseFromString(s string) (Case, error) {
	v, ok := stringToCase[strings.TrimSpace(s)]
	if !ok {
		return 0, fmt.Errorf("unknown case: %s", s)
	}
	return v, nil
}

func (c Case) String() string {
	return caseToString[c]
}

const (
	_ Case = iota
	// SnakeCase separates words with underscores. e.g., "snake_case_example"
	SnakeCase
	// ScreamingSnakeCase separates words with underscores and uses uppercase letters. e.g., "SCREAMING_SNAKE_CASE"
	ScreamingSnakeCase
	// KebabCase separates words with hyphens. e.g., "kebab-case-example"
	KebabCase
	// ScreamingKebabCase separates words with hyphens and uses uppercase letters. e.g., "SCREAMING-KEBAB-CASE"
	ScreamingKebabCase
	// PascalCase uses uppercase for the first letter of each word. e.g., "PascalCaseExample"
	PascalCase
	// CamelCase uses uppercase for the first letter of each word except the first word. e.g., "camelCaseExample"
	CamelCase
)

const (
	SnakeCaseStr          = "snake_case"
	ScreamingSnakeCaseStr = "SNAKE_CASE"
	KebabCaseStr          = "kebab-case"
	ScreamingKebabCaseStr = "KEBAB-CASE"
	PascalCaseStr         = "PascalCase"
	CamelCaseStr          = "camelCase"
)

func FormatWithCase(caser *strcase.Caser, s string, c Case) string {
	switch c {
	case PascalCase:
		return caser.ToPascal(s)
	case SnakeCase:
		return caser.ToSnake(s)
	case KebabCase:
		return caser.ToKebab(s)
	case CamelCase:
		return caser.ToCamel(s)
	case ScreamingSnakeCase:
		return caser.ToSNAKE(s)
	case ScreamingKebabCase:
		return caser.ToKEBAB(s)
	default:
		return s
	}
}

var stringToCase = map[string]Case{
	SnakeCaseStr:          SnakeCase,
	ScreamingSnakeCaseStr: ScreamingSnakeCase,
	KebabCaseStr:          KebabCase,
	ScreamingKebabCaseStr: ScreamingKebabCase,
	PascalCaseStr:         PascalCase,
	CamelCaseStr:          CamelCase,
}

var caseToString = map[Case]string{
	SnakeCase:          SnakeCaseStr,
	ScreamingSnakeCase: ScreamingSnakeCaseStr,
	KebabCase:          KebabCaseStr,
	ScreamingKebabCase: ScreamingKebabCaseStr,
	PascalCase:         PascalCaseStr,
	CamelCase:          CamelCaseStr,
}

// InitConfig reads in Cfg file and ENV variables if set.
func InitConfig() (*Config, error) {
	var c Config
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to reading config file: %w", err)
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	_ = filepath.Walk("./", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, walkErr)
			return walkErr
		}
		if _, ok := configFileName[info.Name()]; !ok {
			return nil
		}
		data, err := os.ReadFile(info.Name())
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, &c); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing Cfg file:", err)
		}
		return filepath.SkipAll
	})
	return &c, nil
}

// Caser returns a strcase.Caser instance configured with the initialism settings.
func (c *Config) Caser() *strcase.Caser {
	c.initialCaserOnce.Do(func() {
		initialism := make(map[string]bool)
		for _, v := range c.Initialism.Enable {
			initialism[v] = true
		}
		for _, v := range c.Initialism.Disable {
			initialism[v] = false
		}
		c.caser = strcase.NewCaser(true, initialism, nil)
	})
	return c.caser
}
