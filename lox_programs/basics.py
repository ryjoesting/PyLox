import math as __lox_math
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
        if __lox_math.isnan(value):
            return "NaN"
        if __lox_math.isinf(value):
            return "+Inf" if value > 0 else "-Inf"
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


class _LoxEnv:
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

    def ancestor(self, distance):
        environment = self
        for _ in range(distance):
            environment = environment.enclosing
        return environment

    def get_at(self, distance, name):
        environment = self.ancestor(distance)
        if name not in environment.values:
            raise RuntimeError("Undefined variable '" + name + "'.")
        return environment.values[name]

    def assign(self, name, value):
        if name in self.values:
            self.values[name] = value
            return value
        if self.enclosing is not None:
            return self.enclosing.assign(name, value)
        raise RuntimeError("Undefined variable '" + name + "'.")

    def assign_at(self, distance, name, value):
        environment = self.ancestor(distance)
        if name not in environment.values:
            raise RuntimeError("Undefined variable '" + name + "'.")
        environment.values[name] = value
        return value


class _LoxFunction:
    def __init__(self, implementation, closure, name, arity, initializer=False):
        self.implementation = implementation
        self.closure = closure
        self.name = name
        self.arity = arity
        self.initializer = initializer

    def call(self, arguments):
        if len(arguments) != self.arity:
            raise RuntimeError("Expected " + str(self.arity) + " arguments but got " + str(len(arguments)) + ".")
        result = self.implementation(_LoxEnv(self.closure), arguments)
        if self.initializer:
            return self.closure.get("this")
        return result

    def __str__(self):
        return "<fn " + self.name + ">"

    def bind(self, instance):
        closure = _LoxEnv(self.closure)
        closure.define("this", instance)
        return _LoxFunction(self.implementation, closure, self.name, self.arity, self.initializer)


class _LoxClass:
    def __init__(self, name, superclass, methods):
        self.name = name
        self.superclass = superclass
        self.methods = methods

    def find_method(self, name):
        if name in self.methods:
            return self.methods[name]
        if self.superclass is not None:
            return self.superclass.find_method(name)
        return None

    def arity(self):
        initializer = self.find_method("init")
        return 0 if initializer is None else initializer.arity

    def call(self, arguments):
        expected = self.arity()
        if len(arguments) != expected:
            raise RuntimeError("Expected " + str(expected) + " arguments but got " + str(len(arguments)) + ".")
        instance = _LoxInstance(self)
        initializer = self.find_method("init")
        if initializer is not None:
            initializer.bind(instance).call(arguments)
        return instance

    def __str__(self):
        return self.name


class _LoxInstance:
    def __init__(self, klass):
        self.klass = klass
        self.fields = {}

    def get(self, name):
        if name in self.fields:
            return self.fields[name]
        method = self.klass.find_method(name)
        if method is not None:
            return method.bind(self)
        raise RuntimeError("Undefined property '" + name + "'.")

    def set(self, name, value):
        self.fields[name] = value
        return value

    def __str__(self):
        return self.klass.name + " instance"


def __lox_call(callee, arguments):
    if isinstance(callee, (_LoxFunction, _LoxClass)):
        return callee.call(arguments)
    raise RuntimeError("Can only call functions and classes.")


def __lox_get(object, name):
    if not isinstance(object, _LoxInstance):
        raise RuntimeError("Only instances have properties.")
    return object.get(name)


def __lox_set(object, name, value):
    if not isinstance(object, _LoxInstance):
        raise RuntimeError("Only instances have fields.")
    return object.set(name, value)


def __lox_super(environment, name):
    superclass = environment.get("super")
    if not isinstance(superclass, _LoxClass):
        raise RuntimeError("Superclass must be a class.")
    instance = environment.get("this")
    method = superclass.find_method(name)
    if method is None:
        raise RuntimeError("Undefined property '" + name + "'.")
    return method.bind(instance)


def __lox_print(value):
    print(__lox_string(value))



__lox_globals = _LoxEnv()
__lox_globals.define("radius", 3.0)
__lox_globals.define("area", __lox_binary("*", __lox_binary("*", 3.14159, __lox_globals.get("radius")), __lox_globals.get("radius")))
__lox_print("Area:")
__lox_print(__lox_globals.get("area"))
__lox_print(__lox_logic("and", lambda: (__lox_binary("<", __lox_globals.get("radius"), 5.0)), lambda: (True)))
if __lox_truthy(__lox_binary(">", __lox_globals.get("area"), 20.0)):
    __lox_env_1 = _LoxEnv(__lox_globals)
    __lox_print("The circle is large.")
else:
    __lox_env_2 = _LoxEnv(__lox_globals)
    __lox_print("The circle is small.")
