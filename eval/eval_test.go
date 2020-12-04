package eval

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lucasepe/g2d/lexer"
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/parser"
	"github.com/lucasepe/g2d/utils"
)

func assertEvaluated(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	assert := assert.New(t)

	switch expected.(type) {
	case nil:
		if _, ok := actual.(*object.Null); ok {
			assert.True(ok)
		} else {
			assert.Equal(expected, actual)
		}
	case int:
		if i, ok := actual.(*object.Integer); ok {
			assert.Equal(int64(expected.(int)), i.Value)
		} else {
			assert.Equal(expected, actual)
		}
	case error:
		if e, ok := actual.(*object.Integer); ok {
			assert.Equal(expected.(error).Error(), e.Value)
		} else {
			assert.Equal(expected, actual)
		}
	case string:
		if s, ok := actual.(*object.String); ok {
			assert.Equal(expected.(string), s.Value)
		} else {
			assert.Equal(expected, actual)
		}
	default:
		t.Fatalf("unsupported type for expected got=%T", expected)
	}
}

func TestEvalExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"!1", false},
		{"~1", -2},
		{"5 % 2", 1},
		{"1 | 2", 3},
		{"2 ^ 4", 6},
		{"3 & 6", 2},
		{`" " * 4`, "    "},
		{`4 * " "`, "    "},
		{"1 << 2", 4},
		{"4 >> 2", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if expected, ok := tt.expected.(int64); ok {
			testIntegerObject(t, evaluated, expected)
		} else if expected, ok := tt.expected.(bool); ok {
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"true && true", true},
		{"false && true", false},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"false || true", true},
		{"true || false", true},
		{"false || false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(1 <= 2) == true", true},
		{"(1 <= 2) == false", false},
		{"(1 >= 2) == true", false},
		{"(1 >= 2) == false", true},
		{`"a" == "a"`, true},
		{`"a" < "b"`, true},
		{`"abc" == "abc"`, true},
	}

	for _, tt := range tests {
		t.Log(tt.input)
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNullExpression(t *testing.T) {
	evaluated := testEval("null")
	testNullObject(t, evaluated)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 < 2) { 10 } else if (1 == 2) { 20 }", 10},
		{"if (1 > 2) { 10 } else if (1 == 2) { 20 } else { 30 }", 30},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"while (false) { }", nil},
		{"n := 0; while (n < 10) { n := n + 1 }; n", 10},
		{"n := 10; while (n > 0) { n := n - 1 }; n", 0},
		{"n := 0; while (n < 10) { n := n + 1 }", nil},
		{"n := 10; while (n > 0) { n := n - 1 }", nil},
		{"n := 0; while (n < 10) { n = n + 1 }; n", 10},
		{"n := 10; while (n > 0) { n = n - 1 }; n", 0},
		{"n := 0; while (n < 10) { n = n + 1 }", nil},
		{"n := 10; while (n > 0) { n = n - 1 }", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"unknown operator: int + bool",
		},
		{
			"5 + true; 5;",
			"unknown operator: int + bool",
		},
		{
			"-true",
			"unknown operator: -bool",
		},
		{
			"true + false;",
			"unknown operator: bool + bool",
		},
		{
			"5; true + false; 5",
			"unknown operator: bool + bool",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: bool + bool",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: bool + bool",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: str - str",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: fn",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestIndexAssignmentStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"xs := [1, 2, 3]; xs[1] = 4; xs[1];", 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"a := 0; a = 5;", nil},
		{"a := 0; a = 5; a;", 5},
		{"a := 0; a = 5 * 5;", nil},
		{"a := 0; a = 5 * 5; a;", 25},
		{"a := 0; a = 5; b := 0; b = a;", nil},
		{"a := 0; a = 5; b := 0; b = a; b;", 5},
		{"a := 0; a = 5; b := 0; b = a; c := 0; c = a + b + 5;", nil},
		{"a := 0; a = 5; b := 0; b = a; c := 0; c = a + b + 5; c;", 15},
		{"a := 5; b := a; a = 0;", nil},
		{"a := 5; b := a; a = 0; b;", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBindExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a := 5; a;", 5},
		{"a := 5 * 5; a;", 25},
		{"a := 5; b := a; b;", 5},
		{"a := 5; b := a; c := a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"identity := fn(x) { x; }; identity(5);", 5},
		{"identity := fn(x) { return x; }; identity(5);", 5},
		{"double := fn(x) { x * 2; }; double(5);", 10},
		{"add := fn(x, y) { x + y; }; add(5, 5);", 10},
		{"add := fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	newAdder := fn(x) {
	  fn(y) { x + y };
	};

	addTwo := newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`"abc"[0]`,
			"a",
		},
		{
			`"abc"[1]`,
			"b",
		},
		{
			`"abc"[2]`,
			"c",
		},
		{
			`i := 0; "abc"[i];`,
			"a",
		},
		{
			`"abc"[1 + 1];`,
			"c",
		},
		{
			`myString := "abc"; myString[0] + myString[1] + myString[2];`,
			"abc",
		},
		{
			`"abc"[3]`,
			"",
		},
		{
			`"foo"[-1]`,
			"",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, string(str))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, errors.New("TypeError: object of type 'int' has no len()")},
		{`len("one", "two")`, errors.New("TypeError: len() takes exactly 1 argument (2 given)")},
		{`len("∑")`, 1},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, errors.New("TypeError: first() expected argument #1 to be `array` got `int`")},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, errors.New("TypeError: last() expected argument #1 to be `array` got `int`")},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, errors.New("TypeError: push() expected argument #1 to be `array` got `int`")},
		{`print("Hello World")`, nil},
		{`input()`, ""},
		{`pop([])`, errors.New("IndexError: pop from an empty array")},
		{`pop([1])`, 1},
		{`bool(1)`, true},
		{`bool(0)`, false},
		{`bool(true)`, true},
		{`bool(false)`, false},
		{`bool(null)`, false},
		{`bool("")`, false},
		{`bool("foo")`, true},
		{`bool([])`, false},
		{`bool([1, 2, 3])`, true},
		{`bool({})`, false},
		{`bool({"a": 1})`, true},
		{`int(true)`, 1},
		{`int(false)`, 0},
		{`int(1)`, 1},
		{`int("10")`, 10},
		{`str(null)`, "null"},
		{`str(true)`, "true"},
		{`str(false)`, "false"},
		{`str(10)`, "10"},
		{`str("foo")`, "foo"},
		{`str([1, 2, 3])`, "[1, 2, 3]"},
		{`str({"a": 1})`, "{\"a\": 1}"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, expected)
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		case error:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected.Error() {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayDuplication(t *testing.T) {
	input := "[1] * 3"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 1)
	testIntegerObject(t, result.Elements[2], 1)
}

