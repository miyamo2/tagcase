package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

type vetCallback func([]analysis.Diagnostic)

type fmtCallback func(*astutil.Cursor, ast.Node)

func visitor[T vetCallback | fmtCallback](config *Config, callBack T) func(c *astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		typeSpec, ok := c.Node().(*ast.TypeSpec)
		if !ok {
			return true
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		replacedFields := make([]*ast.Field, 0, len(structType.Fields.List))
		diagnostics := make([]analysis.Diagnostic, 0, len(structType.Fields.List))

		for i, field := range structType.Fields.List {
			replacedField := &ast.Field{
				Doc:     field.Doc,
				Names:   field.Names,
				Type:    field.Type,
				Comment: field.Comment,
			}
			replacedFields = append(replacedFields, replacedField)
			if field.Tag == nil {
				continue
			}
			replacedField.Tag = &ast.BasicLit{
				Kind:  field.Tag.Kind,
				Value: field.Tag.Value,
			}
			structTag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			if slices.Contains(strings.Split(structTag.Get("tagcase"), ","), "ignore") {
				continue
			}
			for key, rule := range config.Tags {
				previousTagStr, ok := structTag.Lookup(key)
				if !ok {
					continue
				}
				target := previousTagStr
				if rule.Delimit != nil {
					split := strings.Split(previousTagStr, rule.Delimit.Delimiter)
					if len(split) <= rule.Delimit.Index {
						continue
					}
					target = split[rule.Delimit.Index]
				}
				replaced := FormatWithCase(config.Caser(), target, rule.Case)
				if target == replaced {
					continue
				}
				presentTagStr := strings.ReplaceAll(previousTagStr, target, replaced)

				fullyPreviousTagStr := fmt.Sprintf("%s:%q", key, previousTagStr)
				fullyPresentTagStr := fmt.Sprintf("%s:%q", key, presentTagStr)

				replacedFields[i].Tag.Value = strings.Replace(
					replacedFields[i].Tag.Value,
					fullyPreviousTagStr,
					fullyPresentTagStr, 1)

				posStart := strings.Index(field.Tag.Value, fullyPreviousTagStr) +
					strings.Index(fullyPreviousTagStr, target)
				diagnostics = append(diagnostics, analysis.Diagnostic{
					Pos:     field.Tag.Pos() + token.Pos(posStart),
					Message: fmt.Sprintf("%q should be %q in %s tags.", target, replaced, key),
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: fmt.Sprintf("Correct %q to %q", previousTagStr, replaced),
							TextEdits: []analysis.TextEdit{
								{
									Pos:     field.Tag.Pos() + token.Pos(posStart),
									End:     field.Tag.Pos() + token.Pos(posStart+len(target)),
									NewText: []byte(replaced),
								},
							},
						},
					},
				})
			}
		}
		switch f := any(callBack).(type) {
		case fmtCallback:
			structType.Fields.List = replacedFields
			f(c, typeSpec)
		case vetCallback:
			f(diagnostics)
		}
		return true
	}
}
