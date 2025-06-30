package main

import (
	"fmt"
	"os"

	"github.com/miyamo2/tagcase/internal"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	config, err := internal.InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing configuration: %v\n", err)
		os.Exit(1)
	}
	unitchecker.Main(internal.Analyzer(config))
}
