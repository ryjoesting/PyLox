# PyLox

PyLox transpiles Lox programs into standalone Python source files. The Go
front end scans, parses, and resolves a Lox program before emitting Python
3.10+ code with the helpers needed to preserve Lox semantics.

## Requirements

- Go 1.27.1 or later
- Python 3.10 or later for running generated files and the transpiler tests
- No third-party Go dependencies

## Transpile a program

Run these commands from the repository root:

```powershell
go run . transpile "lox_programs\basics.lox"
python "lox_programs\basics.py"
```

By default, the `.py` file is written beside the input with the same base
name. Choose a different path with `-o`:

```powershell
go run . transpile "lox_programs\classes.lox" -o "output\classes.py"
python "output\classes.py"
```

Generated files include their compatibility runtime and do not depend on an
installed PyLox package. Python output files are ignored by Git under
`lox_programs/` to keep generated artifacts out of the source tree.

The CLI also retains the Go interpreter for comparison:

```powershell
go run . run "lox_programs\basics.lox"
```

The legacy `go run . <file.lox>` invocation still runs the interpreter, and
`go run .` starts its REPL.

## Test

Run the Go test suite and static analysis from the repository root:

```powershell
go test ./...
go vet ./...
```

The transpiler tests use Python to syntax-check and execute generated output,
then compare its output with the Go interpreter. They cover every `.lox`
example in `lox_programs/` plus focused cases for expression semantics,
lexical scopes, closures, functions, classes, inheritance, and CLI output.

## Examples

The [`lox_programs/`](./lox_programs/) directory contains small programs used
by the test suite:

- `basics.lox` — variables, arithmetic, logic, and conditionals
- `expressions.lox` — truthiness, short-circuiting, conditionals, and comma
  expressions
- `functions.lox` — recursion and mutable closures
- `scopes.lox` — lexical binding and shadowing
- `loops.lox` — loops and function calls
- `closures.lox` — closures capturing and updating local state
- `classes.lox` — fields, methods, inheritance, and `super`

## Project structure

- `scanner.go`, `parser.go`, `ast.go`, and `statement.go` implement the Lox
  front end.
- `resolver.go` records lexical bindings used by both the interpreter and
  transpiler.
- `compiler.go` provides the shared scan/parse/resolve pipeline.
- `python_emitter.go` and `python_runtime.go` generate standalone Python.
- `transpile_cli.go` handles transpilation and output-file selection.
- `interpreter.go`, `function.go`, and `class.go` retain the Go interpreter as
  a compatibility command and a behavioral reference for tests.
