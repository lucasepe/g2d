package eval

// Package eval implements the evaluator -- a tree-walker implemtnation that
// recursively walks the parsed AST (abstract syntax tree) and evaluates
// the nodes according to their semantic meaning

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	"github.com/lucasepe/g2d/ast"
	"github.com/lucasepe/g2d/builtins"
	"github.com/lucasepe/g2d/lexer"
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/parser"
	"github.com/lucasepe/g2d/token"
	"github.com/lucasepe/g2d/utils"
)

var (
	// TRUE is a cached Boolean object holding the `true` value
	TRUE = &object.Boolean{Value: true}

	// FALSE is a cached Boolean object holding the `false` value
	FALSE = &object.Boolean{Value: false}

	// NULL is a cached Null object
	NULL = &object.Null{}
)

// This program's lexer used for error location in Eval(program)
var lex *lexer.Lexer

func fromNativeBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	if lex == nil {
		return &object.Error{Message: fmt.Sprintf(format, a...)}
	}

	// get the token position from the error node and append the offending line to the error message
	lineNum, _, errorLine := lex.ErrorLine(tok.Position)
	errorPosition := fmt.Sprintf("\n    * line: %d\t%s", lineNum, errorLine)
	return &object.Error{Message: fmt.Sprintf(format, a...) + errorPosition}
}

/*
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
*/

// EvalModule evaluates the named module and returns a *object.Module object
func EvalModule(tok token.Token, name string) object.Object {
	filename := utils.FindModule(name)
	if filename == "" {
		return newError(tok, "ImportError: no module named '%s'", name)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return newError(tok, "IOError: error reading module '%s': %s", name, err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	module := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return newError(tok, "ParseError: %s", p.Errors())
	}

	env := object.NewEnvironment()
	Eval(module, env)

	return env.ExportedHash()
}

// BeginEval (program, env, lexer) object.Object
// REPL and testing modules call this function to init the global lexer pointer for error location
// NB. Eval(node, env) is recursive
func BeginEval(program ast.Node, env *object.Environment, lexer *lexer.Lexer) object.Object {
	// global lexer
	lex = lexer
	// run the evaluator
	return Eval(program, env)
}

// Eval evaluates the node and returns an object
func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return fromNativeBoolean(node.Value)
	case *ast.Null:
		return NULL

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Token, node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Token, node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ImportExpression:
		return evalImportExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(node.Token, env, function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.BindExpression:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			if immutable, ok := value.(object.Immutable); ok {
				env.Set(ident.Value, immutable.Clone())
			} else {
				env.Set(ident.Value, value)
			}

			return NULL
		}
		return newError(node.Token, "expected identifier on left got=%T", node.Left)

	case *ast.AssignmentExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			env.Set(ident.Value, value)
		} else if ie, ok := node.Left.(*ast.IndexExpression); ok {
			obj := Eval(ie.Left, env)
			if isError(obj) {
				return obj
			}

			if array, ok := obj.(*object.Array); ok {
				index := Eval(ie.Index, env)
				if isError(index) {
					return index
				}
				if idx, ok := index.(*object.Integer); ok {
					array.Elements[idx.Value] = value
				} else {
					return newError(node.Token, "cannot index array with %#v", index)
				}
			} else if hash, ok := obj.(*object.Hash); ok {
				key := Eval(ie.Index, env)
				if isError(key) {
					return key
				}
				if hashKey, ok := key.(object.Hashable); ok {
					hashed := hashKey.HashKey()
					hash.Pairs[hashed] = object.HashPair{Key: key, Value: value}
				} else {
					return newError(node.Token, "cannot index hash with %T", key)
				}
			} else {
				return newError(node.Token, "object type %T does not support item assignment", obj)
			}
		} else {
			return newError(node.Token, "expected identifier or index expression got=%T", left)
		}

		return NULL

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(node.Token, left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.SwitchExpression:
		return evalSwitchStatement(node, env)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.Return); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(tok token.Token, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(tok, right)

	default:
		return newError(tok, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return newError(tok, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	switch {

	// " " * 4
	case operator == "*" && left.Type() == object.STRING && right.Type() == object.INTEGER:
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.Integer).Value
		return &object.String{Value: strings.Repeat(leftVal, int(rightVal))}

	// 4 * " "
	case operator == "*" && left.Type() == object.INTEGER && right.Type() == object.STRING:
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: strings.Repeat(rightVal, int(leftVal))}

	case operator == "==":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == 0)
	case operator == "!=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) != 0)
	case operator == "<=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) < 1)
	case operator == ">=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) > -1)
	case operator == "<":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == -1)
	case operator == ">":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == 1)

	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evalBooleanInfixExpression(tok, operator, left, right)
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(tok, operator, left, right)
	case left.Type() == object.FLOAT && right.Type() == object.FLOAT:
		return evalFloatInfixExpression(tok, operator, left, right)
	case left.Type() == object.INTEGER && right.Type() == object.FLOAT:
		return evalIntegerFloatInfixExpression(tok, operator, left, right)
	case left.Type() == object.FLOAT && right.Type() == object.INTEGER:
		return evalFloatIntegerInfixExpression(tok, operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(tok, operator, left, right)

	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return fromNativeBoolean(leftVal && rightVal)
	case "||":
		return fromNativeBoolean(leftVal || rightVal)
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(_ token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalFloatInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerFloatInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := float64(left.(*object.Integer).Value)
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatIntegerInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := float64(right.(*object.Integer).Value)

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			result = Eval(we.Consequence, env)
			if isError(result) {
				return result
			}
		} else {
			break
		}
	}

	if result != nil {
		return result
	}
	return NULL
}

