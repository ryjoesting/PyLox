package main

import "testing"

func TestCompileSourceReturnsResolvedProgram(t *testing.T) {
	program, err := compileSource(`fun outer() { var value = 1; fun inner() { print value; } }`)
	if err != nil {
		t.Fatalf("compile source: %v", err)
	}
	if len(program.statements) != 1 {
		t.Fatalf("compiled statements: got %d, want 1", len(program.statements))
	}
	if len(program.locals) == 0 {
		t.Fatal("compiled program has no local bindings")
	}
}

func TestCompileSourceReturnsFrontEndErrors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{name: "scan error", source: `@`},
		{name: "parse error", source: `var value = ;`},
		{name: "resolve error", source: `return 1;`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := compileSource(test.source)
			if err == nil {
				t.Fatal("compileSource succeeded, want a front-end error")
			}
			if program != nil {
				t.Fatalf("compileSource returned a program after an error: %#v", program)
			}
		})
	}
}
