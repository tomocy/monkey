package object

import "fmt"

const (
	integer = "Integer"
	boolean = "Boolean"
	null    = "Null"
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
	return integer
}

func (i IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type BooleanObject struct {
	Value bool
}

func (b BooleanObject) Type() ObjectType {
	return boolean
}

func (b BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type NullObject struct {
}

func (n NullObject) Type() ObjectType {
	return null
}

func (n NullObject) Inspect() string {
	return "null"
}
