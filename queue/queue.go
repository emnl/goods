// Package queue offers a nice and simple interface
// to create and use queues. It relies on a linkedlist
// under the hood and maintains a very low complexity.
package queue

import "github.com/emnl/goods/linkedlist"

// Queue uses a linkedlist to behave as a first-in-first-out
// queue.
type Queue struct {
	linkedlist.LinkedList
}

// Elem is used as a generic for any type of value.
type Elem linkedlist.Elem

// New is used as a constructor for the Queue
// struct.
func New() *Queue {
	return &Queue{*linkedlist.New()}
}

// Offer places an element last in the queue.
func (Q *Queue) Offer(V Elem) {
	Q.AddLast(V)
}

// Poll returns the first element in the queue
// and removes it.
func (Q *Queue) Poll() Elem {
	if Q.Empty() {
		return nil
	}

	result := Q.First()
	Q.RemoveFirst()
	return result
}

// Peek returns the first element in the queue
// without removing it.
func (Q *Queue) Peek() Elem {
	if Q.Empty() {
		return nil
	}

	return Q.First()
}

// Enqueue is a synonym for Offer().
func (Q *Queue) Enqueue(V Elem) { Q.Offer(V) }

// Dequeue is a synonym for Poll().
func (Q *Queue) Dequeue() Elem { return Q.Poll() }
