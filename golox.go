package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runPrompt()
		return
	}

	switch args[0] {
	case "run":
		if len(args) != 2 {
			usage()
			os.Exit(64)
		}
		os.Exit(runFile(args[1]))
	case "transpile":
		inputPath, outputPath, ok := parseTranspileArgs(args[1:])
		if !ok {
			usage()
			os.Exit(64)
		}
		os.Exit(runTranspileFile(inputPath, outputPath))
	}

	if len(args) == 1 {
		os.Exit(runFile(args[0]))
	}
	usage()
	os.Exit(64)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: pylox [run <script> | transpile <input.lox> [-o <output.py>]]")
}
