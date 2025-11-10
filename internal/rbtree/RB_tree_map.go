package rbtree

import (
	"cmp"
	"iter"
)

type color bool

const (
	BLACK color = false
	RED   color = true
)

func less[K cmp.Ordered](a, b K) bool {
	return a < b
}

type Node[K cmp.Ordered, V any] struct {
	key    K
	value  V
	color  color
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
}

type RBTreeMap[K cmp.Ordered, V any] struct {
	root     *Node[K, V]
	sentinel *Node[K, V]
	size     int
	compare  func(a, b K) bool
}

func New[K cmp.Ordered, V any]() *RBTreeMap[K, V] {
	nilNode := &Node[K, V]{color: BLACK}
	return &RBTreeMap[K, V]{
		root:     nilNode,
		sentinel: nilNode,
		compare:  less[K],
	}
}

func NewWithCompare[K cmp.Ordered, V any](compare func(a, b K) bool) *RBTreeMap[K, V] {
	nilNode := &Node[K, V]{color: BLACK}
	return &RBTreeMap[K, V]{
		root:     nilNode,
		sentinel: nilNode,
		compare:  compare,
	}
}

func (r *RBTreeMap[K, V]) Size() int {
	return r.size
}

func (r *RBTreeMap[K, V]) search(key K) *Node[K, V] {
	current := r.root
	for current != r.sentinel {
		if key == current.key {
			return current
		}
		if r.compare(key, current.key) {
			current = current.left
		} else {
			current = current.right
		}
	}
	return r.sentinel
}

func (r *RBTreeMap[K, V]) Insert(key K, value V) {
	parent := r.sentinel
	current := r.root

	for current != r.sentinel {
		parent = current
		if key == current.key {
			current.value = value
			return
		}
		if r.compare(key, current.key) {
			current = current.left
		} else {
			current = current.right
		}
	}

	newNode := &Node[K, V]{
		key:    key,
		value:  value,
		color:  RED,
		parent: parent,
		left:   r.sentinel,
		right:  r.sentinel,
	}
	if parent == r.sentinel {
		r.root = newNode
	} else if r.compare(newNode.key, parent.key) {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	r.size++
	r.fixInsert(newNode)
}

func (r *RBTreeMap[K, V]) Get(key K) (V, bool) {
	node := r.search(key)
	if node != r.sentinel {
		return node.value, true
	}
	var zero V
	return zero, false
}

func (r *RBTreeMap[K, V]) Remove(key K) {
	z := r.search(key)
	if z == r.sentinel {
		return
	}
	r.size--

	var x *Node[K, V]
	y := z
	yOriginalColor := y.color

	if z.left == r.sentinel {
		x = z.right
		r.transplant(z, z.right)
	} else if z.right == r.sentinel {
		x = z.left
		r.transplant(z, z.left)
	} else {
		y = r.minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			r.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		r.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yOriginalColor == BLACK {
		r.fixDelete(x)
	}
}

func (r *RBTreeMap[K, V]) ContainsKey(key K) bool {
	return r.search(key) != r.sentinel
}

func (r *RBTreeMap[K, V]) InOrder() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		stack := make([]*Node[K, V], 0)
		current := r.root
		for {
			for current != r.sentinel {
				stack = append(stack, current)
				current = current.left
			}
			if len(stack) == 0 {
				return
			}
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !yield(node.key, node.value) {
				return
			}
			current = node.right
		}
	}
}

func (r *RBTreeMap[K, V]) LowerBound(key K) (K, V, bool) {
	result := r.sentinel
	current := r.root

	for current != r.sentinel {
		if !r.compare(current.key, key) {
			result = current
			current = current.left
		} else {
			current = current.right
		}
	}

	if result != r.sentinel {
		return result.key, result.value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

func (r *RBTreeMap[K, V]) UpperBound(key K) (K, V, bool) {
	result := r.sentinel
	current := r.root

	for current != r.sentinel {
		if r.compare(key, current.key) {
			result = current
			current = current.left
		} else {
			current = current.right
		}
	}

	if result != r.sentinel {
		return result.key, result.value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

func (r *RBTreeMap[K, V]) minimum(node *Node[K, V]) *Node[K, V] {
	for node.left != r.sentinel {
		node = node.left
	}
	return node
}

func (r *RBTreeMap[K, V]) fixInsert(node *Node[K, V]) {
	for node.parent.color == RED {
		if node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle.color == RED {
				node.parent.color = BLACK
				uncle.color = BLACK
				node.parent.parent.color = RED
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					r.rotateLeft(node)
				}
				node.parent.color = BLACK
				node.parent.parent.color = RED
				r.rotateRight(node.parent.parent)
			}
		} else {
			uncle := node.parent.parent.left
			if uncle.color == RED {
				node.parent.color = BLACK
				uncle.color = BLACK
				node.parent.parent.color = RED
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					r.rotateRight(node)
				}
				node.parent.color = BLACK
				node.parent.parent.color = RED
				r.rotateLeft(node.parent.parent)
			}
		}
	}
	r.root.color = BLACK
}

func (r *RBTreeMap[K, V]) rotateLeft(x *Node[K, V]) {
	y := x.right
	x.right = y.left
	if y.left != r.sentinel {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == r.sentinel {
		r.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (r *RBTreeMap[K, V]) rotateRight(y *Node[K, V]) {
	x := y.left
	y.left = x.right
	if x.right != r.sentinel {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == r.sentinel {
		r.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x
}

func (r *RBTreeMap[K, V]) transplant(u, v *Node[K, V]) {
	if u.parent == r.sentinel {
		r.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (r *RBTreeMap[K, V]) fixDelete(x *Node[K, V]) {
	for x != r.root && x.color == BLACK {
		if x == x.parent.left {
			sibling := x.parent.right
			if sibling.color == RED {
				sibling.color = BLACK
				x.parent.color = RED
				r.rotateLeft(x.parent)
				sibling = x.parent.right
			}
			if sibling.left.color == BLACK && sibling.right.color == BLACK {
				sibling.color = RED
				x = x.parent
			} else {
				if sibling.right.color == BLACK {
					sibling.left.color = BLACK
					sibling.color = RED
					r.rotateRight(sibling)
					sibling = x.parent.right
				}
				sibling.color = x.parent.color
				x.parent.color = BLACK
				sibling.right.color = BLACK
				r.rotateLeft(x.parent)
				x = r.root
			}
		} else {
			sibling := x.parent.left
			if sibling.color == RED {
				sibling.color = BLACK
				x.parent.color = RED
				r.rotateRight(x.parent)
				sibling = x.parent.left
			}
			if sibling.right.color == BLACK && sibling.left.color == BLACK {
				sibling.color = RED
				x = x.parent
			} else {
				if sibling.left.color == BLACK {
					sibling.right.color = BLACK
					sibling.color = RED
					r.rotateLeft(sibling)
					sibling = x.parent.left
				}
				sibling.color = x.parent.color
				x.parent.color = BLACK
				sibling.left.color = BLACK
				r.rotateRight(x.parent)
				x = r.root
			}
		}
	}
	x.color = BLACK
}
