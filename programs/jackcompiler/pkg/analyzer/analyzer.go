package analyzer

import (
	"os"
	"strings"
)

// Analyzer handles the top level of analysis
type Analyzer struct {
	inputPath string
	isDir     bool
}

// NewAnalyzer constructs an analyzer from an input file
func NewAnalyzer(inputPath string) *Analyzer {
	// Let's determine if this is a directory or a file
	info, err := os.Stat(inputPath)

	if err != nil {
		panic(err)
		return nil
	}

	isDir := info.IsDir()

	return &Analyzer{inputPath: inputPath, isDir: isDir}
}

// Analyze will analyze the input file(s) and output the xml file(s)
func (a *Analyzer) Analyze() {
	jackFiles := make([]string, 0)

	if a.isDir {
		// Get all jack files in the directory
		files, err := os.ReadDir(a.inputPath)
		if err != nil {
			panic(err)
			return
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jack") {
				jackFiles = append(jackFiles, a.inputPath+"/"+file.Name())
			}
		}
	} else {
		jackFiles = append(jackFiles, a.inputPath)
	}

	// Now we will loop through all files and run analysis
	for _, jackFile := range jackFiles {
		// Create the engine
		engine := NewEngine(jackFile)
		// Process and write
		engine.WriteXML()

	}

}
