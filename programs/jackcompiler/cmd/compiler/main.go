package main

import (
	"flag"
	"fmt"
	. "jackcompiler/pkg/analyzer"
	"os"
)

func main() {
	flag.Parse()

	// Make sure we have an input path
	if flag.NArg() != 1 {
		_, err := fmt.Fprintf(os.Stderr, "Usage: jackcompiler <inputPath>\n")
		if err != nil {
			return
		}
		return
	}

	inputPath := flag.Arg(0)
	//inputPath := "testsButcher/Square/Main.jack"

	if inputPath == "" {
		panic("inputPath is required")
		return
	}

	analyzer := NewAnalyzer(inputPath)

	analyzer.Analyze()

}
