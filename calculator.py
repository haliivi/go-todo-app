#!/usr/bin/env python3
"""
Simple safe calculator CLI and REPL.
- Evaluates arithmetic expressions using ast parsing (no eval of arbitrary code).
- Supports +, -, *, /, %, //, **, parentheses, unary +/-, and numbers.
"""
import ast
import operator
import sys

# Allowed operators mapping
BIN_OPS = {
    ast.Add: operator.add,
    ast.Sub: operator.sub,
    ast.Mult: operator.mul,
    ast.Div: operator.truediv,
    ast.FloorDiv: operator.floordiv,
    ast.Mod: operator.mod,
    ast.Pow: operator.pow,
}

UNARY_OPS = {
    ast.UAdd: operator.pos,
    ast.USub: operator.neg,
}


class EvalExpr(ast.NodeVisitor):
    def visit(self, node):
        if isinstance(node, ast.Expression):
            return self.visit(node.body)
        return super().visit(node)

    def visit_BinOp(self, node):
        left = self.visit(node.left)
        right = self.visit(node.right)
        op_type = type(node.op)
        if op_type in BIN_OPS:
            try:
                return BIN_OPS[op_type](left, right)
            except Exception as e:
                raise ValueError(f"Error in binary operation: {e}")
        raise ValueError(f"Unsupported binary operator: {op_type.__name__}")

    def visit_UnaryOp(self, node):
        operand = self.visit(node.operand)
        op_type = type(node.op)
        if op_type in UNARY_OPS:
            return UNARY_OPS[op_type](operand)
        raise ValueError(f"Unsupported unary operator: {op_type.__name__}")

    def visit_Num(self, node):
        return node.n

    def visit_Constant(self, node):
        if isinstance(node.value, (int, float)):
            return node.value
        raise ValueError("Only numeric constants are allowed")

    def visit_Expr(self, node):
        return self.visit(node.value)

    def generic_visit(self, node):
        raise ValueError(f"Unsupported expression: {type(node).__name__}")


def eval_expr(expr: str):
    """Safely evaluate arithmetic expression and return result."""
    try:
        parsed = ast.parse(expr, mode="eval")
    except SyntaxError as e:
        raise ValueError(f"Syntax error: {e}")
    visitor = EvalExpr()
    return visitor.visit(parsed)


def repl():
    print("Simple Calculator REPL. Type 'quit' or 'exit' to leave.")
    while True:
        try:
            s = input('> ').strip()
        except (EOFError, KeyboardInterrupt):
            print()  # newline
            break
        if not s:
            continue
        if s.lower() in ("quit", "exit"):
            break
        try:
            result = eval_expr(s)
            print(result)
        except Exception as e:
            print(f"Error: {e}")


def main(argv=None):
    argv = argv or sys.argv[1:]
    if not argv:
        repl()
        return
    # join args in case expression contains spaces
    expr = " ".join(argv)
    try:
        res = eval_expr(expr)
        print(res)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(2)


if __name__ == '__main__':
    main()
