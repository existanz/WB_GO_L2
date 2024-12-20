/*
Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*

	Паттерн "посетитель" позволяет добавлять новые операции к объектам, не изменяя их классы.
	При этом отделяя алгоритмы от объектов, над которыми они работают.

	+ гибкость: позволяет добавлять новые операции без изменения существующих классов.
	- сложность: делает код менее читаемым, усложняет структуру кода.
	- гибкость: добавление новых классов объектов, требует изменение всех посетителей.
	- SO из SOLID: нарушеает принципы единственной ответственности и открытости/закрытости.

	Системы отчетности: в системах, где необходимо генерировать разные отчеты по различным типам данных.
  Компиляторы: в компиляторах для выполнения различных операций над узлами синтаксического дерева.
*/

package pattern

import "fmt"

type Node interface {
	Accept(visitor Visitor)
}

type AddNode struct {
	Left  Node
	Right Node
}

func (n *AddNode) Accept(visitor Visitor) {
	visitor.VisitAdd(n)
}

type SubNode struct {
	Left  Node
	Right Node
}

func (n *SubNode) Accept(visitor Visitor) {
	visitor.VisitSub(n)
}

type NumberNode struct {
	Value int
}

func (n *NumberNode) Accept(visitor Visitor) {
	visitor.VisitNumber(n)
}

type Visitor interface {
	VisitAdd(*AddNode)
	VisitSub(*SubNode)
	VisitNumber(*NumberNode)
}

type Evaluator struct {
	Result int
}

func (e *Evaluator) VisitAdd(n *AddNode) {
	n.Left.Accept(e)
	leftResult := e.Result
	n.Right.Accept(e)
	rightResult := e.Result
	e.Result = leftResult + rightResult
}

func (e *Evaluator) VisitSub(n *SubNode) {
	n.Left.Accept(e)
	leftResult := e.Result
	n.Right.Accept(e)
	rightResult := e.Result
	e.Result = leftResult - rightResult
}

func (e *Evaluator) VisitNumber(n *NumberNode) {
	e.Result = n.Value
}

func mainVisitor() {
	root := &AddNode{
		Left: &NumberNode{Value: 5},
		Right: &SubNode{
			Left:  &NumberNode{Value: 10},
			Right: &NumberNode{Value: 3},
		},
	}

	evaluator := &Evaluator{}
	root.Accept(evaluator)
	fmt.Println(evaluator.Result)
}