func evalImportExpression(ie *ast.ImportExpression, env *object.Environment) object.Object {
	name := Eval(ie.Name, env)
	if isError(name) {
		return name
	}

	if s, ok := name.(*object.String); ok {
		attrs := EvalModule(ie.Token, s.Value)
		if isError(attrs) {
			return attrs
		}
		return &object.Module{Name: s.Value, Attrs: attrs}
	}
	return newError(ie.Token, "ImportError: invalid import path '%s'", name)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins.Builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Token, "identifier `%s` not found", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(tok token.Token, env *object.Environment, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		fnEnv, err := extendFunctionEnv(tok, fn, args)
		if err != nil {
			return err
		}

		evaluated := Eval(fn.Body, fnEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		if result := fn.Fn(env, args...); result != nil {
			return result
		}
		return NULL

	default:
		return newError(tok, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(tok token.Token, fn *object.Function, args []object.Object) (*object.Environment, *object.Error) {
	env := fn.Env.Clone()

	for paramIdx, param := range fn.Parameters {
		argumentPassed := len(args) > paramIdx

		if !argumentPassed {
			return nil, newError(tok, "argument `%s` to function `%s` is missing", param.Value, fn)
		}

		env.Set(param.Value, args[paramIdx])
	}

	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}

	return obj
}

/*
func evalIndexAssignmentExpression(left, index, value object.Object) object.Object {
	switch {
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
*/
func evalIndexExpression(tok token.Token, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(tok, left, index)
	case left.Type() == object.MODULE:
		return evalModuleIndexExpression(tok, left, index)
	default:
		return newError(tok, "index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(tok token.Token, hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError(tok, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalModuleIndexExpression(tok token.Token, module, index object.Object) object.Object {
	moduleObject := module.(*object.Module)
	return evalHashIndexExpression(tok, moduleObject.Attrs, index)
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return &object.String{Value: ""}
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError(node.Token, "unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalSwitchStatement(se *ast.SwitchExpression, env *object.Environment) object.Object {

	// Get the value.
	obj := Eval(se.Value, env)

	// Try all the choices
	for _, opt := range se.Choices {
		// skipping the default-case, which we'll
		// handle later.
		if opt.Default {
			continue
		}

		// Look at any expression we've got in this case.
		for _, val := range opt.Expr {
			// Get the value of the case
			out := Eval(val, env)

			// Is is a boolean and true?
			if (out.Type() == object.BOOLEAN) && (out.Inspect() == "true") {
				return evalBlockStatement(opt.Block, env)
			}

			// Is it a literal match?
			if (obj.Type() == out.Type()) && (obj.Inspect() == out.Inspect()) {
				// Evaluate the block and return the value
				return evalBlockStatement(opt.Block, env)
			}
		}
	}

	// No match?  Handle default if present
	for _, opt := range se.Choices {

		// skip default
		if opt.Default {

			out := evalBlockStatement(opt.Block, env)
			return out
		}
	}

	return nil
}
