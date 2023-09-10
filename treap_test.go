package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicFunctionality(t *testing.T) {
	for _, nVals := range []int{1, 5, 10, 100, 1000} {
		t.Run(fmt.Sprintf("basic functionality with %d values", nVals), func(t *testing.T) {
			tc := testCase{nVals: nVals}
			tc.run(t)
		})
	}
}

func TestEmptyTreap(t *testing.T) {
	empty := NewTreap[*val]()

	require.Nil(t, empty.Min())
	require.Nil(t, empty.Max())
	require.Nil(t, empty.LeastGTE(&val{i: 0}))
	require.Nil(t, empty.GreatestLTE(&val{i: 0}))
}

// Helpers below

type testCase struct {
	nVals int
	// defaults to 10000 - range will be [-rangeVals, rangeVals] both inclusive
	rangeVals int
}

func (tc *testCase) applyDefaults() {
	if tc.rangeVals <= 0 {
		tc.rangeVals = 10000
	}
}

func (tc *testCase) run(t *testing.T) {
	tc.applyDefaults()

	require.Greater(t, tc.nVals, 0)

	treap := NewTreap[*val]()

	var m, M int
	vals := make([]int, tc.nVals)
	for i := 0; i < tc.nVals; i++ {
		v := rand.Intn(2*tc.rangeVals+1) - tc.rangeVals

		vals[i] = v
		node, isNewNode := treap.Insert(&val{
			i: v,
		})
		require.NotNil(t, node)
		require.Equal(t, v, node.Value.i)
		require.Equal(t, node.Value.Counter() == 1, isNewNode)

		if i == 0 {
			m = v
			M = v
		} else {
			m = min(m, v)
			M = max(M, v)
		}

		mi, ma := assertTreapWellFormed(t, treap)
		require.Equal(t, m, mi)
		require.Equal(t, M, ma)
	}

	sort.Ints(vals)

	require.Equal(t, vals[0], treap.Min().Value.i)
	require.Equal(t, vals[tc.nVals-1], treap.Max().Value.i)

	// successors
	current := treap.Min()
	for i := 0; i < tc.nVals; {
		j := 1
		for ; i+j < tc.nVals && vals[i+j] == vals[i]; j++ {
		}

		require.Equal(t, vals[i], current.Value.i)
		require.Equal(t, j, current.Value.Counter())
		i += j

		current = current.Successor()
	}
	require.Nil(t, current)

	// predecessors
	current = treap.Max()
	for i := tc.nVals - 1; i >= 0; {
		j := 1
		for ; i-j >= 0 && vals[i-j] == vals[i]; j++ {
		}

		require.Equal(t, vals[i], current.Value.i)
		require.Equal(t, j, current.Value.Counter())
		i -= j

		current = current.Predecessor()
	}
	require.Nil(t, current)

	// LeastGTE/GreatestLTE
	var previous *TreapNode[*val]
	for i := 0; i < tc.nVals; {
		j := 1
		for ; i+j < tc.nVals && vals[i+j] == vals[i]; j++ {
		}

		v := vals[i]
		current := treap.LeastGTE(&val{i: v})
		require.Equal(t, v, current.Value.i, i)

		if i == 0 || vals[i-1] != v-1 {
			middle := &val{i: v - 1}
			require.Equal(t, previous, treap.GreatestLTE(middle))
			require.Equal(t, current, treap.LeastGTE(middle))

			n1, n2 := treap.Neighbors(middle)
			require.Equal(t, previous, n1)
			require.Equal(t, current, n2)
		}

		n1, n2 := treap.Neighbors(&val{i: v})
		assert.Equal(t, current, n1)
		assert.Equal(t, current, n2)

		i += j
		previous = current
	}
	biggerThanAll := &val{i: vals[tc.nVals-1] + 1}
	require.Equal(t, previous, treap.GreatestLTE(biggerThanAll))
	require.Nil(t, treap.LeastGTE(biggerThanAll))
}

type val struct {
	i       int
	counter int
}

var _ TreapValue[*val] = &val{}

func (v *val) Compare(other *val) int {
	return v.i - other.i
}

func (v *val) Merge(other *val) {
	v.counter = v.Counter() + other.Counter()
}

func (v *val) Counter() int {
	if v.counter == 0 {
		v.counter = 1
	}
	return v.counter
}

// also returns min and max found in the treap
func assertTreapWellFormed(t *testing.T, treap *Treap[*val]) (int, int) {
	if treap.root != nil {
		assertHeap(t, treap.root)
		return assertBST(t, treap.root)
	}
	return 0, 0
}

func assertHeap(t *testing.T, node *TreapNode[*val]) {
	if node.left != nil {
		require.GreaterOrEqual(t, node.priority, node.left.priority)
		assertHeap(t, node.left)
	}
	if node.right != nil {
		require.GreaterOrEqual(t, node.priority, node.right.priority)
		assertHeap(t, node.right)
	}
}

// the 2 returned values are the min and the max of the subtree rooted at node
func assertBST(t *testing.T, node *TreapNode[*val]) (m, M int) {
	m = node.Value.i
	M = node.Value.i

	if node.left != nil {
		leftMin, leftMax := assertBST(t, node.left)
		require.Greater(t, node.Value.i, leftMax)
		require.LessOrEqual(t, leftMin, leftMax)
		m = leftMin
	}

	if node.right != nil {
		rightMin, rightMax := assertBST(t, node.right)
		require.Less(t, node.Value.i, rightMin)
		require.LessOrEqual(t, rightMin, rightMax)
		M = rightMax
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
