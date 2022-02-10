package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/intercloud/gobinsec/gobinsec"
)

const (
	CodeVulnerable = 1
	CodeError      = 2
)

var Version = "NONE"

func main() {
	version := flag.Bool("version", false, "Print gobinsec version")
	verbose := flag.Bool("verbose", false, "Print additional information in terminal")
	strict := flag.Bool("strict", false, "Vulnerabilities without version are exposed")
	config := flag.String("config", "", "Configuration file")
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if len(flag.Args()) < 1 {
		println("ERROR you must pass binary/ies to analyze on command line")
		os.Exit(CodeError)
	}
	if err := gobinsec.LoadConfig(*config, *strict); err != nil {
		println(fmt.Sprintf("ERROR %v", err))
		os.Exit(CodeError)
	}
	issue := false
	for _, path := range flag.Args() {
		binary, err := gobinsec.NewBinary(path)
		if err != nil {
			gobinsec.ColorRed.Print("ERROR")
			fmt.Printf("analyzing %s: %v\n", path, err)
			issue = true
		} else {
			binary.Report(*verbose)
			if binary.Vulnerable {
				issue = true
			}
		}
	}
	if issue {
		os.Exit(CodeVulnerable)
	}
}
