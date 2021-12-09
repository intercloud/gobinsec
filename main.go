package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	codeVulnerable = 1
	codeError      = 2
)

func main() {
	verbose := flag.Bool("verbose", false, "Print additional information in terminal")
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("ERROR you must pass one binary to analyze")
	}
	path := flag.Args()[0]
	binary, err := NewBinary(path)
	if err != nil {
		fmt.Printf("ERROR analyzing %s: %v\n", path, err)
		os.Exit(codeError)
	}
	binary.Report(*verbose)
	if binary.Vulnerable {
		os.Exit(codeVulnerable)
	}
}
