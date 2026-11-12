package main

type compiledProgram struct {
	statements []Stmt
	locals     map[Expr]int
	globals    map[Expr]struct{}
}

func compileSource(source string) (*compiledProgram, error) {
	statements, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	resolver := NewResolver(nil)
	if err := resolver.Resolve(statements); err != nil {
		return nil, err
	}

	return &compiledProgram{
		statements: statements,
		locals:     resolver.locals,
		globals:    resolver.globals,
	}, nil
}
