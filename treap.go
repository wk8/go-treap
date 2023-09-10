package main

import (
	"math/rand"
)

// TODO wkpo translate to py as well?

type TreapValue[T any] interface {
	// Compare must return 0 if equal, < 0 if other is greater, > 0 if other is lesser
	Compare(other T) int
	// Merge will be called when inserting the same element (as characterized by Compare)
	Merge(other T)
}

type Treap[T TreapValue[T]] struct {
	root *TreapNode[T]
}

type TreapNode[T TreapValue[T]] struct {
	Value T
	// random treap priority
	priority float64
	parent   *TreapNode[T]
	left     *TreapNode[T]
	right    *TreapNode[T]
}

func NewTreap[T TreapValue[T]]() *Treap[T] {
	return &Treap[T]{}
}

func newTreapNode[T TreapValue[T]](value T) *TreapNode[T] {
	return &TreapNode[T]{
		Value:    value,
		priority: rand.Float64(),
	}
}

// Insert returns the node that was inserted or amended,
// together with a bool saying whether it was a new node or an amended one
func (t *Treap[T]) Insert(value T) (*TreapNode[T], bool) {
	if t.root == nil {
		t.root = newTreapNode(value)
		return t.root, true
	}

	node, isNewNode := t.root.insert(value)
	if node.parent == nil {
		// new root
		t.root = node
	}

	return node, isNewNode
}

// LeastGTE returns the "first" node greater than or equal to the given value, ie
// the node that compares >= 0 to value and < 0 to any other node
// that compares >= 0 to value.
// May return nil.
func (t *Treap[T]) LeastGTE(value T) *TreapNode[T] {
	_, result := t.Neighbors(value)
	return result
}

// GreatestLTE returns the "first" node lesser than or equal to the given value, ie
// the node that compares <= 0 to value and > 0 to any other node
// that compares <= 0 to value.
// May return nil.
func (t *Treap[T]) GreatestLTE(value T) *TreapNode[T] {
	result, _ := t.Neighbors(value)
	return result
}

// Neighbors returns both results from GreatestLTE and LeastGTE in one operation.
func (t *Treap[T]) Neighbors(value T) (*TreapNode[T], *TreapNode[T]) {
	var lastLeftParent, lastRightParent *TreapNode[T]

	current := t.root
	for current != nil {
		comparison := current.Value.Compare(value)
		if comparison == 0 {
			return current, current
		}

		if comparison > 0 {
			lastLeftParent = current
			current = current.left
		} else {
			lastRightParent = current
			current = current.right
		}
	}

	return lastRightParent, lastLeftParent
}

func (t *Treap[T]) Min() *TreapNode[T] {
	if t.root == nil {
		return nil
	}

	n := t.root
	for ; n.left != nil; n = n.left {
	}
	return n
}

func (t *Treap[T]) Max() *TreapNode[T] {
	if t.root == nil {
		return nil
	}

	n := t.root
	for ; n.right != nil; n = n.right {
	}
	return n
}

// Predecessor may return nil if the node is the treap's minimum
// see https://en.wikipedia.org/wiki/Binary_search_tree#Successor_and_predecessor
func (n *TreapNode[T]) Predecessor() *TreapNode[T] {
	if current := n.left; current != nil {
		for ; current.right != nil; current = current.right {
		}
		return current
	}

	parent := n.parent
	for current := n; parent != nil && current == parent.left; current, parent = parent, parent.parent {
	}

	return parent
}

// Successor may return nil if the node is the treap's maximum
// see https://en.wikipedia.org/wiki/Binary_search_tree#Successor_and_predecessor
func (n *TreapNode[T]) Successor() *TreapNode[T] {
	if current := n.right; current != nil {
		for ; current.left != nil; current = current.left {
		}
		return current
	}

	parent := n.parent
	for current := n; parent != nil && current == parent.right; current, parent = parent, parent.parent {
	}

	return parent
}

func (n *TreapNode[T]) insert(value T) (*TreapNode[T], bool) {
	comparison := n.Value.Compare(value)
	if comparison == 0 {
		n.Value.Merge(value)
		return n, false
	}

	var nextNode **TreapNode[T]
	if comparison > 0 {
		nextNode = &n.left
	} else {
		nextNode = &n.right
	}

	if *nextNode == nil {
		newNode := newTreapNode(value)
		newNode.parent = n
		*nextNode = newNode

		newNode.heapify()
		return newNode, true
	}

	return (*nextNode).insert(value)
}

// rotates until the treap is a heap again wrt random priorities
func (n *TreapNode[T]) heapify() {
	for n.parent != nil && n.parent.priority < n.priority {
		parent := n.parent
		grandParent := parent.parent

		if n.parent.left == n {
			n.rotateRight()
		} else {
			n.rotateLeft()
		}

		n.parent = grandParent
		if grandParent != nil {
			if grandParent.left == parent {
				grandParent.left = n
			} else {
				grandParent.right = n
			}
		}
	}
}

func (n *TreapNode[T]) rotateRight() {
	rightChild := n.right

	n.parent.parent = n
	n.right = n.parent

	n.parent.left = rightChild
	if rightChild != nil {
		rightChild.parent = n.parent
	}
}

func (n *TreapNode[T]) rotateLeft() {
	leftChild := n.left

	n.parent.parent = n
	n.left = n.parent

	n.parent.right = leftChild
	if leftChild != nil {
		leftChild.parent = n.parent
	}
}
