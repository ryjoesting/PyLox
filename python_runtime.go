package main

const pythonValueRuntime = `import math as __lox_math
from decimal import Decimal as __lox_Decimal


def __lox_truthy(value):
    return value is not None and value is not False


def __lox_logic(operator, left, right):
    value = left()
    if operator == "or" and __lox_truthy(value):
        return value
    if operator == "and" and not __lox_truthy(value):
        return value
    return right()


def __lox_conditional(condition, then_branch, else_branch):
    if __lox_truthy(condition):
        return then_branch()
    return else_branch()


def __lox_string(value):
    if value is None:
        return "nil"
    if value is True:
        return "true"
    if value is False:
        return "false"
    if type(value) is float:
        text = format(__lox_Decimal(repr(value)), "f")
        if "." in text:
            text = text.rstrip("0").rstrip(".")
        return text
    return str(value)


def __lox_equal(left, right):
    if left is None or right is None:
        return left is None and right is None
    if type(left) is not type(right):
        return False
    if type(left) in (bool, float, str):
        return left == right
    return left is right


def __lox_number(value, operator):
    if type(value) is not float:
        raise RuntimeError("Operand must be a number." if operator in ("-", "!") else "Operands must be numbers.")
    return value


def __lox_binary(operator, left, right):
    if operator == "+":
        if type(left) is float and type(right) is float:
            return left + right
        if type(left) is str and type(right) is str:
            return left + right
        raise RuntimeError("Operands must be two numbers or two strings.")
    if operator in ("-", "*", "/", "<", "<=", ">", ">="):
        left = __lox_number(left, operator)
        right = __lox_number(right, operator)
        if operator == "-":
            return left - right
        if operator == "*":
            return left * right
        if operator == "/":
            if right == 0.0:
                if left == 0.0:
                    return float("nan")
                sign = __lox_math.copysign(1.0, left) * __lox_math.copysign(1.0, right)
                return __lox_math.copysign(float("inf"), sign)
            return left / right
        if operator == "<":
            return left < right
        if operator == "<=":
            return left <= right
        if operator == ">":
            return left > right
        return left >= right
    if operator == "==":
        return __lox_equal(left, right)
    if operator == "!=":
        return not __lox_equal(left, right)
    raise RuntimeError("Unsupported operator: " + operator)


def __lox_unary(operator, value):
    if operator == "!":
        return not __lox_truthy(value)
    if operator == "-":
        return -__lox_number(value, operator)
    raise RuntimeError("Unsupported unary operator: " + operator)
`

const pythonEnvironmentRuntime = `

class __LoxEnv:
    def __init__(self, enclosing=None):
        self.enclosing = enclosing
        self.values = {}

    def define(self, name, value):
        self.values[name] = value
        return value

    def get(self, name):
        if name in self.values:
            return self.values[name]
        if self.enclosing is not None:
            return self.enclosing.get(name)
        raise RuntimeError("Undefined variable '" + name + "'.")

    def assign(self, name, value):
        if name in self.values:
            self.values[name] = value
            return value
        if self.enclosing is not None:
            return self.enclosing.assign(name, value)
        raise RuntimeError("Undefined variable '" + name + "'.")


class __LoxFunction:
    def __init__(self, implementation, closure, name, arity):
        self.implementation = implementation
        self.closure = closure
        self.name = name
        self.arity = arity

    def call(self, arguments):
        if len(arguments) != self.arity:
            raise RuntimeError("Expected " + str(self.arity) + " arguments but got " + str(len(arguments)) + ".")
        return self.implementation(__LoxEnv(self.closure), arguments)

    def __str__(self):
        return "<fn " + self.name + ">"


def __lox_call(callee, arguments):
    if not isinstance(callee, __LoxFunction):
        raise RuntimeError("Can only call functions and classes.")
    return callee.call(arguments)


def __lox_print(value):
    print(__lox_string(value))
`
