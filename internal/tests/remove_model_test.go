package tests

import (
	"math/rand"
	"rb-tree-map/internal/rbtree"
	"slices"
	"testing"
	"time"
)

func TestRemoveAndInOrderTraversal(t *testing.T) {
	tree := rbtree.New[int, string]()
	initialKeys := []int{10, 85, 15, 70, 20, 60, 30, 50, 65, 80, 90, 40, 5, 55}
	for _, key := range initialKeys {
		tree.Insert(key, "v"+string(rune(key)))
	}

	keysToRemove := []int{85, 10, 5, 50}
	for _, key := range keysToRemove {
		tree.Remove(key)
	}

	expectedOrder := []int{15, 20, 30, 40, 55, 60, 65, 70, 80, 90}

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if tree.Size() != len(expectedOrder) {
		t.Errorf("Expected tree size to be %d after removals, but got %d", len(expectedOrder), tree.Size())
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		t.Errorf("In-order traversal is incorrect after removals.\nExpected: %v\nGot:      %v", expectedOrder, actualOrder)
	}
}

func TestRemoveNonExistentKey(t *testing.T) {
	tree := rbtree.New[int, bool]()
	keys := []int{10, 20, 30}
	for _, k := range keys {
		tree.Insert(k, true)
	}

	initialSize := tree.Size()
	initialOrder := make([]int, 0, initialSize)
	for k, _ := range tree.InOrder() {
		initialOrder = append(initialOrder, k)
	}

	tree.Remove(99)

	if tree.Size() != initialSize {
		t.Errorf("Expected size to remain %d after removing a non-existent key, but got %d", initialSize, tree.Size())
	}

	currentOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		currentOrder = append(currentOrder, k)
	}

	if !slices.Equal(initialOrder, currentOrder) {
		t.Errorf("Tree content should not change after attempting to remove a non-existent key.\nBefore: %v\nAfter:  %v", initialOrder, currentOrder)
	}
}

func TestRemoveAllElements(t *testing.T) {
	tree := rbtree.New[int, int]()

	keysToInsert := []int{4, 2, 6, 1, 3, 5, 7}
	for _, k := range keysToInsert {
		tree.Insert(k, k)
	}

	if tree.Size() != len(keysToInsert) {
		t.Fatalf("Tree size after insertion is incorrect")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(keysToInsert), func(i, j int) {
		keysToInsert[i], keysToInsert[j] = keysToInsert[j], keysToInsert[i]
	})

	for _, k := range keysToInsert {
		tree.Remove(k)
	}

	if tree.Size() != 0 {
		t.Errorf("Expected tree size to be 0 after removing all elements, but got %d", tree.Size())
	}

	count := 0
	for _, _ = range tree.InOrder() {
		count++
	}
	if count != 0 {
		t.Errorf("In-order traversal of an empty tree should yield 0 items, but got %d", count)
	}
}
