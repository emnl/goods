// Package binarytree provides the basic datastructure
// binary search tree. It is not self-balanced.
package binarytree

import (
	"errors"
	"fmt"
	"github.com/emnl/goods/queue"
	"github.com/emnl/goods/stack"
)

// A BinaryTree has a size, a pointer to the root node, and
// a user defined function which is used to compare the node's element.
type BinaryTree struct {
	less LessFunc
	size int
	root *node
}

// The binarytree is made up of nodes with an element,
// a pointer to the left (smaller) node, and a pointer to the right (bigger) node
type node struct {
	elem  Elem
	left  *node
	right *node
}

// Elem is used as a generic for any type of value.
type Elem interface{}

// LessFunc is used as a user function to compare elements in the list.
// It must return true if the first parameter is less then the second.
// False, if the first and second are equal.
type LessFunc func(a, b interface{}) bool

// New is used as an optional constructor for the BinaryTree
// struct.
func New(lf LessFunc) *BinaryTree {
	bt := BinaryTree{lf, 0, nil}
	return &bt
}

// Size returns the size of the tree.
func (T *BinaryTree) Size() int {
	return T.size
}

// Empty returns true if the tree is empty.
func (T *BinaryTree) Empty() bool {
	return T.size == 0
}

// Add adds the given element to the tree.
func (T *BinaryTree) Add(E Elem) error {
	oldsize := T.size
	T.insert(E)
	if oldsize == T.size {
		return errors.New("Item already exists in Tree")
	}
	return nil
}

// Remove removes the given element from the tree.
func (T *BinaryTree) Remove(E Elem) error {
	rem := T.remove(E)
	if !rem {
		return errors.New("Item does not exist in Tree.")
	}

	T.size--
	return nil
}

// Contains returns true if the given element exists within
// the tree.
func (T *BinaryTree) Contains(E Elem) bool {
	return T.get(E) != nil
}

// First returns the leftmost element in the tree.
func (T *BinaryTree) First() Elem {
	if T.Empty() {
		return nil
	}
	return T.root.findMin().elem
}

// Last returns the rightmost element in the tree.
func (T *BinaryTree) Last() Elem {
	if T.Empty() {
		return nil
	}
	return T.root.findMax().elem
}

// PrintTree prints the tree in the console. It is used as a
// debugging tool.
func (T *BinaryTree) PrintTree() {
	if T.Empty() {
		fmt.Println("Empty tree")
		return
	}
	fmt.Print("\n")
	print(T.root, 0)
	fmt.Print("\n")
}

// InOrder returns an iterator over the tree depth-first inorder:
// Visit the root.
// Traverse the left subtree.
// Traverse the right subtree.
func (T *BinaryTree) InOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		nodes := stack.New()
		currentNode := T.root

		for true {
			if currentNode != nil {
				nodes.Push(currentNode)
				currentNode = currentNode.left
			} else {
				if !nodes.Empty() {
					currentNode = nodes.Pop().(*node)
					ch <- currentNode.elem
					currentNode = currentNode.right
				} else {
					break
				}
			}
		}

		close(ch)
	}()
	return ch
}

// PreOrder returns an iterator over the tree depth-first in
// preorder:
// Traverse the left subtree.
// Visit the root.
// Traverse the right subtree.
func (T *BinaryTree) PreOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		if T.Empty() {
			close(ch)
			return
		}

		nodes := stack.New()
		nodes.Push(T.root)

		for !nodes.Empty() {
			currentNode := nodes.Pop().(*node)
			ch <- currentNode.elem

			if currentNode.right != nil {
				nodes.Push(currentNode.right)
			}
			if currentNode.left != nil {
				nodes.Push(currentNode.left)
			}
		}

		close(ch)
	}()
	return ch
}

