package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/intercloud/gobinsec/gobinsec"
)

const (
	codeVulnerable = 1
	codeError      = 2
)

func main() {
	verbose := flag.Bool("verbose", false, "Print additional information in terminal")
	config := flag.String("config", "", "Configuration file")
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("ERROR you must pass one binary to analyze")
		os.Exit(codeError)
	}
	if err := gobinsec.LoadConfig(*config); err != nil {
		println(fmt.Sprintf("ERROR loading configuration: %v", err))
		os.Exit(codeError)
	}
	path := flag.Args()[0]
	binary, err := gobinsec.NewBinary(path)
	if err != nil {
		println(fmt.Sprintf("ERROR analyzing %s: %v", path, err))
		os.Exit(codeError)
	}
	binary.Report(*verbose)
	if binary.Vulnerable {
		os.Exit(codeVulnerable)
	}
}
