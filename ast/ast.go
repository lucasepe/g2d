package ast

// Packge ast implement the Abstract Syntax Tree that represents the parsed
// source code before being passed on to the interpreter for evaluation.

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lucasepe/g2d/token"
)

// Node defines an interface for all nodes in the AST.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement defines the interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression defines the interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node. All programs consist of a slice of Statement(s)
type Program struct {
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a stringified version of the AST for debugging
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ReturnStatement represenets the `return` statement node
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a stringified version of the AST for debugging
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression statement and holds an
// expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a stringified version of the AST for debugging
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block statement and holds one or more other
// statements
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns a stringified version of the AST for debugging
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier represents an identiifer and holds the name of the identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns a stringified version of the AST for debugging
func (i *Identifier) String() string { return i.Value }

// Null represents a null value
type Null struct {
	Token token.Token
}

func (n *Null) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (n *Null) TokenLiteral() string { return n.Token.Literal }

// String returns a stringified version of the AST for debugging
func (n *Null) String() string { return n.Token.Literal }

// Boolean represents a boolean value and holds the underlying boolean value
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns a stringified version of the AST for debugging
func (b *Boolean) String() string { return b.Token.Literal }

// IntegerLiteral represents a literal integare and holds an integer value
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the AST for debugging
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// FloatLiteral holds a floating-point number
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a stringified version of the AST for debugging
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

// StringLiteral represents a literal string and holds a string value
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String returns a stringified version of the AST for debugging
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// PrefixExpression represents a prefix expression and holds the operator
// as well as the right-hand side expression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a stringified version of the AST for debugging
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an infix expression and holds the left-hand
// expression, operator and right-hand expression
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfExpression represents an `if` expression and holds the condition,
// consequence and alternative expressions
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// WhileExpression represents an `while` expression and holds the condition,
// and consequence expression
type WhileExpression struct {
	Token       token.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }

// String returns a stringified version of the AST for debugging
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())

	return out.String()
}

// FunctionLiteral represents a literal functions and holds the function's
// formal parameters and boy of the function as a block statement
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a stringified version of the AST for debugging
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fmt.Sprintf("%s %s", fl.TokenLiteral(), fl.Name))
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents a call expression and holds the function to be
// called as well as the arguments to be passed to that function
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// ArrayLiteral represents the array literal and holds a list of expressions
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String returns a stringified version of the AST for debugging
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// BindExpression represents a binding expression of the form:
// x := 1
type BindExpression struct {
	Token token.Token // The := token
	Left  Expression
	Value Expression
}

func (be *BindExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (be *BindExpression) TokenLiteral() string { return be.Token.Literal }

// String returns a stringified version of the AST for debugging
func (be *BindExpression) String() string {
	var out bytes.Buffer

	out.WriteString(be.Left.String())
	out.WriteString(be.TokenLiteral())
	out.WriteString(be.Value.String())

	return out.String()
}

// AssignmentExpression represents an assignment expression of the form:
// x = 1 or xs[1] = 2
type AssignmentExpression struct {
	Token token.Token // The = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(ae.TokenLiteral())
	out.WriteString(ae.Value.String())

	return out.String()
}

// IndexExpression represents an index operator expression, e.g: xs[2]
// and holds the left expression and index expression
type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// CaseExpression handles the case within a switch statement
type CaseExpression struct {
	// Token is the actual token
	Token token.Token

	// Default branch?
	Default bool

	// The thing we match
	Expr []Expression

	// The code to execute if there is a match
	Block *BlockStatement
}

func (ce *CaseExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (ce *CaseExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns this object as a string.
func (ce *CaseExpression) String() string {
	var out bytes.Buffer

	if ce.Default {
		out.WriteString("default ")
	} else {
		out.WriteString("case ")

		tmp := []string{}
		for _, exp := range ce.Expr {
			tmp = append(tmp, exp.String())
		}
		out.WriteString(strings.Join(tmp, ","))
	}
	out.WriteString(ce.Block.String())
	return out.String()
}

// SwitchExpression handles a switch statement
type SwitchExpression struct {
	// Token is the actual token
	Token token.Token

	// Value is the thing that is evaluated to determine
	// which block should be executed.
	Value Expression

	// The branches we handle
	Choices []*CaseExpression
}

func (se *SwitchExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (se *SwitchExpression) TokenLiteral() string { return se.Token.Literal }

// String returns this object as a string.
func (se *SwitchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("\nswitch (")
	out.WriteString(se.Value.String())
	out.WriteString(")\n{\n")

	for _, tmp := range se.Choices {
		if tmp != nil {
			out.WriteString(tmp.String())
		}
	}
	out.WriteString("}\n")

	return out.String()
}
