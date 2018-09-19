package ast

type modifier func(Node) Node

func Modify(node Node, modifier modifier) Node {
	switch node := node.(type) {
	case *Program:
		return modifyProgram(node, modifier)
	case *ExpressionStatement:
		return modifyExpressionStatement(node, modifier)
	case *BlockStatement:
		return modifyBlockStatement(node, modifier)
	case *LetStatement:
		return modifyLetStatement(node, modifier)
	case *ReturnStatement:
		return modifyReturnStatement(node, modifier)
	case *If:
		return modifyIf(node, modifier)
	case *Prefix:
		return modifyPrefix(node, modifier)
	case *Infix:
		return modifyInfix(node, modifier)
	case *Function:
		return modifyFunction(node, modifier)
	case *Array:
		return modifyArray(node, modifier)
	case *Hash:
		return modifyHash(node, modifier)
	case *Subscript:
		return modifySubscript(node, modifier)
	default:
		return modifier(node)
	}
}

func modifyProgram(node *Program, modifier modifier) Node {
	for i, stmt := range node.Statements {
		node.Statements[i], _ = Modify(stmt, modifier).(Statement)
	}

	return node
}

func modifyExpressionStatement(node *ExpressionStatement, modifier modifier) Node {
	node.Value, _ = Modify(node.Value, modifier).(Expression)

	return node
}

func modifyBlockStatement(node *BlockStatement, modifier modifier) Node {
	for i, stmt := range node.Statements {
		node.Statements[i], _ = Modify(stmt, modifier).(Statement)
	}

	return node
}

func modifyLetStatement(node *LetStatement, modifier modifier) Node {
	node.Value, _ = Modify(node.Value, modifier).(Expression)

	return node
}

func modifyReturnStatement(node *ReturnStatement, modifier modifier) Node {
	node.Value, _ = Modify(node.Value, modifier).(Expression)

	return node
}

func modifyIf(node *If, modifier modifier) Node {
	node.Condition, _ = Modify(node.Condition, modifier).(Expression)
	node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
	if node.Alternative != nil {
		node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
	}

	return node
}

func modifyPrefix(node *Prefix, modifier modifier) Node {
	node.RightValue, _ = Modify(node.RightValue, modifier).(Expression)

	return node
}

func modifyInfix(node *Infix, modifier modifier) Node {
	node.LeftValue, _ = Modify(node.LeftValue, modifier).(Expression)
	node.RightValue, _ = Modify(node.LeftValue, modifier).(Expression)

	return node
}

func modifyFunction(node *Function, modifier modifier) Node {
	for i, param := range node.Parameters {
		node.Parameters[i], _ = Modify(param, modifier).(*Identifier)
	}
	node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	return node
}

func modifyArray(node *Array, modifier modifier) Node {
	for i, elem := range node.Elements {
		node.Elements[i], _ = Modify(elem, modifier).(Expression)
	}

	return node
}

func modifyHash(node *Hash, modifier modifier) Node {
	values := make(map[Expression]Expression)
	for key, value := range node.Values {
		newKey, _ := Modify(key, modifier).(Expression)
		newValue, _ := Modify(value, modifier).(Expression)
		values[newKey] = newValue
	}

	node.Values = values

	return node
}

func modifySubscript(node *Subscript, modifier modifier) Node {
	node.LeftValue, _ = Modify(node.LeftValue, modifier).(Expression)
	node.Index, _ = Modify(node.Index, modifier).(Expression)

	return node
}
