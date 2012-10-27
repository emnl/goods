// Package redblacktree provides a self-balanced
// red-black tree datastructure.
package redblacktree

import (
	"errors"
	"fmt"
	"github.com/emnl/goods/queue"
	"github.com/emnl/goods/stack"
	"math"
)

// A redblacktree has a size, a pointer to the root node, and
// a user defined function which is used to compare the node's element.
//
// It has the following requirements:
// 1. A node is either red or black.
// 2. The root is black.
// 3. All leaves are the same color as the root.
// 4. Both children of every red node are black.
// 5. Every simple path from a given node to any of its descendant leaves
//    contains the same number of black nodes.
//
type RedBlackTree struct {
	less LessFunc
	size int
	root *node
}

// The redblacktree is made up of nodes with an element,
// a pointer to the left (smaller) node, a pointer to the right (bigger) node,
// a pointer to the parent node, and a color (red/black).
type node struct {
	elem   Elem
	left   *node
	right  *node
	parent *node
	red    bool
}

// Elem is used as a generic for any type of value.
type Elem interface{}

// LessFunc is used as a user function to compare elements in the list.
// It must return true if the first parameter is less then the second.
// False, if the first and second are equal.
//
// e.g. intLess func(a,b interface{}) { return (a.(int) < b.(int)) }
//
type LessFunc func(a, b interface{}) bool

// New is used as an optional constructor for the BinaryTree
// struct.
//
// e.g. mytree := redblacktree.New(intLess)
//
func New(lf LessFunc) *RedBlackTree {
	rbt := RedBlackTree{lf, 0, nil}
	return &rbt
}

// Size returns the size of the Tree.
//
// e.g. (2 (1) (3)).Size() => 3
//
func (T *RedBlackTree) Size() int {
	return T.size
}

// Empty returns true if the Tree is empty.
//
// e.g. (2 (1) (3)).Empty() => false
//      ().Empty() => true
//
func (T *RedBlackTree) Empty() bool {
	return T.root == nil
}

// Add inserts an element into the Tree
// and keeps the invariant of a redblacktree.
//
// e.g. (2 () ()).Add(3) => (2 () (3))
//
func (T *RedBlackTree) Add(E Elem) error {
	oldsize := T.size
	T.insert(E)
	if oldsize == T.size {
		return errors.New("Item already exists in Tree.")
	}
	return nil
}

// Remove deletes an element from the Tree
// and keeps the invariant of a redblacktree.
//
// e.g. (2 (1) (3)).Remove(2) => (1 () (3))
//
func (T *RedBlackTree) Remove(E Elem) error {
	oldsize := T.size
	T.delete(E)
	if oldsize == T.size {
		return errors.New("Item not found in Tree.")
	}
	return nil
}

// Contains returns true if the given element exists
// within the Tree.
//
// e.g. (2 (1) (3)).Contains(1) => true
//      (2 (1) (3)).Contains(4) => false
//
func (T *RedBlackTree) Contains(E Elem) bool {
	return T.get(E) != nil
}

// First returns the left-most (smallest) element in the Tree.
//
// e.g. (2 (1) (3)).First() => 1
//
func (T *RedBlackTree) First() Elem {
	if T.Empty() {
		return nil
	}
	return T.root.findMin().elem
}

// Last returns the right-most (largest) element in the Tree.
//
// e.g. (2 (1) (3)).Last() => 3
//
func (T *RedBlackTree) Last() Elem {
	if T.Empty() {
		return nil
	}
	return T.root.findMax().elem
}

// Depth returns the logical depth of the Tree.
//
// e.g. Log2(tree.Size())
//
func (T *RedBlackTree) Depth() float64 {
	return math.Log2(float64(T.size))
}

// Height is a synonym for Depth().
func (T *RedBlackTree) Height() float64 {
	return T.Depth()
}

