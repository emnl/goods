package linkedlist

import (
	"testing"
)

func TestNew(t *testing.T) {
	list := New()

	if list.Size() != 0 || list.First() != nil {
		t.Errorf("New constructor is broken.")
	}
}

func TestSize(t *testing.T) {
	list := New()
	list.AddFirst(10)

	if list.Size() != 1 {
		t.Errorf("A new list with one added item should have the size of 1.")
	}
}

func TestEmpty(t *testing.T) {
	list := New()

	if !list.Empty() {
		t.Errorf("Empty should return true on a new list.")
	}

	list.AddFirst(10)

	if list.Empty() {
		t.Errorf("Empty should return false when the list is not empty, obviously.")
	}
}

func TestAddFirst(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(15)

	if list.First() != 15 {
		t.Errorf("AddFirst should append an item first in the list.")
	}
}

func TestAddLast(t *testing.T) {
	list := New()

	list.AddLast(10)
	list.AddLast(15)

	if list.Get(1) != 15 {
		t.Errorf("AddLast should append an item last in the list.")
	}
}

func TestContains(t *testing.T) {
	list := New()

	list.AddFirst(10)

	if !list.Contains(10) {
		t.Errorf("Contains should return true if the object exists within the list.")
	}

}

func TestIndex(t *testing.T) {
	list := New()

	if list.Index(10) != -1 {
		t.Errorf("Index should return -1 if the object does not exist within the list.")
	}

	list.AddFirst(10)
	list.AddFirst(5)

	if list.Index(5) != 0 {
		t.Errorf("Index should return the correct index position of the given item.")
	}
}

func TestGet(t *testing.T) {
	list := New()

	if list.First() != nil {
		t.Errorf("Get should return nil if the given index does not exist within the list.")
	}

	list.AddFirst(10)
	list.AddFirst(5)

	if list.First() != 5 {
		t.Errorf("Get should return the item at the given index.")
	}
}

func TestSet(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(5)

	list.Set(0, 15)

	if list.First() != 15 {
		t.Errorf("Set should replace the item at the given index.")
	}
}

func TestFirst(t *testing.T) {
	list := New()

	if list.First() != nil {
		t.Errorf("First should return nil if the list is empty.")
	}

	list.AddFirst(10)
	list.AddFirst(5)

	if list.First() != 5 {
		t.Errorf("First should return the first element in the list.")
	}
}

func TestLast(t *testing.T) {
	list := New()

	if list.Last() != nil {
		t.Errorf("Last should return nil if the list is empty.")
	}

	list.AddFirst(10)
	list.AddFirst(5)

	if list.Last() != 10 {
		t.Errorf("Last should return the last element in the list.")
	}
}
func TestRemoveFirst(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(5)

	list.RemoveFirst()

	if list.First() != 10 {
		t.Errorf("RemoveFirst should remove the first element in the list.")
	}
}

func TestRemoveLast(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(5)

	list.RemoveLast()

	if list.First() != 5 {
		t.Errorf("RemoveLast should remove the last element in the list.")
	}
}

func TestRemove(t *testing.T) {
	list := New()

	list.AddFirst(5)
	list.AddFirst(10)
	list.AddFirst(5)

	list.Remove(5)

	if list.First() != 10 || list.Last() != 5 {
		t.Errorf("Remove should remove the first occurrence of the given element.")
	}
}

func TestFastRemove(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(5)

	list.FastRemove(10)

	if list.Last() != 5 {
		t.Errorf("Remove should remove an occurrence of the given element.")
	}
}

func TestToSlice(t *testing.T) {
	list := New()

	list.AddFirst(10)
	list.AddFirst(5)

	sc := list.ToSlice()

	if sc[0] != 5 || sc[1] != 10 {
		t.Errorf("ToSlice should convert the list to a native slice.")
	}
}

func TestFromSlice(t *testing.T) {
	sc := []int{5, 10}
	list := FromSlice(sc)

	if list.First() != 5 || list.Last() != 10 {
		t.Errorf("FromSlice should create a list from a native slice.")
	}
}

func TestConc(t *testing.T) {
	l1 := New()
	l2 := New()

	l1.AddFirst(1)
	l2.AddFirst(2)

	l1.Conc(l2)

	if l1.First() != 1 || l1.Last() != 2 {
		t.Errorf("Conc should append the given list to the end of the current list.")
	}
}
func TestSerialize(t *testing.T) {
	// TODO
}
func TestDeserialize(t *testing.T) {
	// TODO
}
