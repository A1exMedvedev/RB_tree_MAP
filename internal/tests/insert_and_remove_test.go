package tests

import (
	"math/rand"
	"rb-tree-map/internal/rbtree"
	"slices"
	"testing"
	"time"
)

func TestCombinedInsertRemoveSequence(t *testing.T) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	t.Logf("Running combined test with random seed: %d", seed)

	tree := rbtree.New[int, bool]()
	trackingMap := make(map[int]struct{})

	initialSize := 100
	churnCycles := 500
	for i := 0; i < initialSize; i++ {
		key := rng.Intn(initialSize * 10)
		if _, exists := trackingMap[key]; exists {
			i--
			continue
		}
		tree.Insert(key, true)
		trackingMap[key] = struct{}{}
	}

	if tree.Size() != len(trackingMap) {
		t.Fatalf("Initial population failed. Expected size %d, got %d", len(trackingMap), tree.Size())
	}
	for i := 0; i < churnCycles; i++ {
		if rng.Intn(2) == 0 {
			key := rng.Intn(initialSize * 10)
			tree.Insert(key, true)
			trackingMap[key] = struct{}{}
		} else {
			if len(trackingMap) == 0 {
				continue
			}
			keysInMap := getKeysFromMap(trackingMap)
			keyToRemove := keysInMap[rng.Intn(len(keysInMap))]

			tree.Remove(keyToRemove)
			delete(trackingMap, keyToRemove)
		}
	}
	if tree.Size() != len(trackingMap) {
		t.Errorf("Final size is incorrect after churn. Expected %d, got %d", len(trackingMap), tree.Size())
	}
	expectedOrder := getKeysFromMap(trackingMap)
	slices.Sort(expectedOrder)

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		t.Errorf("In-order traversal is incorrect after combined operations.\nExpected: %v\nGot:      %v", expectedOrder, actualOrder)
	}
}

func getKeysFromMap(m map[int]struct{}) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestLargeScaleBuildUpAndRandomTearDown(t *testing.T) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	t.Logf("Running large build-up/tear-down test with seed: %d", seed)

	tree := rbtree.New[int, int]()
	const numElements = 5000

	keys := rng.Perm(numElements)
	for _, key := range keys {
		tree.Insert(key, key)
	}

	if tree.Size() != numElements {
		t.Fatalf("Expected size %d after insertion, but got %d", numElements, tree.Size())
	}

	expectedOrder := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		expectedOrder[i] = i
	}

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		t.Fatal("Tree is not correctly sorted after large-scale insertion.")
	}

	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for i, key := range keys {
		tree.Remove(key)
		expectedSize := numElements - i - 1
		if tree.Size() != expectedSize {
			t.Fatalf("Size mismatch during removal of key %d. Expected %d, got %d", key, expectedSize, tree.Size())
		}
	}

	if tree.Size() != 0 {
		t.Errorf("Expected tree to be empty after removing all elements, but size is %d", tree.Size())
	}
}

func TestSequentialInsertAndRemove(t *testing.T) {
	const sequenceSize = 2500
	t.Run("AscInsert_AscRemove", func(t *testing.T) {
		tree := rbtree.New[int, bool]()
		for i := 0; i < sequenceSize; i++ {
			tree.Insert(i, true)
		}
		if tree.Size() != sequenceSize {
			t.Fatalf("Expected size %d, got %d", sequenceSize, tree.Size())
		}
		for i := 0; i < sequenceSize; i++ {
			tree.Remove(i)
		}
		if tree.Size() != 0 {
			t.Errorf("Tree should be empty, but size is %d", tree.Size())
		}
	})

	t.Run("AscInsert_DescRemove", func(t *testing.T) {
		tree := rbtree.New[int, bool]()
		for i := 0; i < sequenceSize; i++ {
			tree.Insert(i, true)
		}
		if tree.Size() != sequenceSize {
			t.Fatalf("Expected size %d, got %d", sequenceSize, tree.Size())
		}
		for i := sequenceSize - 1; i >= 0; i-- {
			tree.Remove(i)
		}
		if tree.Size() != 0 {
			t.Errorf("Tree should be empty, but size is %d", tree.Size())
		}
	})
}

func TestSustainedChurnAndIntermittentVerification(t *testing.T) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	t.Logf("Running sustained churn test with seed: %d", seed)

	const initialPopulation = 1000
	const churnCycles = 20000
	const checkInterval = 1000

	tree := rbtree.New[int, bool]()
	trackingMap := make(map[int]struct{})

	for i := 0; i < initialPopulation; i++ {
		key := rng.Intn(initialPopulation * 5)
		if _, exists := trackingMap[key]; exists {
			i--
			continue
		}
		tree.Insert(key, true)
		trackingMap[key] = struct{}{}
	}

	for i := 0; i < churnCycles; i++ {
		operation := rng.Intn(100)
		if operation < 50 {
			key := rng.Intn(initialPopulation * 5)
			tree.Insert(key, true)
			trackingMap[key] = struct{}{}
		} else { // 50% шанс на удаление
			if len(trackingMap) == 0 {
				continue
			}
			keysInMap := getKeysFromMapForTest(trackingMap)
			keyToRemove := keysInMap[rng.Intn(len(keysInMap))]
			tree.Remove(keyToRemove)
			delete(trackingMap, keyToRemove)
		}

		if i > 0 && i%checkInterval == 0 {
			if tree.Size() != len(trackingMap) {
				t.Fatalf("Size mismatch at cycle %d. Expected %d, got %d", i, len(trackingMap), tree.Size())
			}
			verifyOrder(t, tree, trackingMap)
		}
	}

	t.Run("FinalVerification", func(t *testing.T) {
		if tree.Size() != len(trackingMap) {
			t.Errorf("Final size is incorrect. Expected %d, got %d", len(trackingMap), tree.Size())
		}
		verifyOrder(t, tree, trackingMap)
	})
}

func verifyOrder(t *testing.T, tree *rbtree.RBTreeMap[int, bool], m map[int]struct{}) {
	t.Helper()
	expectedOrder := getKeysFromMapForTest(m)
	slices.Sort(expectedOrder)

	actualOrder := make([]int, 0, tree.Size())
	for k, _ := range tree.InOrder() {
		actualOrder = append(actualOrder, k)
	}

	if !slices.Equal(expectedOrder, actualOrder) {
		limit := 20
		if len(expectedOrder) < limit {
			limit = len(expectedOrder)
		}
		t.Fatalf("Order verification failed.\nExpected head: %v\nGot head:      %v", expectedOrder[:limit], actualOrder[:limit])
	}
}

func getKeysFromMapForTest(m map[int]struct{}) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
