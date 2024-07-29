package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zer0yu/xrecon/pkg/runner"
)

func main() {
	options := &runner.Options{}

	flag.StringVar(&options.URL, "url", "", "Single URL to scan")
	flag.StringVar(&options.URLFile, "file", "", "File containing URLs to scan")
	flag.StringVar(&options.OutputFormat, "output", "terminal", "Output format: terminal, csv or txt")
	flag.StringVar(&options.OutputFile, "o", "", "Output file name (without extension)")
	flag.Parse()

	if options.URL == "" && options.URLFile == "" {
		fmt.Println("Please provide either a URL (-url) or a file with URLs (-file)")
		os.Exit(1)
	}

	runner, err := runner.NewRunner(options)
	if err != nil {
		fmt.Printf("Error creating runner: %v\n", err)
		os.Exit(1)
	}

	err = runner.Run()
	if err != nil {
		fmt.Printf("Error running xrecon: %v\n", err)
		os.Exit(1)
	}
}
