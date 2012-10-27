package redblacktree

import "testing"

func intLess(a, b interface{}) bool {
	return a.(int) < b.(int)
}

func TestNew(t *testing.T) {
	tree := New(intLess)

	if tree.size != 0 || tree.root != nil {
		t.Errorf("New constructor is broken.")
	}
}

func TestSize(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(20)

	if tree.size != 2 {
		t.Errorf("Size should return 2.")
	}
}

func TestEmpty(t *testing.T) {
	tree := New(intLess)

	if !tree.Empty() {
		t.Errorf("Empty should return true.")
	}

	tree.Add(10)

	if tree.Empty() {
		t.Errorf("Empty should return false")
	}
}

func TestAdd(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(20)
	tree.Add(30)
	tree.Add(int(20.0)) // Same as second add

	if tree.size != 3 {
		t.Errorf("Add should add elements.")
	}
}

func TestNRemove(t *testing.T) {
	tree := New(intLess)

	tree.Remove(10)

	if tree.size != 0 {
		t.Errorf("Tree should not be affected.")
	}

	tree.Add(10)
	tree.Remove(20)

	if tree.size != 1 {
		t.Errorf("Nothing should have been removed.")
	}

	if tree.Remove(10) != nil {
		t.Errorf("Remove didn't work.")
	}
}

func TestContains(t *testing.T) {
	tree := New(intLess)
	tree.Add(10)

	if tree.Contains(20) {
		t.Errorf("Tree does not contain 20.")
	}

	if !tree.Contains(10) {
		t.Errorf("Tree does contain 10.")
	}
}

func TestFirst(t *testing.T) {
	tree := New(intLess)

	if tree.First() != nil {
		t.Errorf("An empty tree should return nil on First.")
	}

	for x := 0; x < 10; x++ {
		tree.Add(x)
	}

	if tree.First() != 0 {
		t.Errorf("First element should be 0.")
	}
}

func TestLast(t *testing.T) {
	tree := New(intLess)

	if tree.Last() != nil {
		t.Errorf("An empty tree should return nil on Last.")
	}

	for x := 0; x <= 10; x++ {
		tree.Add(x)
	}

	if tree.Last() != 10 {
		t.Errorf("Last element should be 10.")
	}
}

func TestInOrder(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(5)
	tree.Add(15)

	i := []int{}

	for item := range tree.InOrder() {
		i = append(i, item.(int))
	}

	if i[0] != 5 {
		t.Errorf("InOrder should visit the left subtree first.")
	}
	if i[1] != 10 {
		t.Errorf("InOrder should visit the root after the left subtree.")
	}
	if i[2] != 15 {
		t.Errorf("InOrder should visit the right subtree last.")
	}

}

func TestPreOrder(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(5)
	tree.Add(15)

	i := []int{}

	for item := range tree.PreOrder() {
		i = append(i, item.(int))
	}

	if i[0] != 10 {
		t.Errorf("PreOrder should visit the root first.")
	}
	if i[1] != 5 {
		t.Errorf("PreOrder should visit the left subtree after the root.")
	}
	if i[2] != 15 {
		t.Errorf("PreOrder should visit the right subtree last.")
	}

}

func TestPostOrder(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(5)
	tree.Add(15)

	i := []int{}

	for item := range tree.PostOrder() {
		i = append(i, item.(int))
	}

	if i[0] != 5 {
		t.Errorf("PostOrder should visit the left subtree first.")
	}
	if i[1] != 15 {
		t.Errorf("PostOrder should visit the right subtree after the left subtree.")
	}
	if i[2] != 10 {
		t.Errorf("PostOrder should visit the root last.")
	}

}

func TestLevelOrder(t *testing.T) {
	tree := New(intLess)

	tree.Add(10)
	tree.Add(5)
	tree.Add(15)

	i := []int{}

	for item := range tree.LevelOrder() {
		i = append(i, item.(int))
	}

	if i[0] != 10 {
		t.Errorf("LevelOrder should visit the first (root) level first.")
	}
	if i[1] != 5 {
		t.Errorf("LevelOrder should visit the second level with the left-most node.")
	}
	if i[2] != 15 {
		t.Errorf("LevelOrder should visit the second level with the right-most node last.")
	}

}
