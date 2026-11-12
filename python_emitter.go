package main

import (
	"fmt"
	"strconv"
	"strings"
)

type pythonEmitter struct {
	nextFunction int
}

func (emitter *pythonEmitter) expression(expression Expr, environment string) (string, error) {
	switch expression := expression.(type) {
	case *Assign:
		value, err := emitter.expression(expression.Value, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.assign(%s, %s)", environment, strconv.Quote(expression.Name.Lexeme), value), nil
	case *Binary:
		left, err := emitter.expression(expression.Left, environment)
		if err != nil {
			return "", err
		}
		right, err := emitter.expression(expression.Right, environment)
		if err != nil {
			return "", err
		}
		if expression.Operator.Type == COMMA {
			return fmt.Sprintf("((%s), (%s))[1]", left, right), nil
		}
		return fmt.Sprintf("__lox_binary(%s, %s, %s)", strconv.Quote(expression.Operator.Lexeme), left, right), nil
	case *Call:
		callee, err := emitter.expression(expression.Callee, environment)
		if err != nil {
			return "", err
		}
		arguments := make([]string, 0, len(expression.Arguments))
		for _, argument := range expression.Arguments {
			value, err := emitter.expression(argument, environment)
			if err != nil {
				return "", err
			}
			arguments = append(arguments, value)
		}
		return fmt.Sprintf("__lox_call(%s, [%s])", callee, strings.Join(arguments, ", ")), nil
	case *Conditional:
		condition, err := emitter.expression(expression.Condition, environment)
		if err != nil {
			return "", err
		}
		thenBranch, err := emitter.expression(expression.Then, environment)
		if err != nil {
			return "", err
		}
		elseBranch, err := emitter.expression(expression.Else, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("__lox_conditional(%s, lambda: (%s), lambda: (%s))", condition, thenBranch, elseBranch), nil
	case *Get:
		object, err := emitter.expression(expression.Object, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("__lox_get(%s, %s)", object, strconv.Quote(expression.Name.Lexeme)), nil
	case *Grouping:
		inner, err := emitter.expression(expression.Expression, environment)
		if err != nil {
			return "", err
		}
		return "(" + inner + ")", nil
	case *Literal:
		switch value := expression.Value.(type) {
		case nil:
			return "None", nil
		case bool:
			if value {
				return "True", nil
			}
			return "False", nil
		case float64:
			text := strconv.FormatFloat(value, 'g', -1, 64)
			if !strings.ContainsAny(text, ".eE") {
				text += ".0"
			}
			return text, nil
		case string:
			return strconv.Quote(value), nil
		default:
			return "", fmt.Errorf("unsupported Lox literal type %T", expression.Value)
		}
	case *Logical:
		left, err := emitter.expression(expression.Left, environment)
		if err != nil {
			return "", err
		}
		right, err := emitter.expression(expression.Right, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("__lox_logic(%s, lambda: (%s), lambda: (%s))",
			strconv.Quote(expression.Operator.Lexeme), left, right), nil
	case *Set:
		object, err := emitter.expression(expression.Object, environment)
		if err != nil {
			return "", err
		}
		value, err := emitter.expression(expression.Value, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("__lox_set(%s, %s, %s)", object, strconv.Quote(expression.Name.Lexeme), value), nil
	case *Super:
		return fmt.Sprintf("__lox_super(%s, %s)", environment, strconv.Quote(expression.Method.Lexeme)), nil
	case *This:
		return environment + `.get("this")`, nil
	case *Unary:
		right, err := emitter.expression(expression.Right, environment)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("__lox_unary(%s, %s)", strconv.Quote(expression.Operator.Lexeme), right), nil
	case *Variable:
		return fmt.Sprintf("%s.get(%s)", environment, strconv.Quote(expression.Name.Lexeme)), nil
	default:
		return "", fmt.Errorf("unsupported Lox expression %T", expression)
	}
}