// InOrder returns an iterator over the tree depth-first inorder:
// Traverse the left subtree.
// Visit the root.
// Traverse the right subtree.
//
// e.g. for x := range (2 (1) (3)).InOrder() { x } => 1, 2, 3
//
func (T *RedBlackTree) InOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		nodes := stack.New()
		currentNode := T.root

		for {
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
// Visit the root.
// Traverse the left subtree.
// Traverse the right subtree.
//
// e.g. for x := range (2 (1) (3)).PreOrder() { x } => 2, 1, 3
//
func (T *RedBlackTree) PreOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		if T.Empty() {
			close(ch)
			return
		}

		nodes := stack.New()
		nodes.Push(T.root)

		for !nodes.Empty() {
			currentNode, _ := nodes.Pop().(*node)

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
//
// e.g. for x := range (2 (1) (3)).PostOrder() { x } => 1, 3, 2
//
func (T *RedBlackTree) PostOrder() chan Elem {
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
//
// e.g. for x := range (2 (1) (3)).LevelOrder() { x } => 2, 1, 3
//
func (T *RedBlackTree) LevelOrder() chan Elem {
	ch := make(chan Elem, T.size)
	go func() {

		if T.Empty() {
			close(ch)
			return
		}

		nodes := queue.New()
		nodes.Offer(T.root)

		for !nodes.Empty() {
			current := nodes.Poll().(*node)
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

// PrintTree prints the tree in the console. It is used as a
// debugging tool.
func (T *RedBlackTree) PrintTree() {
	if T.Empty() {
		fmt.Println("Empty tree")
		return
	}
	fmt.Print("\n")
	print(T.root, 0)
	fmt.Print("\n")
}

// isRed returns true if the given node is red.
// The leafs of a redblacktree are always considered black,
// therefore nil return false. This is important.
func isRed(n *node) bool {
	if n == nil {
		return false // leaf nodes are considered black
	}
	return n.red
}

// get returns the node given an element.
func (T *RedBlackTree) get(E Elem) *node {
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

// rotateLeft replaces the given node with the right node
// and then rotates the subtree to the left.
//
//		1
//		 \
//		  2
//		   \
//		    3
//
//		2
//	   / \
//	  1   3
//
func (T *RedBlackTree) rotateLeft(n *node) {
	right := n.right
	T.replaceNode(n, right)
	n.right = right.left
	if right.left != nil {
		right.left.parent = n
	}
	right.left = n
	n.parent = right
}

// rotateRight replaces the given node with the left node
// and then rotates the subtree to the right.
//
//		1
//	   /
//	  2
//	 /
//	3
//
//		2
//	   / \
//	  1   3
//
func (T *RedBlackTree) rotateRight(n *node) {
	left := n.left
	T.replaceNode(n, left)
	n.left = left.right
	if left.right != nil {
		left.right.parent = n
	}
	left.right = n
	n.parent = left
}

// replaceNode replaces an old node for a new one and
// keeps the order in the Tree.
func (T *RedBlackTree) replaceNode(oldn, newn *node) {
	if oldn.parent == nil {
		T.root = newn // the old node was the Tree-root
	} else {
		if oldn == oldn.parent.left {
			oldn.parent.left = newn
		} else {
			oldn.parent.right = newn
		}
	}
	if newn != nil {
		newn.parent = oldn.parent
	}
}

// insert takes the given element and inserts
// it into the Tree. A new node is always inserted as
// red.
func (T *RedBlackTree) insert(E Elem) {
	newn := &node{E, nil, nil, nil, true}

	if T.root == nil {
		T.root = newn
	} else {
		n := T.root
		for true {
			if T.less(newn.elem, n.elem) {
				if n.left == nil {
					n.left = newn
					break
				} else {
					n = n.left
				}
			} else if T.less(n.elem, newn.elem) {
				if n.right == nil {
					n.right = newn
					break
				} else {
					n = n.right
				}
			} else {
				n.elem = newn.elem
				return
			}
		}
		newn.parent = n
	}

	T.size += 1 // A node will be added
	T.insertCase1(newn)
}

// insertCase1 keeps the redblacktree invariant:
// "the root node must always be black".
func (T *RedBlackTree) insertCase1(newn *node) {
	if newn.parent == nil {
		newn.red = false
	} else {
		T.insertCase2(newn)
	}
}

// insertCase2 breakes the insert if the tree is valid.
func (T *RedBlackTree) insertCase2(newn *node) {
	if isRed(newn.parent) == false {
		// if parent is black
		return // Valid tree
	} else {
		T.insertCase3(newn)
	}
}

// insertCase3 repaints the parent and uncle if both are red.
// Also, their grandparent becomes red.
func (T *RedBlackTree) insertCase3(newn *node) {
	if newn.uncle() != nil && newn.uncle().red {
		newn.parent.red = false
		newn.uncle().red = false
		newn.grandparent().red = true
		T.insertCase1(newn.grandparent())
	} else {
		T.insertCase4(newn)
	}
}

// insertCase4 takes care of the situation where the parent is red
// but the uncle is black. Also, the new node is a left child.
// It rotates the Tree to fit the requirements.
func (T *RedBlackTree) insertCase4(newn *node) {
	if newn == newn.parent.right && newn.parent == newn.grandparent().left {
		T.rotateLeft(newn.parent)
		newn = newn.left
	} else if newn == newn.parent.left && newn.parent == newn.grandparent().right {
		T.rotateRight(newn.parent)
		newn = newn.right
	}
	T.insertCase5(newn)
}

// insertCase5 takes care of the situation where the parent is red
// but the uncle is black. Also, the new node is a right child.
// It rotates the Tree to fit the requirements.
func (T *RedBlackTree) insertCase5(newn *node) {
	newn.parent.red = false
	newn.grandparent().red = true

	if newn == newn.parent.left {
		T.rotateRight(newn.grandparent())
	} else {
		T.rotateLeft(newn.grandparent())
	}
}

// delete removes a node from the Tree given an input
// element.
func (T *RedBlackTree) delete(E Elem) {
	dnode := T.get(E)

	if T.Empty() || dnode == nil {
		return
	}

	if dnode.left != nil && dnode.right != nil {
		pred := dnode.left.findMax()
		dnode.elem = pred.elem
		dnode = pred
	}

	var child *node
	if dnode.right == nil {
		child = dnode.left
	} else {
		child = dnode.right
	}

	if !isRed(dnode) {
		dnode.red = isRed(child)
		T.deleteCase1(dnode)
	}
	T.replaceNode(dnode, child)

	if isRed(T.root) {
		T.root.red = false
	}

	T.size -= 1
}

// deleteCase1 checks if the deleted node is the root.
// If it is, we're done.
func (T *RedBlackTree) deleteCase1(dnode *node) {
	if dnode.parent == nil {
		return
	} else {
		T.deleteCase2(dnode)
	}
}

// deleteCase2 rotates and repaint if the input node is red.
func (T *RedBlackTree) deleteCase2(dnode *node) {
	if isRed(dnode.sibling()) {
		dnode.parent.red = true
		dnode.sibling().red = false
		if dnode == dnode.parent.left {
			T.rotateLeft(dnode.parent)
		} else {
			T.rotateRight(dnode.parent)
		}
	}
	T.deleteCase3(dnode)
}

// deleteCase3 handles the case where the input node,
// the parent node and the input node's children are black.
func (T *RedBlackTree) deleteCase3(dnode *node) {
	if isRed(dnode.parent) == false &&
		isRed(dnode.sibling()) == false &&
		isRed(dnode.sibling().left) == false &&
		isRed(dnode.sibling().right) == false {
		dnode.sibling().red = true
		T.deleteCase1(dnode.parent)
	} else {
		T.deleteCase4(dnode)
	}
}

// deleteCase4 the input node and its children are black,
// but the parent is red.
func (T *RedBlackTree) deleteCase4(dnode *node) {
	if isRed(dnode.parent) &&
		isRed(dnode.sibling()) == false &&
		isRed(dnode.sibling().left) == false &&
		isRed(dnode.sibling().right) == false {
		dnode.sibling().red = true
		dnode.parent.red = false
	} else {
		T.deleteCase5(dnode)
	}
}

// deleteCase5 the input node is black but its left child is
// red.
func (T *RedBlackTree) deleteCase5(dnode *node) {
	if dnode == dnode.parent.left &&
		isRed(dnode.sibling()) == false &&
		isRed(dnode.sibling().left) == true &&
		isRed(dnode.sibling().right) == false {
		dnode.sibling().red = true
		dnode.sibling().left.red = false
		T.rotateRight(dnode.sibling())
	} else if dnode == dnode.parent.right &&
		isRed(dnode.sibling()) == false &&
		isRed(dnode.sibling().right) == true &&
		isRed(dnode.sibling().left) == false {
		dnode.sibling().red = true
		dnode.sibling().right.red = false
		T.rotateLeft(dnode.sibling())
	}
	T.deleteCase6(dnode)
}

// deleteCase5 the input node is black but its right child is
// red.
func (T *RedBlackTree) deleteCase6(dnode *node) {
	dnode.sibling().red = isRed(dnode.parent)
	dnode.parent.red = false

	if dnode == dnode.parent.left {
		dnode.sibling().right.red = false
		T.rotateLeft(dnode.parent)
	} else {
		dnode.sibling().left.red = false
		T.rotateRight(dnode.parent)
	}
}

// findMax returns the rightmost (biggest) node in
// the subtree.
func (N *node) findMax() *node {
	found := N
	for found.right != nil {
		found = found.right
	}
	return found
}

// findMin returns the leftmost (smallest) node in
//the subtree.
func (N *node) findMin() *node {
	found := N
	for found.left != nil {
		found = found.left
	}
	return found
}

// uncle returns the parent's sibling().
func (N *node) uncle() *node {
	if N.parent == nil {
		return nil
	}
	return N.parent.sibling()
}

// grandparent returns the parentnode's parent.
func (N *node) grandparent() *node {
	if N.parent == nil {
		return nil
	}
	return N.parent.parent
}

// sibling returns the parent's other child.
func (N *node) sibling() *node {
	if N.parent == nil {
		return nil
	}
	if N.parent.left == N {
		return N.parent.right
	} else {
		return N.parent.left
	}
	return nil
}

// print is used with debugging. It prints a simple tree
// representation.
func print(N *node, padding int) {
	if N != nil {
		newp := padding + 5
		print(N.right, newp)
		for i := 0; i < padding; i++ {
			fmt.Print("-")
		}
		if N.red {
			fmt.Printf("(%d) \n", N.elem)
		} else {
			fmt.Printf("|%d| \n", N.elem)
		}
		print(N.left, newp)
	}
}
