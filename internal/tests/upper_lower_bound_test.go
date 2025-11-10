package tests

import (
	"rb-tree-map/internal/rbtree"
	"testing"
)

func createTestTree() *rbtree.RBTreeMap[int, string] {
	tree := rbtree.New[int, string]()
	keys := []int{30, 20, 50, 10, 60}
	for _, k := range keys {
		tree.Insert(k, "v"+string(rune(k)))
	}
	return tree
}

func TestLowerBound(t *testing.T) {
	tree := createTestTree()

	t.Run("Empty Tree", func(t *testing.T) {
		emptyTree := rbtree.New[int, string]()
		_, _, ok := emptyTree.LowerBound(10)
		if ok {
			t.Error("LowerBound on an empty tree should return ok=false")
		}
	})

	t.Run("Key Exists", func(t *testing.T) {
		k, v, ok := tree.LowerBound(30)
		if !ok || k != 30 || v != "v"+string(rune(30)) {
			t.Errorf("Expected LowerBound(30) to find key 30, but got k=%v, v=%v, ok=%v", k, v, ok)
		}
	})

	t.Run("Key Does Not Exist (In the middle)", func(t *testing.T) {
		k, v, ok := tree.LowerBound(35)
		if !ok || k != 50 || v != "v"+string(rune(50)) {
			t.Errorf("Expected LowerBound(35) to find key 50, but got k=%v, v=%v, ok=%v", k, v, ok)
		}
	})

	t.Run("Key Smaller Than All Elements", func(t *testing.T) {
		k, _, ok := tree.LowerBound(5)
		if !ok || k != 10 {
			t.Errorf("Expected LowerBound(5) to find the minimum key 10, but got k=%v, ok=%v", k, ok)
		}
	})

	t.Run("Key Larger Than All Elements", func(t *testing.T) {
		_, _, ok := tree.LowerBound(100)
		if ok {
			t.Error("LowerBound for a key larger than all elements should return ok=false")
		}
	})

	t.Run("Key is the Minimum Element", func(t *testing.T) {
		k, _, ok := tree.LowerBound(10)
		if !ok || k != 10 {
			t.Errorf("Expected LowerBound(10) to find key 10, but got k=%v, ok=%v", k, ok)
		}
	})

	t.Run("Key is the Maximum Element", func(t *testing.T) {
		k, _, ok := tree.LowerBound(60)
		if !ok || k != 60 {
			t.Errorf("Expected LowerBound(60) to find key 60, but got k=%v, ok=%v", k, ok)
		}
	})
}

func TestUpperBound(t *testing.T) {
	tree := createTestTree()

	t.Run("Empty Tree", func(t *testing.T) {
		emptyTree := rbtree.New[int, string]()
		_, _, ok := emptyTree.UpperBound(10)
		if ok {
			t.Error("UpperBound on an empty tree should return ok=false")
		}
	})

	t.Run("Key Exists", func(t *testing.T) {
		k, v, ok := tree.UpperBound(30)
		if !ok || k != 50 || v != "v"+string(rune(50)) {
			t.Errorf("Expected UpperBound(30) to find key 50, but got k=%v, v=%v, ok=%v", k, v, ok)
		}
	})

	t.Run("Key Does Not Exist (In the middle)", func(t *testing.T) {
		k, v, ok := tree.UpperBound(35)
		if !ok || k != 50 || v != "v"+string(rune(50)) {
			t.Errorf("Expected UpperBound(35) to find key 50, but got k=%v, v=%v, ok=%v", k, v, ok)
		}
	})

	t.Run("Key Smaller Than All Elements", func(t *testing.T) {
		k, _, ok := tree.UpperBound(5)
		if !ok || k != 10 {
			t.Errorf("Expected UpperBound(5) to find the minimum key 10, but got k=%v, ok=%v", k, ok)
		}
	})

	t.Run("Key is the Maximum Element", func(t *testing.T) {
		_, _, ok := tree.UpperBound(60)
		if ok {
			t.Error("UpperBound for the maximum element should return ok=false")
		}
	})

	t.Run("Key Larger Than All Elements", func(t *testing.T) {
		_, _, ok := tree.UpperBound(100)
		if ok {
			t.Error("UpperBound for a key larger than all elements should return ok=false")
		}
	})

	t.Run("Key Right Before an Existing Element", func(t *testing.T) {
		k, _, ok := tree.UpperBound(29)
		if !ok || k != 30 {
			t.Errorf("Expected UpperBound(29) to find key 30, but got k=%v, ok=%v", k, ok)
		}
	})
}
