package main

import (
	"strings"
	"testing"
)

func TestPythonEmitterExpressions(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{source: `nil`, want: "None"},
		{source: `true`, want: "True"},
		{source: `123`, want: "123.0"},
		{source: `"text"`, want: `"text"`},
		{source: `a = 1`, want: `env.assign("a", 1.0)`},
		{source: `1 + 2 * 3`, want: `__lox_binary("+", 1.0, __lox_binary("*", 2.0, 3.0))`},
		{source: `a and b`, want: `__lox_logic("and", lambda: (env.get("a")), lambda: (env.get("b")))`},
		{source: `a or b`, want: `__lox_logic("or", lambda: (env.get("a")), lambda: (env.get("b")))`},
		{source: `a ? b : c`, want: `__lox_conditional(env.get("a"), lambda: (env.get("b")), lambda: (env.get("c")))`},
		{source: `(a, b)`, want: `(((env.get("a")), (env.get("b")))[1])`},
		{source: `f(1, 2)`, want: `__lox_call(env.get("f"), [1.0, 2.0])`},
		{source: `object.field`, want: `__lox_get(env.get("object"), "field")`},
		{source: `object.field = 2`, want: `__lox_set(env.get("object"), "field", 2.0)`},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			expression, err := parseExpressionSource(test.source)
			if err != nil {
				t.Fatalf("parse expression: %v", err)
			}
			got, err := (&pythonEmitter{}).expression(expression, "env")
			if err != nil {
				t.Fatalf("emit expression: %v", err)
			}
			if got != test.want {
				t.Fatalf("emitted expression:\n got: %s\nwant: %s", got, test.want)
			}
		})
	}
}

func TestPythonEmitterUsesLoxLazyHelpers(t *testing.T) {
	for _, helper := range []string{"__lox_logic", "__lox_conditional"} {
		if !strings.Contains(pythonValueRuntime, "def "+helper+"(") {
			t.Errorf("Python runtime does not include %s", helper)
		}
	}
}
