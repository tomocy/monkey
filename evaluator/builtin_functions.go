package evaluator

import "github.com/tomocy/monkey/object"

var builtinFns = map[string]*object.BuiltinFunctionObject{
	"len": &object.BuiltinFunctionObject{
		Function: builtinLen,
	},
	"first": &object.BuiltinFunctionObject{
		Function: builtinFirst,
	},
	"last": &object.BuiltinFunctionObject{
		Function: builtinLast,
	},
	"rest": &object.BuiltinFunctionObject{
		Function: builtinRest,
	},
	"push": &object.BuiltinFunctionObject{
		Function: builtinPush,
	},
}

func builtinLen(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("invalid number of arguments to len: expected 1, but got %d", len(objs))
	}
	switch obj := objs[0].(type) {
	case *object.StringObject:
		return &object.IntegerObject{Value: int64(len(obj.Value))}
	case *object.ArrayObject:
		return &object.IntegerObject{Value: int64(len(obj.Elements))}
	default:
		return newError("unknown operation: len(%s)", obj.Type())
	}
}

func builtinFirst(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("invalid number of arguments to first: expected 1, but got %d", len(objs))
	}

	obj := objs[0]
	array, ok := obj.(*object.ArrayObject)
	if !ok {
		return newError("unknown operation: first(%s)", obj.Type())
	}

	if len(array.Elements) <= 0 {
		return nullObj
	}

	return array.Elements[0]
}

func builtinLast(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("invalid number of arguments to last: expected 1, but got %d", len(objs))
	}

	obj := objs[0]
	array, ok := obj.(*object.ArrayObject)
	if !ok {
		return newError("unknown operation: last(%s)", obj.Type())
	}

	if len(array.Elements) <= 0 {
		return nullObj
	}

	return array.Elements[len(array.Elements)-1]
}

func builtinRest(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("invalid number of arguments to rest: expected 1, but got %d", len(objs))
	}

	obj := objs[0]
	array, ok := obj.(*object.ArrayObject)
	if !ok {
		return newError("unknown operation: rest(%s)", obj.Type())
	}

	arrayLen := len(array.Elements)
	if arrayLen <= 0 {
		return nullObj
	}

	if arrayLen == 1 {
		return &object.ArrayObject{Elements: make([]object.Object, 0)}
	}

	newElems := make([]object.Object, arrayLen-1, arrayLen-1)
	copy(newElems, array.Elements[1:])

	return &object.ArrayObject{Elements: newElems}
}

func builtinPush(objs ...object.Object) object.Object {
	if len(objs) != 2 {
		return newError("invalid number of arguments to push: expected 2, but got %d", len(objs))
	}

	srcArray := objs[0]
	newElem := objs[1]
	array, ok := srcArray.(*object.ArrayObject)
	if !ok {
		return newError("unknown operation: push(%s, %s)", srcArray.Type(), newElem.Type())
	}

	arrayLen := len(array.Elements)
	newElems := make([]object.Object, arrayLen+1, arrayLen+1)
	copy(newElems, array.Elements)
	newElems[arrayLen] = newElem

	return &object.ArrayObject{Elements: newElems}
}
