// Package Stack offers a nice and simple interface
// to create and use stacks. It relies on a linkedlist
// under the hood and maintains a very low complexity.
package stack

import "github.com/emnl/goods/linkedlist"

// Stack uses a linkedlist to behave as a last-in-first-out
// stack.
type Stack struct {
	linkedlist.LinkedList
}

// Elem is used as a generic for any type of value.
type Elem linkedlist.Elem

// New is used as a constructor for the Stack
// struct.
func New() *Stack {
	return &Stack{*linkedlist.New()}
}

// Push pushes an element onto the stack.
func (S *Stack) Push(V Elem) {
	S.AddFirst(V)
}

// Pop returns the first element on the stack and removes it
func (S *Stack) Pop() Elem {
	if S.Empty() {
		return nil
	}

	result := S.First()
	S.RemoveFirst()
	return result
}

// Peek returns the first element on the stack
// without removing it.
func (S *Stack) Peek() Elem {
	if S.Empty() {
		return nil
	}

	return S.First()
}