func TestArrayMerging(t *testing.T) {
	input := "[1] + [2]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 2 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 2)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"i := 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"myArray := [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray := [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray := [1, 2, 3]; i := myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `two := "two";
    {
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6 / 2,
        4: 4,
        true: 5,
        false: 6
    }`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashMerging(t *testing.T) {
	input := `{"a": 1} + {"b": 2}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "a"}).HashKey(): 1,
		(&object.String{Value: "b"}).HashKey(): 2,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashSelectorExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}.foo`,
			5,
		},
		{
			`{"foo": 5}.bar`,
			nil,
		},
		{
			`{}.foo`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`key := "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestImportExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`mod := import("../testdata/mod"); mod.A`, 5},
		{`mod := import("../testdata/mod"); mod.Sum(2, 3)`, 5},
		{`mod := import("../testdata/mod"); mod.a`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		assertEvaluated(t, tt.expected, evaluated)
	}
}

func TestImportSearchPaths(t *testing.T) {
	utils.AddPath("../testdata")

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`mod := import("mod"); mod.A`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		assertEvaluated(t, tt.expected, evaluated)
	}
}

func TestExamples(t *testing.T) {
	matches, err := filepath.Glob("./examples/*.monkey")
	if err != nil {
		t.Error(err)
	}

	for _, match := range matches {
		b, err := ioutil.ReadFile(match)
		if err != nil {
			t.Error(err)
		}
		testEval(string(b))
	}
}
