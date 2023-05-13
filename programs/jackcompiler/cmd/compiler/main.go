package main

import (
	. "jackcompiler/pkg/analyzer"
)

func main() {
	analyzer := NewAnalyzer("testsButcher/Square/Square.jack")

	analyzer.Analyze()

}
