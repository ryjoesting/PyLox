package main

const pythonValueRuntime = `import math as __lox_math
from decimal import Decimal as __lox_Decimal


def __lox_truthy(value):
    return value is not None and value is not False


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
