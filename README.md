# PyLox

PyLox transpiles Lox programs into standalone, self-contained Python source files. This project is an extension of my [GoLox Interpreter](https://github.com/ryjoesting/GoLox) and includes the interpreter. The user can choose to execute a Lox source file, transpile it to Python, or open the Lox read-execute-print-loop (REPL).

## Install

To avoid cumbersome commands, I created a shortcut command to run the project's different modes (transpile, run, or REPL for any Lox input).

From the repository root, install the `pylox` command:

```powershell
go install .
```

On Windows, Go installs `pylox.exe` under the `GOBIN` or, by default,
`%USERPROFILE%\go\bin`. Make sure that directory is on the `PATH`, then open a
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
Python files include their Lox compatibility helpers. This may make output Python programs difficult to read but it helps keep the source file self-contained. This behavior is expected in most transpilers (fundamentally, differences in the language and syntax must be resolved). 

## Test

The automated tests for this project use the built in Go test suite. Needs Go 1.27.1+ and Python 3.10+.

```powershell
go test ./...
go vet ./...
```

The tests syntax-check and execute generated Python, then compare its output
with the Go interpreter. They cover each `.lox` example in `lox_programs/` and
focused cases for expressions, lexical scope, closures, functions, classes,
inheritance, and CLI output.

## How the Transpiler Works

At a high-level, this transpiler works by (fill in)

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
