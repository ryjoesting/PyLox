package main

import (
	"strings"
	"testing"
)

func TestPythonValueRuntimeIncludesLoxSemanticHelpers(t *testing.T) {
	for _, helper := range []string{
		"def __lox_truthy(value):",
		"def __lox_string(value):",
		"def __lox_equal(left, right):",
		"def __lox_binary(operator, left, right):",
		"def __lox_unary(operator, value):",
	} {
		if !strings.Contains(pythonValueRuntime, helper) {
			t.Errorf("Python runtime does not include %q", helper)
		}
	}
}
