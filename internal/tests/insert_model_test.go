package tests

import (
	"rb-tree-map/internal/rbtree"
	"slices"
	"testing"
)

func TestInsertAndInOrderTraversal(t *testing.T) {
	tree := rbtree.New[int, string]()

	keysToInsert := []int{10, 85, 15, 70, 20, 60, 30, 50, 65, 80, 90, 40, 5, 55}

	for _, key := range keysToInsert {
		tree.Insert(key, "v"+string(rune(key)))
	}

	expectedOrder := []int{5, 10, 15, 20, 30, 40, 50, 55, 60, 65, 70, 80, 85, 90}

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if tree.Size() != len(expectedOrder) {
		t.Errorf("Expected tree size to be %d, but got %d", len(expectedOrder), tree.Size())
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		t.Errorf("In-order traversal is incorrect.\nExpected: %v\nGot:      %v", expectedOrder, actualOrder)
	}
}

func TestInsertDuplicates(t *testing.T) {
	tree := rbtree.New[string, int]()

	tree.Insert("apple", 10)
	tree.Insert("banana", 20)
	tree.Insert("cherry", 30)

	initialSize := tree.Size()
	if initialSize != 3 {
		t.Fatalf("Expected initial size to be 3, got %d", initialSize)
	}

	tree.Insert("banana", 99)

	if tree.Size() != initialSize {
		t.Errorf("Expected size to remain %d after inserting a duplicate key, but got %d", initialSize, tree.Size())
	}

	value, found := tree.Get("banana")
	if !found {
		t.Error("Expected to find key 'banana' after update, but it was not found")
	}
	if value != 99 {
		t.Errorf("Expected value for key 'banana' to be updated to 99, but got %d", value)
	}

	value, _ = tree.Get("apple")
	if value != 10 {
		t.Errorf("Value for key 'apple' should not have changed, got %d", value)
	}
}

func TestInsertWithNegativeAndZeroValues(t *testing.T) {
	tree := rbtree.New[int, bool]()

	keysToInsert := []int{10, -5, 0, 20, -15, 5}

	for _, key := range keysToInsert {
		tree.Insert(key, true)
	}

	expectedOrder := []int{-15, -5, 0, 5, 10, 20}

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		t.Errorf("In-order traversal with negative numbers is incorrect.\nExpected: %v\nGot:      %v", expectedOrder, actualOrder)
	}
}
