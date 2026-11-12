package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseTranspileArgs(args []string) (string, string, bool) {
	if len(args) == 1 {
		return args[0], "", true
	}
	if len(args) == 3 && args[1] == "-o" && args[2] != "" {
		return args[0], args[2], true
	}
	return "", "", false
}

func defaultPythonOutputPath(inputPath string) string {
	extension := filepath.Ext(inputPath)
	return strings.TrimSuffix(inputPath, extension) + ".py"
}

func runTranspileFile(inputPath, outputPath string) int {
	source, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 65
	}

	program, err := compileSource(string(source))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 65
	}
	generated, err := emitPython(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Python source: %v\n", err)
		return 70
	}

	if outputPath == "" {
		outputPath = defaultPythonOutputPath(inputPath)
	}
	if err := os.WriteFile(outputPath, []byte(generated), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing Python output: %v\n", err)
		return 74
	}
	return 0
}
