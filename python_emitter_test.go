package main

import (
	"os/exec"
	"strings"
	"testing"
)

func assertPythonSyntax(t *testing.T, source string) {
	t.Helper()
	python, err := exec.LookPath("python")
	if err != nil {
		t.Skip("Python is unavailable for generated-source syntax validation")
	}
	command := exec.Command(python, "-c", "import ast,sys; ast.parse(sys.stdin.read())")
	command.Stdin = strings.NewReader(source)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated Python is invalid: %v\n%s", err, output)
	}
}

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

func TestPythonEmitterStatementsAndClosures(t *testing.T) {
	source, err := transpileSource(`
		fun makeCounter() {
			var count = 0;
			fun increment() {
				count = count + 1;
				return count;
			}
			return increment;
		}
		var counter = makeCounter();
		print counter();
	`)
	if err != nil {
		t.Fatalf("transpile source: %v", err)
	}

	for _, expected := range []string{
		"class __LoxEnv:",
		"class __LoxFunction:",
		"def __lox_function_1(__lox_local, __lox_arguments):",
		"def __lox_function_2(__lox_local, __lox_arguments):",
		`__lox_local.define("count", 0.0)`,
		`__lox_local.assign("count", __lox_binary("+", __lox_local.get("count"), 1.0))`,
		`__lox_globals.define("counter", __lox_call(__lox_globals.get("makeCounter"), []))`,
		"__lox_print(__lox_call(__lox_globals.get(\"counter\"), []))",
	} {
		if !strings.Contains(source, expected) {
			t.Errorf("generated Python does not contain %q", expected)
		}
	}
	assertPythonSyntax(t, source)
}

func TestPythonEmitterPreservesBlockScopeAndControlFlow(t *testing.T) {
	source, err := transpileSource(`
		var i = 0;
		while (i < 2) {
			if (i == 0) print "first"; else print "later";
			i = i + 1;
		}
	`)
	if err != nil {
		t.Fatalf("transpile source: %v", err)
	}
	for _, expected := range []string{
		"while __lox_truthy(__lox_binary(\"<\", __lox_globals.get(\"i\"), 2.0)):",
		"__lox_env_1 = __LoxEnv(__lox_globals)",
		"if __lox_truthy(__lox_binary(\"==\", __lox_env_1.get(\"i\"), 0.0)):",
		"else:",
	} {
		if !strings.Contains(source, expected) {
			t.Errorf("generated Python does not contain %q", expected)
		}
	}
	assertPythonSyntax(t, source)
}

func TestPythonEmitterRejectsClassesUntilClassStage(t *testing.T) {
	if _, err := transpileSource(`class A {}`); err == nil {
		t.Fatal("transpileSource accepted a class before class emission is implemented")
	}
}
