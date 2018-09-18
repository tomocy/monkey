package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/tomocy/monkey/ast"
)

const (
	Integer         = "Integer"
	Boolean         = "Boolean"
	String          = "String"
	Array           = "Array"
	Null            = "Null"
	Return          = "Return"
	Error           = "Error"
	Function        = "Function"
	BuiltinFunction = "Builtin Function"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type IntegerObject struct {
	Value int64
}

func (i IntegerObject) Type() ObjectType {
	return Integer
}

func (i IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type BooleanObject struct {
	Value bool
}

func (b BooleanObject) Type() ObjectType {
	return Boolean
}

func (b BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type NullObject struct {
}

func (n NullObject) Type() ObjectType {
	return Null
}

func (n NullObject) Inspect() string {
	return "null"
}

type ReturnObject struct {
	Value Object
}

func (r ReturnObject) Type() ObjectType {
	return Return
}

func (r ReturnObject) Inspect() string {
	return r.Value.Inspect()
}

type ErrorObject struct {
	Message string
}

func (e ErrorObject) Type() ObjectType {
	return Error
}

func (e ErrorObject) Inspect() string {
	return "Error: " + e.Message
}

type FunctionObject struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f FunctionObject) Type() ObjectType {
	return Function
}

func (f FunctionObject) Inspect() string {
	b := make([]byte, 0, 10)
	b = append(b, "fn ("...)
	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}
	b = append(b, strings.Join(params, ",")...)
	b = append(b, ") "...)
	b = append(b, f.Body.String()...)

	return string(b)
}

type StringObject struct {
	Value string
}

func (s StringObject) Type() ObjectType {
	return String
}

func (s StringObject) Inspect() string {
	return s.Value
}

type BuiltinFunctionObject struct {
	Function func(objs ...Object) Object
}

func (bf BuiltinFunctionObject) Type() ObjectType {
	return BuiltinFunction
}

func (bf BuiltinFunctionObject) Inspect() string {
	return "builtin function"
}

type ArrayObject struct {
	Elements []Object
}

func (a ArrayObject) Type() ObjectType {
	return Array
}

func (a ArrayObject) Inspect() string {
	b := make([]byte, 0, 10)
	b = append(b, '[')
	elms := make([]string, len(a.Elements))
	for i, elm := range a.Elements {
		elms[i] = elm.Inspect()
	}
	b = append(b, strings.Join(elms, ",")...)
	b = append(b, ']')

	return string(b)
}

type HashKeyable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (i IntegerObject) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

var (
	HashKeyTrue  = HashKey{Type: Boolean, Value: 1}
	HashKeyFalse = HashKey{Type: Boolean, Value: 0}
)

func (b BooleanObject) HashKey() HashKey {
	if b.Value {
		return HashKeyTrue
	}

	return HashKeyFalse
}

func (s StringObject) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}
