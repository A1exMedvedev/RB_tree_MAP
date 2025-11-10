package tests

import (
	"math/rand"
	"rb-tree-map/internal/rbtree"
	"testing"
)

func BenchmarkBulkInsert(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = rng.Int()
	}

	tree := rbtree.New[int, int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i], i)
	}
}

func BenchmarkBulkRemove(b *testing.B) {
	rng := rand.New(rand.NewSource(2))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = rng.Int()
	}

	b.StopTimer()
	tree := rbtree.New[int, int]()
	for _, key := range keys {
		tree.Insert(key, key)
	}
	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.Remove(keys[i])
	}
}

func benchmarkChurn(N int, b *testing.B) {
	tree := rbtree.New[int, int]()
	keysInTree := make([]int, N)
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < N; i++ {
		key := rng.Int()
		tree.Insert(key, i)
		keysInTree[i] = key
	}

	keysToInsert := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		keysToInsert[i] = rng.Int()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		keyIndexToRemove := i % N
		keyToRemove := keysInTree[keyIndexToRemove]

		tree.Remove(keyToRemove)

		newKey := keysToInsert[i]
		tree.Insert(newKey, i)

		keysInTree[keyIndexToRemove] = newKey
	}
}

func BenchmarkChurn_10k(b *testing.B)  { benchmarkChurn(10000, b) }
func BenchmarkChurn_100k(b *testing.B) { benchmarkChurn(100000, b) }
func BenchmarkChurn_1m(b *testing.B)   { benchmarkChurn(1000000, b) }
