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
	case *ReturnStatement:
		return modifyReturnStatement(node, modifier)
	case *If:
		return modifyIf(node, modifier)
	case *Prefix:
		return modifyPrefix(node, modifier)
	case *Infix:
		return modifyInfix(node, modifier)
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

func modifySubscript(node *Subscript, modifier modifier) Node {
	node.LeftValue, _ = Modify(node.LeftValue, modifier).(Expression)
	node.Index, _ = Modify(node.Index, modifier).(Expression)

	return node
}
