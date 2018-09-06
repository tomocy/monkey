package object

import "fmt"

const (
	integer = "Integer"
	boolean = "Boolean"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i Integer) Type() ObjectType {
	return integer
}

func (i Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b Boolean) Type() ObjectType {
	return boolean
}

func (b Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
