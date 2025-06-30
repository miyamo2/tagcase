package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var configFile string

var (
	writeFlag   bool
	diffFlag    bool
	versionFlag bool
	initFlag    bool
)

var (
	Version  = "develop"
	Revision = "unknown"
)

// ErrNoFilesProvided is returned when no files are provided for formatting or diffing.
var ErrNoFilesProvided = errors.New("no files provided")

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tagcase",
		Short: "Checking and standardizing the case conventions used in naming Go struct tags.",
		Long:  `Checking and standardizing the case conventions used in naming Go struct tags.`,
		RunE:  run,
	}
	cmd.PersistentFlags().
		StringVar(&configFile, "config", "", "config file (default is <working directory>/.tagcase.yaml)")
	cmd.Flags().
		BoolVarP(&writeFlag, "write", "w", false, "Correct formatting issues and save changes to files instead of just reporting them.")
	cmd.Flags().
		BoolVarP(&diffFlag, "diff", "d", false, "Show the diff of the changes instead of writing them to files.")
	cmd.Flags().BoolVar(&versionFlag, "version", false, `Print the version of tagcase.`)
	cmd.Flags().BoolVar(&initFlag, "init", false, `Initialize a new configuration file.`)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	switch {
	case versionFlag:
		_, _ = fmt.Fprintf(
			cmd.OutOrStdout(),
			"tagcase version %q, revision %q\n",
			Version,
			Revision,
		)
	case initFlag:
		_, err := os.ReadFile(".tagcase.yaml")
		if err == nil {
			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				".tagcase.yaml already exists, skipping initialization.\n",
			)
			return nil
		}
		if !os.IsNotExist(err) {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error reading .tagcase.yaml: %v\n", err)
			os.Exit(1)
		}
		template := generateConfigTemplate()
		file, err := os.Create(".tagcase.yaml")
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error creating .tagcase.yaml: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			_ = file.Close()
		}()
		_, err = file.WriteString(template)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error writing to .tagcase.yaml: %v\n", err)
			os.Exit(1)
		}
		_, _ = fmt.Fprintf(
			cmd.OutOrStdout(),
			"Initialized .tagcase.yaml with default configuration.\n",
		)
	case diffFlag:
		config, err := InitConfig()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "failed to load config: %v\n", err)
		}
		if len(args) == 0 {
			fmt.Fprintf(cmd.OutOrStderr(), "no files provided\n")
			return nil
		}
		dir, file := filePathSplit(args[0])
		diff, err := Diff(config, dir, file)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error generating diff: %v\n", err)
			os.Exit(1)
		}
		if diff != "" {
			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				"unformatted files found, see diff below.\n%s",
				diff,
			)
			os.Exit(1)
		}
	case writeFlag:
		config, err := InitConfig()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "failed to load config: %v\n", err)
		}
		if len(args) == 0 {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "no files provided\n")
			return nil
		}
		dir, file := filePathSplit(args[0])
		err = Format(config, dir, file)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error formatting file: %v\n", err)
			os.Exit(1)
		}
	default:
		return cmd.Usage()
	}
	return nil
}

var configFileName = map[string]struct{}{
	".tagcase.yaml": {},
	".tagcase.yml":  {},
}

func filePathSplit(path string) (string, string) {
	dir, file := filepath.Split(path)
	if !strings.HasSuffix(file, ".go") {
		return path, ""
	}
	return dir, file
}

func generateConfigTemplate() string {
	ref := "main"
	if Version != "develop" {
		ref = fmt.Sprintf("refs/tags/%s", Version)
	} else if Revision != "unknown" {
		ref = fmt.Sprintf("refs/heads/%s", Revision)
	}
	return fmt.Sprintf(
		`# yaml-language-server: $schema=https://raw.githubusercontent.com/miyamo2/tagcase/%s/schema.json
# List of tags to check and their expected case.
tags:
  db:
    case: snake_case
  someother:
    # support: snake_case, camelCase, PascalCase, kebab-case, SNAKE_CASE, KEBAB-CASE
    case: camelCase
# List of initialisms to treat as special cases. default: See https://github.com/golang/lint/blob/83fdc39ff7b56453e3793356bcff3070b9b96445/lint.go#L770-L809
#initialism:
#  enable:
#    - BLAH
#  disable:
#    - ID
`,
		ref,
	)
}
