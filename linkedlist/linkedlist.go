// Package linkedlist provides an expandeble and generic
// interface to linkedlists. It uses two-way non-cirkular
// linkedlist and is thread safe.
package linkedlist

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"sync"
)

// A linkedlist has a size, a pointer to the first node, and
// a pointer to the last node in the list.
type LinkedList struct {
	size  int
	first *node
	last  *node
	mu    sync.RWMutex
}

// The linkedlist's chain is made up of nodes with an element,
// a pointer to the previous node, and a pointer to the next node
type node struct {
	Value Elem
	next  *node
	prev  *node
}

// Elem is used as a generic for any type of value.
type Elem interface{}

// New is used as an optional constructor for the LinkedList
// struct.
func New() *LinkedList {
	return &LinkedList{}
}

// Size returns the size of the list.
func (L *LinkedList) Size() int {
	L.mu.RLock()
	defer L.mu.RUnlock()

	return L.size
}

// Len is a synonym for Size().
func (L *LinkedList) Len() int {
	return L.Size()
}

// Empty returns true if the list is empty.
func (L *LinkedList) Empty() bool {
	return L.Size() == 0
}

// AddFirst adds a node at the start of the list with
// the given element.
func (L *LinkedList) AddFirst(V Elem) {
	L.mu.Lock()
	n := node{V, nil, nil}

	if L.size == 0 {
		L.last = &n
	} else {
		n.next = L.first
		L.first.prev = &n
	}

	L.first = &n
	L.size += 1
	L.mu.Unlock()
}

// AddLast adds a node at the end of the list with
// the given element.
func (L *LinkedList) AddLast(V Elem) {
	L.mu.Lock()
	n := node{V, nil, L.last}

	if L.size == 0 {
		L.first = &n
	} else {
		L.last.next = &n
	}

	L.last = &n
	L.size += 1
	L.mu.Unlock()
}

// Contains returns true if the list has at least one
// node with the given element.
func (L *LinkedList) Contains(V Elem) bool {
	L.mu.RLock()
	defer L.mu.RUnlock()

	res := L.fastGet(V)
	if res == nil {
		return false
	}
	return true
}

// Index returns the index of the first occurrence of the
// node with the given element.
func (L *LinkedList) Index(V Elem) int {
	i := 0
	for n := range L.iter() {
		if n == V {
			return i
		}
		i++
	}
	return -1
}

// Get returns the node's element at the given index.
func (L *LinkedList) Get(i int) Elem {
	L.mu.RLock()
	defer L.mu.RUnlock()

	node, err := L.getNode(i)
	if err != nil {
		return nil
	}
	return node.Value
}

// Set updates the node's element at the given index.
func (L *LinkedList) Set(i int, V Elem) error {
	L.mu.Lock()
	defer L.mu.Unlock()

	node, err := L.getNode(i)
	if err == nil {
		node.Value = V
	}
	return err
}

// First returns the first nodes' element.
func (L *LinkedList) First() Elem {
	L.mu.RLock()
	defer L.mu.RUnlock()
	if L.size == 0 {
		return nil
	}

	return L.first.Value
}

// Last returns the last node's element.
func (L *LinkedList) Last() Elem {
	L.mu.RLock()
	defer L.mu.RUnlock()

	if L.size == 0 {
		return nil
	}

	return L.last.Value
}

// RemoveFirst deletes the first node in the
// list.
func (L *LinkedList) RemoveFirst() error {
	L.mu.Lock()
	defer L.mu.Unlock()

	if L.size == 0 {
		return errors.New("List is empty.")
	}

	L.removeNode(L.first)
	return nil
}

// RemoveLast deletes the last node in the
// list.
func (L *LinkedList) RemoveLast() error {
	L.mu.Lock()
	defer L.mu.Unlock()

	if L.size == 0 {
		return errors.New("List is empty.")
	}

	L.removeNode(L.last)
	return nil
}

// Remove deletes the first occurence of a node in
// the linkedlist by its value.
func (L *LinkedList) Remove(V Elem) error {
	L.mu.Lock()
	defer L.mu.Unlock()

	res := L.slowGet(V)
	if res == nil {
		return errors.New("Item not found in list.")
	}

	L.removeNode(res)
	return nil
}

// FastRemove deletes an occurence of a node in the
// linkedlist by its value.
func (L *LinkedList) FastRemove(V Elem) error {
	L.mu.Lock()
	defer L.mu.Unlock()

	res := L.fastGet(V)
	if res == nil {
		return errors.New("Item not found in list.")
	}

	L.removeNode(res)
	return nil
}

