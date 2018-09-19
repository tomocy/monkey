package ast

type modifier func(Node) Node

func Modify(node Node, modifier modifier) Node {
	switch node := node.(type) {
	case *Program:
		return modifyProgram(node, modifier)
	case *ExpressionStatement:
		return modifyExpressionStatement(node, modifier)
	case *Infix:
		return modifyInfix(node, modifier)
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

func modifyInfix(node *Infix, modifier modifier) Node {
	node.LeftValue, _ = Modify(node.LeftValue, modifier).(Expression)
	node.RightValue, _ = Modify(node.LeftValue, modifier).(Expression)

	return node
}
