package tests

import (
	"rb-tree-map/internal/rbtree"
	"testing"
)

func TestGetBasic(t *testing.T) {
	tree := rbtree.New[string, int]()
	val, ok := tree.Get("any_key")
	if ok {
		t.Error("Get on an empty tree should return ok=false")
	}
	if val != 0 {
		t.Errorf("Get on an empty tree should return the zero value for the type, but got %d", val)
	}

	tree.Insert("apple", 10)
	tree.Insert("banana", 20)
	tree.Insert("cherry", 30)

	val, ok = tree.Get("banana")
	if !ok {
		t.Error("Expected to find key 'banana', but it was not found")
	}
	if val != 20 {
		t.Errorf("Expected value for 'banana' to be 20, but got %d", val)
	}

	val, ok = tree.Get("apple")
	if !ok {
		t.Error("Expected to find key 'apple', but it was not found")
	}
	if val != 10 {
		t.Errorf("Expected value for 'apple' to be 10, but got %d", val)
	}

	val, ok = tree.Get("grape")
	if ok {
		t.Errorf("Did not expect to find key 'grape', but got value %d", val)
	}
	if val != 0 {
		t.Errorf("Get for a non-existent key should return the zero value, but got %d", val)
	}
}

func TestGetAfterModification(t *testing.T) {
	tree := rbtree.New[int, string]()
	tree.Insert(100, "original")
	tree.Insert(200, "to_be_deleted")

	tree.Insert(100, "updated")
	val, ok := tree.Get(100)
	if !ok {
		t.Fatal("Key 100 should exist after update")
	}
	if val != "updated" {
		t.Errorf("Expected to get 'updated' value for key 100, but got '%s'", val)
	}

	_, ok = tree.Get(200)
	if !ok {
		t.Fatal("Key 200 should exist before being removed")
	}
	tree.Remove(200)

	val, ok = tree.Get(200)
	if ok {
		t.Errorf("Key 200 should not be found after removal, but got value '%s'", val)
	}

	val, ok = tree.Get(100)
	if !ok || val != "updated" {
		t.Error("Key 100 should still exist and have its correct value after another key was removed")
	}
}

func TestContainsKeyBasic(t *testing.T) {
	tree := rbtree.New[int, bool]()
	if tree.ContainsKey(123) {
		t.Error("ContainsKey on an empty tree should return false")
	}

	tree.Insert(10, true)
	tree.Insert(20, false)
	tree.Insert(30, true)

	if !tree.ContainsKey(20) {
		t.Error("Expected ContainsKey to return true for key 20")
	}

	if !tree.ContainsKey(10) {
		t.Error("Expected ContainsKey to return true for key 10")
	}

	if tree.ContainsKey(99) {
		t.Error("Expected ContainsKey to return false for non-existent key 99")
	}
}

func TestContainsKeyAfterModification(t *testing.T) {
	tree := rbtree.New[string, int]()

	if tree.ContainsKey("apple") {
		t.Fatal("Tree should not contain 'apple' initially")
	}
	tree.Insert("apple", 1)
	if !tree.ContainsKey("apple") {
		t.Error("ContainsKey should return true for 'apple' after it has been inserted")
	}

	tree.Insert("banana", 2)
	if !tree.ContainsKey("banana") {
		t.Fatal("Setup failed: 'banana' should exist before being removed")
	}

	tree.Remove("banana")
	if tree.ContainsKey("banana") {
		t.Error("ContainsKey should return false for 'banana' after it has been removed")
	}

	if !tree.ContainsKey("apple") {
		t.Error("'apple' should still be contained in the tree after 'banana' was removed")
	}
}