// Iter is an iterator to be used for iterate
// the linkedlist (front first) easily.
// Usage: for x := range list.Iter() {}
func (L *LinkedList) Iter() chan Elem {
	L.mu.RLock()
	defer L.mu.RUnlock()

	return L.iter()
}

// ToSlice returns a slice representation of the
// linkedlist.
func (L *LinkedList) ToSlice() []Elem {
	L.mu.RLock()
	defer L.mu.RUnlock()

	res := make([]Elem, L.size)
	i := 0
	for x := range L.iter() {
		res[i] = x
		i++
	}

	return res
}

// FromSlice creates a linkedlist from a go slice.
func FromSlice(slc interface{}) *LinkedList {
	v := reflect.ValueOf(slc)
	newl := New()

	for i := 0; i < v.Len(); i++ {
		newl.AddLast(v.Index(i).Interface())
	}

	return newl
}

// Conc concatenates two linkedlists. The function is
// efficient with a complexity at O(1).
func (L *LinkedList) Conc(other *LinkedList) {
	L.mu.Lock()
	defer L.mu.Unlock()

	if L.size == 0 {
		L.first = other.first
		L.last = other.last
		L.size = other.size
		return
	} else if other.size == 0 {
		return
	}

	L.last.next = other.first
	other.first.prev = L.last
	L.last = other.last
	L.size = L.size + other.size
}

// Serialize WIP TODO
func (L *LinkedList) Serialize() []byte {
	L.mu.RLock()
	defer L.mu.RUnlock()

	m := new(bytes.Buffer)
	gob.NewEncoder(m).Encode(L.ToSlice())

	return m.Bytes()
}

// Deserialize WIP TODO
func Deserialize(bt []byte) *LinkedList {
	p := bytes.NewBuffer(bt)
	dec := gob.NewDecoder(p)

	var slc []Elem
	dec.Decode(&slc)

	return FromSlice(slc)
}

// iter is used internally and is not locked.
func (L *LinkedList) iter() chan Elem {
	ch := make(chan Elem, L.size)
	go func() {
		for n := L.first; n != nil; n = n.next {
			ch <- n.Value
		}
		close(ch)
	}()
	return ch
}

// get searches the list from start to end for the node with
// the given element. This returns the first instance.
func (L *LinkedList) slowGet(E Elem) *node {
	for n := L.first; n != nil; n = n.next {
		if n.Value == E {
			return n
		}
	}
	return nil
}

// fastGet searches the list from both ends concurrently.
// This returns any instance.
func (L *LinkedList) fastGet(E Elem) *node {

	/* Delegate to slower get if the list is small enough */
	if L.size < 1000 {
		return L.slowGet(E)
	}

	found := make(chan *node, 1)
	done := make(chan bool, 2)
	half := L.size / 2

	go func() {
		cur := L.first
		for n := 0; n < half; n++ {
			if E == cur.Value {
				found <- cur
				break
			}
			cur = cur.next
		}
		done <- true
	}()

	go func() {
		cur := L.last
		for n := L.size; n >= half; n-- {
			if E == cur.Value {
				found <- cur
				break
			}
			cur = cur.prev
		}
		done <- true
	}()

	go func() {
		<-done
		<-done
		found <- nil
	}()

	return <-found
}

// getNode retrives a node given an index.
func (L *LinkedList) getNode(i int) (*node, error) {
	if L.size == 0 || i > L.size-1 {
		return nil, errors.New("Index out of bound.")
	}

	var n *node

	if i <= L.size/2 {
		n = L.first
		for p := 0; p != i; p++ {
			n = n.next
		}
	} else {
		n = L.last
		for p := L.size - 1; p != i; p-- {
			n = n.prev
		}
	}

	return n, nil
}

// removeNode deletes the node from the given list.
// The function is considered to be used internally.
func (L *LinkedList) removeNode(N *node) {

	/* Only node */
	if L.size == 1 {
		L.first = nil
		L.last = nil
		L.size--
		return
	}

	/* First node */
	if N.prev == nil {
		N.next.prev = nil
		L.first = N.next
		L.size--
		return
	}

	/* Last node */
	if N.next == nil {
		N.prev.next = nil
		L.last = N.prev
		L.size--
		return
	}

	/* Node in middle of chain */
	N.next.prev = N.prev
	N.prev.next = N.next
	L.size--
	return
}
