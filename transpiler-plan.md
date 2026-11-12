# PyLox Transpiler Plan

## Settled contract

- `golox transpile <input.lox> [-o <output.py>]` writes Python source. Without
  `-o`, the output is beside the input with the same basename and a `.py`
  extension.
- `golox run <input.lox>` runs the Lox interpreter. The existing
  `golox <input.lox>` invocation remains an interpreter alias, and running
  `golox` without arguments continues to open the REPL.
- Generated output is self-contained, syntactically valid Python 3.10+ source;
  it must not rely on an installed PyLox runtime module.
- The deliverable is the Python source file. Running generated files is not a
  project requirement or acceptance test.
- Scope is the language accepted by this repository's current scanner, parser,
  and resolver, including its expressions, lexical scopes, functions,
  closures, classes, inheritance, and desugared `for` loops.

## Stages

Each stage is intended to be a separate commit.

### 1. Isolate the compilation pipeline (complete)

- Separate scanning, parsing, and resolution from execution so the future
  emitter can consume a resolved program without running it.
- Establish and document the CLI, output, and Python-version contract above.
- Preserve existing interpreter and REPL behavior.
- **Complete when:** the compiler pipeline returns a resolved program or a
  front-end error, existing interpreter behavior remains intact, and Go tests
  pass.

### 2. Add Python source support helpers (complete)

- Add only the self-contained Python helper definitions needed to express Lox
  behavior that Python does not directly provide, such as Lox truthiness,
  `nil`, equality, and value formatting.
- Keep generated files independent of a separate runtime package.
- **Complete when:** focused tests verify emitted helper source and generated
  source remains valid Python syntax; generated programs do not need to be run.

### 3. Transpile expressions (complete)

- Implement an AST-based emitter for every expression node.
- Preserve precedence, associativity, evaluation order, short-circuiting, and
  expression values in the generated Python.
- **Complete when:** focused golden tests cover each expression form and
  generated Python passes syntax validation.

### 4. Transpile statements, scopes, and functions (complete)

- Emit declarations, blocks, conditionals, loops, print statements, functions,
  and returns.
- Use resolver binding information and explicit Lox environments to preserve
  shadowing, recursion, and closure structure without exposing Lox identifiers
  as Python variable names.
- **Complete when:** golden tests cover nested scopes, closures, recursion, and
  control flow, and output passes syntax validation.

### 5. Transpile classes and inheritance (complete)

- Emit classes, initializers, methods, `this`, `super`, inheritance, and
  dynamic fields.
- Preserve Lox method lookup and binding in the generated source.
- **Complete when:** tests cover the repository's class and inheritance
  constructs and output passes syntax validation.

### 6. Expose end-to-end transpilation (complete)

- Add the `transpile` CLI and `-o` output option, retaining the agreed run and
  REPL behavior.
- Report source and output-file errors clearly and document the commands in
  the README.
- **Complete when:** CLI tests verify default and explicit output paths,
  representative outputs pass Python syntax validation, and `go test ./...`
  passes.
