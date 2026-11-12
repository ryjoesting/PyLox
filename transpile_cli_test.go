package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTranspileArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		output      string
		shouldParse bool
	}{
		{name: "default output", args: []string{"source.lox"}, input: "source.lox", shouldParse: true},
		{name: "explicit output", args: []string{"source.lox", "-o", "generated.py"}, input: "source.lox", output: "generated.py", shouldParse: true},
		{name: "missing input", args: []string{}, shouldParse: false},
		{name: "missing output", args: []string{"source.lox", "-o"}, shouldParse: false},
		{name: "unknown option", args: []string{"source.lox", "--output", "generated.py"}, shouldParse: false},
		{name: "extra arguments", args: []string{"source.lox", "generated.py"}, shouldParse: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, output, ok := parseTranspileArgs(test.args)
			if ok != test.shouldParse {
				t.Fatalf("parseTranspileArgs success: got %t, want %t", ok, test.shouldParse)
			}
			if ok && (input != test.input || output != test.output) {
				t.Fatalf("parsed paths: got (%q, %q), want (%q, %q)", input, output, test.input, test.output)
			}
		})
	}
}

func TestDefaultPythonOutputPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: filepath.Join("examples", "hello.lox"), want: filepath.Join("examples", "hello.py")},
		{input: "program", want: "program.py"},
	}
	for _, test := range tests {
		if got := defaultPythonOutputPath(test.input); got != test.want {
			t.Errorf("defaultPythonOutputPath(%q): got %q, want %q", test.input, got, test.want)
		}
	}
}

func TestRunTranspileFileWritesDefaultAndExplicitOutputs(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "sample.lox")
	if err := os.WriteFile(input, []byte(`print "transpiled";`), 0600); err != nil {
		t.Fatal(err)
	}

	defaultOutput := defaultPythonOutputPath(input)
	if status := runTranspileFile(input, ""); status != 0 {
		t.Fatalf("default transpile status: got %d, want 0", status)
	}
	generated, err := os.ReadFile(defaultOutput)
	if err != nil {
		t.Fatalf("read default output: %v", err)
	}
	assertPythonSyntax(t, string(generated))

	explicitOutput := filepath.Join(directory, "custom.py")
	if status := runTranspileFile(input, explicitOutput); status != 0 {
		t.Fatalf("explicit transpile status: got %d, want 0", status)
	}
	generated, err = os.ReadFile(explicitOutput)
	if err != nil {
		t.Fatalf("read explicit output: %v", err)
	}
	assertPythonSyntax(t, string(generated))
}

func TestRunTranspileFileErrors(t *testing.T) {
	if status := runTranspileFile("missing-transpile-input.lox", ""); status != 65 {
		t.Fatalf("missing input status: got %d, want 65", status)
	}

	directory := t.TempDir()
	invalidInput := filepath.Join(directory, "invalid.lox")
	if err := os.WriteFile(invalidInput, []byte(`var broken = ;`), 0600); err != nil {
		t.Fatal(err)
	}
	if status := runTranspileFile(invalidInput, ""); status != 65 {
		t.Fatalf("invalid source status: got %d, want 65", status)
	}

	validInput := filepath.Join(directory, "valid.lox")
	if err := os.WriteFile(validInput, []byte(`print 1;`), 0600); err != nil {
		t.Fatal(err)
	}
	if status := runTranspileFile(validInput, filepath.Join(directory, "missing", "out.py")); status != 74 {
		t.Fatalf("output error status: got %d, want 74", status)
	}
}
