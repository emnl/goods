package queue

import (
	"testing"
)

func TestOffer(t *testing.T) {
	queue := New()

	queue.Offer(10)
	queue.Offer(20)

	if queue.Last() != 20 {
		t.Errorf("Offer should add an item last in queue.")
	}
}

func TestPoll(t *testing.T) {
	queue := New()

	queue.Offer(10)

	if queue.Poll() != 10 || !queue.Empty() {
		t.Errorf("Poll should remove and return the first element in the queue.")
	}

	if queue.Poll() != nil {
		t.Errorf("Poll should return nil if the queue is empty.")
	}
}

func TestPeek(t *testing.T) {
	queue := New()

	queue.Offer(10)

	if queue.Peek() != 10 || queue.Empty() {
		t.Errorf("Peek should return the first value, but not remove it.")
	}
}
