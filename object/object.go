package object

import "fmt"

const (
	Integer = "Integer"
	Boolean = "Boolean"
	Null    = "Null"
	Return  = "Return"
	Error   = "Error"
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
