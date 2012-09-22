package stack

import (
	"testing"
)

func TestPush(t *testing.T) {
	stack := New()

	stack.Push(10)
	stack.Push(20)

	if stack.First() != 20 {
		t.Errorf("Push should add the item on top (first) of the stack.")
	}
}

func TestPop(t *testing.T) {
	stack := New()

	stack.Push(10)
	stack.Push(20)

	if stack.Pop() != 20 || stack.Size() != 1 {
		t.Errorf("Pop should return the first item on the stack and then remove it.")
	}
}

func TestPeek(t *testing.T) {
	stack := New()

	stack.Push(10)

	if stack.Peek() != 10 || stack.Empty() {
		t.Errorf("Peek should return the first item on the stack, but not remove it.")
	}
}
