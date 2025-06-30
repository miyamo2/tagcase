package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"iter"
	"os"
	"slices"
	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func Diff(config *Config, dir, source string) (string, error) {
	visitor_ := visitor[fmtCallback](
		config,
		func(c *astutil.Cursor, n ast.Node) {
			c.Replace(n)
		})
	var (
		previous bytes.Buffer
		present  bytes.Buffer
		result   bytes.Buffer
	)
	for info := range fileInfoIter(dir, source) {
		if info.err != nil {
			return "", info.err
		}
		err := printerConfig.Fprint(&previous, info.fset, info.f)
		if err != nil {
			return "", fmt.Errorf("failed to print previous file: %w", err)
		}
		info.f = astutil.Apply(info.f, nil, visitor_).(*ast.File)
		err = printerConfig.Fprint(&present, info.fset, info.f)
		if err != nil {
			return "", fmt.Errorf("failed to print present file: %w", err)
		}
		if diff := cmp.Diff(previous.String(), present.String()); diff != "" {
			result.WriteString(fmt.Sprintf("%s:\n", info.fset.File(info.f.Pos()).Name()))
			result.WriteString(fmt.Sprintf("%s\n", diff))
		}
		previous.Reset()
		present.Reset()
	}
	return result.String(), nil
}

func Format(config *Config, dir, source string) error {
	visitor_ := visitor[fmtCallback](
		config,
		func(c *astutil.Cursor, n ast.Node) {
			c.Replace(n)
		})
	for info := range fileInfoIter(dir, source) {
		if info.err != nil {
			return info.err
		}
		info.f = astutil.Apply(info.f, nil, visitor_).(*ast.File)
		file, err := os.Create(info.fset.File(info.f.Pos()).Name())
		if err != nil {
			return fmt.Errorf(
				"failed to open file %s for writing: %w",
				info.fset.File(info.f.Pos()).Name(),
				err,
			)
		}
		err = printerConfig.Fprint(file, info.fset, info.f)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf(
				"failed to close file %s: %w",
				info.fset.File(info.f.Pos()).Name(),
				err,
			)
		}
	}
	return nil
}

type fileInfo struct {
	fset *token.FileSet
	f    *ast.File
	err  error
}

func fileInfoIter(dir, source string) iter.Seq[*fileInfo] {
	fs := token.NewFileSet()
	if source != "" {
		return func(yield func(*fileInfo) bool) {
			f, err := parser.ParseFile(fs, source, nil, parser.AllErrors|parser.ParseComments)
			if err != nil {
				yield(&fileInfo{
					err: fmt.Errorf("failed to parse file %s: %w", source, err),
				})
				return
			}
			yield(&fileInfo{fset: fs, f: f})
		}
	}
	cfg := &packages.Config{
		Fset: fs,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax,
	}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return func(yield func(*fileInfo) bool) {
			yield(&fileInfo{
				err: fmt.Errorf("failed to load packages from %s: %w", source, err),
			})
		}
	}
	idx := slices.IndexFunc(pkgs, func(pkg *packages.Package) bool {
		return !strings.HasSuffix(pkg.Name, "_test")
	})
	if idx < 0 {
		return func(yield func(*fileInfo) bool) {
			yield(&fileInfo{
				err: fmt.Errorf("no non-test package found in %s", source),
			})
		}
	}
	return func(yield func(*fileInfo) bool) {
		for _, f := range pkgs[idx].Syntax {
			if !yield(&fileInfo{fset: pkgs[idx].Fset, f: f}) {
				return
			}
		}
	}
}

var printerConfig = &printer.Config{
	Tabwidth: 8,
	Mode:     printer.UseSpaces | printer.TabIndent,
}
