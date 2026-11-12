# PyLox

PyLox transpiles Lox programs into standalone Python source files. It scans,
parses, and resolves a Lox program before generating Python 3.10+ code with the
helpers needed to preserve Lox semantics.

## Install

From the repository root, install the `pylox` command:

```powershell
go install .
```

On Windows, Go installs `pylox.exe` under `GOBIN` or, by default,
`%USERPROFILE%\go\bin`. Make sure that directory is on `PATH`, then open a
new terminal.

## Use

Transpile a Lox program:

```powershell
pylox transpile "lox_programs\basics.lox"
```

The generated Python file is written beside the Lox input with the same base
name. Specify another output path with `-o`:

```powershell
pylox transpile "lox_programs\classes.lox" -o "output\classes.py"
```

Run a Lox program with the Go interpreter:

```powershell
pylox run "lox_programs\basics.lox"
```

Running `pylox` without arguments starts the interpreter REPL. Generated
Python files include their Lox compatibility helpers and do not require a
PyLox runtime package. Python files generated under `lox_programs/` are
ignored by Git.

## Test

Requirements: Go 1.27.1+ and Python 3.10+.

```powershell
go test ./...
go vet ./...
```

The tests syntax-check and execute generated Python, then compare its output
with the Go interpreter. They cover each `.lox` example in `lox_programs/` and
focused cases for expressions, lexical scope, closures, functions, classes,
inheritance, and CLI output.

## Examples

The [`lox_programs/`](./lox_programs/) directory contains transpilation
fixtures:

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
- `resolver.go` records lexical bindings used by the transpiler.
- `compiler.go` provides the shared scan/parse/resolve pipeline.
- `python_emitter.go` and `python_runtime.go` generate standalone Python.
- `transpile_cli.go` handles transpilation and output-file selection.
- `interpreter.go`, `function.go`, and `class.go` provide the `pylox run`
  command and a behavioral reference for tests.