// PostOrder returns an iterator over the tree depth-first in
// postorder:
// Traverse the left subtree.
// Traverse the right subtree.
// Visit the root.
func (T *BinaryTree) PostOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		if T.Empty() {
			close(ch)
			return
		}

		nodes := stack.New()
		nodes.Push(T.root)
		var prev *node

		for !nodes.Empty() {
			current := nodes.Peek().(*node)

			if prev == nil || prev.left == current || prev.right == current {
				if current.left != nil {
					nodes.Push(current.left)
				} else if current.right != nil {
					nodes.Push(current.right)
				}
			} else if current.left == prev {
				if current.right != nil {
					nodes.Push(current.right)
				}
			} else {
				ch <- current.elem
				nodes.Pop()
			}
			prev = current
		}

		close(ch)
	}()
	return ch
}

// LevelOrder is an iterator over the levels of the tree.
// Also known as breadth-first traversal.
func (T *BinaryTree) LevelOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		if T.Empty() {
			close(ch)
			return
		}

		nodes := queue.New()
		nodes.Offer(T.root)

		for !nodes.Empty() {
			current, _ := nodes.Poll().(*node)

			ch <- current.elem
			if current.left != nil {
				nodes.Offer(current.left)
			}
			if current.right != nil {
				nodes.Offer(current.right)
			}
		}

		close(ch)
	}()
	return ch
}

// get returns the node of the given element.
func (T *BinaryTree) get(E Elem) *node {
	r := T.root
	for r != nil {
		switch {
		case T.less(E, r.elem):
			r = r.left
		case T.less(r.elem, E):
			r = r.right
		default:
			return r
		}
	}
	return nil
}

// insert addeds an element to the correct position within
// the tree.
func (T *BinaryTree) insert(E Elem) {

	if T.root == nil {
		T.root = &node{E, nil, nil}
		T.size += 1
		return
	}

	for root := T.root; root != nil; {
		if E == root.elem {
			return // Duplicate
		} else if T.less(E, root.elem) {
			if root.left == nil {
				root.left = &node{E, nil, nil}
				T.size += 1
				return
			} else {
				root = root.left
			}
		} else {
			if root.right == nil {
				root.right = &node{E, nil, nil}
				T.size += 1
				return
			} else {
				root = root.right
			}
		}
	}
}

// remove deletes a node from the tree based on an input element.
func (T *BinaryTree) remove(E Elem) bool {
	if T.root == nil {
		return false
	} else {
		if T.root.elem == E {
			dummy := &node{nil, nil, nil}
			dummy.left = T.root
			res := T.root.remove(E, dummy, T.less)
			T.root = dummy.left
			return res
		} else {
			return T.root.remove(E, nil, T.less)
		}
	}
	return true
}

// remove deletes a node from a subtree. It returns false if the
// node isn't found within the subtree. If two children exists
// within the subtree, the root is replaced with the smallest
// value in the right subtree. Else, the removed node sets it's parent
// the right values.
func (N *node) remove(E Elem, parent *node, less LessFunc) bool {
	if less(E, N.elem) {
		if N.left != nil {
			return N.left.remove(E, N, less)
		} else {
			return false
		}
	} else if less(N.elem, E) {
		if N.right != nil {
			return N.right.remove(E, N, less)
		} else {
			return false
		}
	} else {
		if N.left != nil && N.right != nil {
			N.elem = N.right.findMin().elem
			N.right.remove(N.elem, N, less)
		} else if parent.left == N {
			if N.left != nil {
				parent.left = N.left
			} else {
				parent.left = N.right
			}
		} else if parent.right == N {
			if N.left != nil {
				parent.right = N.left
			} else {
				parent.right = N.right
			}
		}
	}
	return true
}

// findMax returns the smallest (most left) node in the
// subtree.
func (N *node) findMin() *node {
	found := N
	for found.left != nil {
		found = found.left
	}
	return found
}

// findMax returns the largest (most right) node in the
// subtree.
func (N *node) findMax() *node {
	found := N
	for found.right != nil {
		found = found.right
	}
	return found
}

// print is used with debugging. It prints a simple tree
// representation.
func print(N *node, padding int) {
	if N != nil {
		newp := padding + 3
		print(N.left, newp)
		for i := 0; i < padding; i++ {
			fmt.Print("-")
		}
		fmt.Printf("%d \n", N.elem)
		print(N.right, newp)
	}
}
