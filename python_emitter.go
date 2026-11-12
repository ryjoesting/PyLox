package main

import (
	"fmt"
	"strconv"
	"strings"
)

type pythonEmitter struct {
	builder         strings.Builder
	nextFunction    int
	nextEnvironment int
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

func transpileSource(source string) (string, error) {
	program, err := compileSource(source)
	if err != nil {
		return "", err
	}
	return emitPython(program)
}

func emitPython(program *compiledProgram) (string, error) {
	emitter := &pythonEmitter{}
	emitter.builder.WriteString(pythonValueRuntime)
	emitter.builder.WriteString(pythonEnvironmentRuntime)
	emitter.builder.WriteString("\n\n__lox_globals = _LoxEnv()\n")
	if err := emitter.emitStatements(program.statements, "__lox_globals", 0); err != nil {
		return "", err
	}
	return emitter.builder.String(), nil
}

func (emitter *pythonEmitter) emitStatements(statements []Stmt, environment string, indent int) error {
	if len(statements) == 0 {
		emitter.writeLine(indent, "pass")
		return nil
	}
	for _, statement := range statements {
		if err := emitter.emitStatement(statement, environment, indent); err != nil {
			return err
		}
	}
	return nil
}

func (emitter *pythonEmitter) emitStatement(statement Stmt, environment string, indent int) error {
	switch statement := statement.(type) {
	case *BlockStmt:
		emitter.nextEnvironment++
		child := fmt.Sprintf("__lox_env_%d", emitter.nextEnvironment)
		emitter.writeLine(indent, child+" = _LoxEnv("+environment+")")
		return emitter.emitStatements(statement.Statements, child, indent)
	case *ClassStmt:
		return emitter.emitClass(statement, environment, indent)
	case *ExpressionStmt:
		expression, err := emitter.expression(statement.Expression, environment)
		if err != nil {
			return err
		}
		emitter.writeLine(indent, "_ = "+expression)
	case *FunctionStmt:
		name, err := emitter.emitFunction(statement, indent)
		if err != nil {
			return err
		}
		emitter.writeLine(indent, fmt.Sprintf("%s.define(%s, _LoxFunction(%s, %s, %s, %d))",
			environment, strconv.Quote(statement.Name.Lexeme), name, environment,
			strconv.Quote(statement.Name.Lexeme), len(statement.Parameters)))
	case *IfStmt:
		condition, err := emitter.expression(statement.Condition, environment)
		if err != nil {
			return err
		}
		emitter.writeLine(indent, "if __lox_truthy("+condition+"):")
		if err := emitter.emitStatement(statement.ThenBranch, environment, indent+1); err != nil {
			return err
		}
		if statement.ElseBranch != nil {
			emitter.writeLine(indent, "else:")
			if err := emitter.emitStatement(statement.ElseBranch, environment, indent+1); err != nil {
				return err
			}
		}
	case *PrintStmt:
		expression, err := emitter.expression(statement.Expression, environment)
		if err != nil {
			return err
		}
		emitter.writeLine(indent, "__lox_print("+expression+")")
	case *ReturnStmt:
		if statement.Value == nil {
			emitter.writeLine(indent, "return None")
		} else {
			expression, err := emitter.expression(statement.Value, environment)
			if err != nil {
				return err
			}
			emitter.writeLine(indent, "return "+expression)
		}
	case *VarStmt:
		value := "None"
		if statement.Initializer != nil {
			var err error
			value, err = emitter.expression(statement.Initializer, environment)
			if err != nil {
				return err
			}
		}
		emitter.writeLine(indent, fmt.Sprintf("%s.define(%s, %s)", environment,
			strconv.Quote(statement.Name.Lexeme), value))
	case *WhileStmt:
		condition, err := emitter.expression(statement.Condition, environment)
		if err != nil {
			return err
		}
		emitter.writeLine(indent, "while __lox_truthy("+condition+"):")
		if err := emitter.emitStatement(statement.Body, environment, indent+1); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported Lox statement %T", statement)
	}
	return nil
}

func (emitter *pythonEmitter) emitFunction(function *FunctionStmt, indent int) (string, error) {
	emitter.nextFunction++
	name := fmt.Sprintf("__lox_function_%d", emitter.nextFunction)
	local := "__lox_local"
	arguments := "__lox_arguments"
	emitter.writeLine(indent, fmt.Sprintf("def %s(%s, %s):", name, local, arguments))
	bodyIndent := indent + 1
	for index, parameter := range function.Parameters {
		emitter.writeLine(bodyIndent, fmt.Sprintf("%s.define(%s, %s[%d])",
			local, strconv.Quote(parameter.Lexeme), arguments, index))
	}
	if err := emitter.emitStatements(function.Body, local, bodyIndent); err != nil {
		return "", err
	}
	return name, nil
}

func (emitter *pythonEmitter) emitClass(statement *ClassStmt, environment string, indent int) error {
	className := statement.Name.Lexeme
	quotedName := strconv.Quote(className)
	emitter.writeLine(indent, fmt.Sprintf("%s.define(%s, None)", environment, quotedName))

	superclass := "None"
	methodClosure := environment
	if statement.Superclass != nil {
		superclassExpression, err := emitter.expression(statement.Superclass, environment)
		if err != nil {
			return err
		}
		emitter.nextEnvironment++
		superclass = fmt.Sprintf("__lox_superclass_%d", emitter.nextEnvironment)
		emitter.writeLine(indent, fmt.Sprintf("%s = %s", superclass, superclassExpression))
		emitter.writeLine(indent, fmt.Sprintf("if not isinstance(%s, _LoxClass):", superclass))
		emitter.writeLine(indent+1, `raise RuntimeError("Superclass must be a class.")`)
		emitter.nextEnvironment++
		methodClosure = fmt.Sprintf("__lox_env_%d", emitter.nextEnvironment)
		emitter.writeLine(indent, fmt.Sprintf("%s = _LoxEnv(%s)", methodClosure, environment))
		emitter.writeLine(indent, fmt.Sprintf("%s.define(\"super\", %s)", methodClosure, superclass))
	}

	methods := make([]string, 0, len(statement.Methods))
	for _, method := range statement.Methods {
		implementation, err := emitter.emitFunction(method, indent)
		if err != nil {
			return err
		}
		initializer := "False"
		if method.Name.Lexeme == "init" {
			initializer = "True"
		}
		methods = append(methods, fmt.Sprintf("%s: _LoxFunction(%s, %s, %s, %d, %s)",
			strconv.Quote(method.Name.Lexeme), implementation, methodClosure,
			strconv.Quote(method.Name.Lexeme), len(method.Parameters), initializer))
	}
	class := fmt.Sprintf("_LoxClass(%s, %s, {%s})", quotedName, superclass, strings.Join(methods, ", "))
	emitter.writeLine(indent, fmt.Sprintf("%s.assign(%s, %s)", environment, quotedName, class))
	return nil
}

func (emitter *pythonEmitter) writeLine(indent int, line string) {
	emitter.builder.WriteString(strings.Repeat("    ", indent))
	emitter.builder.WriteString(line)
	emitter.builder.WriteByte('\n')
}
