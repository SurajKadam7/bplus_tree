package tree

import (
	"fmt"
	"sort"
)

// insertAt inserts a value into the given index, pushing all subsequent values
// forward.
func insertAt[T any](s *[]T, index int, item T) {
	var zero T
	*s = append(*s, zero)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = item
}

// removeAt removes a value at a given index, pulling all subsequent values
// back.
func removeAt[T any](s *[]T, index int) T {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	var zero T
	(*s)[len(*s)-1] = zero
	*s = (*s)[:len(*s)-1]
	return item
}

// pop removes and returns the last element in the list.
func pop[T any](s *[]T) (out T) {
	index := len(*s) - 1
	out = (*s)[index]
	var zero T
	(*s)[index] = zero
	*s = (*s)[:index]
	return
}

// truncate truncates this instance at index so that it contains only the
// first index items. index must be less than or equal to length.
// index = 0 means truncate all
func truncate[T any](s *[]T, index int) {
	var toClear []T
	*s, toClear = (*s)[:index], (*s)[index:]
	var zero T
	for i := 0; i < len(toClear); i++ {
		toClear[i] = zero
	}
}

type node struct {
	isLeaf bool
	keys   []int
	values []int
	next   *node
	parent *node
	childs []*node
}

type tree struct {
	root   *node
	degree int
}

func New(degree int) *tree {
	return &tree{
		degree: degree,
	}
}

func (t *tree) getMinKeys() int { return (t.degree+1)/2 - 1 }

func (t *tree) getMaxKeys() int { return t.degree - 1 }

func (t *tree) getMinChilds() int { return (t.degree + 1) / 2 }

func (t *tree) getMaxChilds() int { return t.degree }

func (n *node) removeChild() {}

func (n *node) del(key int) {}

// find node will return first node where the key may present from given n
func findNode(n *node, key int) (*node, int) {
	for {
		ind := findItemInd(n, key)
		if n.isLeaf || (ind < len(n.keys) && n.keys[ind] == key) {
			return n, ind
		}

		n = n.childs[ind]
	}
}

func findLeafNode(n *node, key int) (*node, int) {
	nn, ii := findNode(n, key)
	if nn.isLeaf {
		return nn, ii
	}
	// finding the leaf node in second attempt as there are only two chances of having key in
	// any node
	nn, ii = findNode(nn.childs[ii+1], key)
	return nn, ii
}

type keyPosition struct {
	n   *node
	ind int
}

func findAllNodes(n *node, key int) []keyPosition {
	kp := make([]keyPosition, 0, 2)
	nn, ii := findNode(n, key)
	kp = append(kp, keyPosition{n: nn, ind: ii})
	if nn.isLeaf {
		return kp
	}
	// finding the leaf node in second attempt as there are only two chances of having key in
	// any node
	nn, ii = findNode(nn.childs[ii+1], key)
	kp = append(kp, keyPosition{n: nn, ind: ii})
	return kp
}

// need to find first grater number
func findItemInd(n *node, key int) int {
	ind := sort.Search(len(n.keys), func(i int) bool {
		return !(n.keys[i] < key)
	})
	return ind
}

// findChildInd takes two arg :
// p : parent node
// key : child value
// it will return the index of child node using its
func findChildInd(p *node, key int) int {
	isEmpty := false
	emptyInd := 0

	i := sort.Search(len(p.childs), func(i int) bool {
		ind := len(p.childs[i].keys) - 1
		if ind < 0 {
			isEmpty = true
			emptyInd = i
			// this means that the key is removed from this child node
			return true
		}
		return !(p.childs[i].keys[ind] < key)
	})

	if isEmpty {
		return emptyInd
	}

	// this case will come when we have deleted rightmost nodes rightmost element
	if len(p.childs) <= i {
		return i - 1
	}

	// this case will come when you delete nodes rightmost element
	if len(p.childs) > i && p.childs[i].keys[0] > key {
		return i - 1
	}

	return i
}

// getLeftNode will take a node and its index in parent
// and retun the left node if prensent else nil
func getLeftNode(n *node, i int) *node {
	p := n.parent
	if p == nil {
		return nil
	}
	if i == 0 {
		return nil
	}
	return n.parent.childs[i-1]
}

// getRightNode will take a node and its index in parent
// and retun the right node if prensent else nil
func getRightNode(n *node, i int) *node {
	p := n.parent
	if p == nil {
		return nil
	}
	if i+1 >= len(p.childs) {
		return nil
	}
	return n.parent.childs[i+1]
}

func getNode(isLeaf bool, size int) *node {
	n := &node{
		isLeaf: isLeaf,
		keys:   make([]int, size),
	}
	if isLeaf {
		n.values = make([]int, size)
		return n
	}
	n.childs = make([]*node, size)
	return n
}

// reAssignPraent will reassign parent for last nChilds
func reAssignParent(n *node, nChilds int) {
	for i := len(n.childs) - 1; nChilds > 0; nChilds-- {
		n.childs[i].parent = n
		i--
	}
}

// TODO it may go into the heap memory handlepointer returns
func split(n *node) (n1 *node, n2 *node) {
	l := len(n.keys)
	mid := l / 2
	n1 = n
	n2 = getNode(n.isLeaf, l-mid)
	// copy keys
	copy(n2.keys, n1.keys[mid:])
	truncate[int](&n1.keys, mid)

	if n.isLeaf {
		// copy values && adjust next pointer
		copy(n2.values, n1.values[mid:])
		truncate[int](&n1.values, mid)
		n1.next, n2.next = n2, n1.next
		return n1, n2
	}

	// copy childs
	copy(n2.childs, n1.childs[mid+1:])
	truncate[*node](&n1.childs, mid+1)
	// reassigning parent of all the n2 childs
	reAssignParent(n2, len(n2.childs))
	return n1, n2
}

func (t *tree) insert(n *node, i int, key int, value int) {

	if i < len(n.keys) && n.keys[i] == key {
		// just update the value for the key
		n.values[i] = value
		return
	}

	mx := t.getMaxKeys()
	insertAt[int](&n.keys, i, key)
	if n.isLeaf {
		insertAt[int](&n.values, i, value)
	}

	if mx >= len(n.keys) {
		return
	}
	// need to grow the parent
	// case 1 : n1 has parent
	// case 2 : n1 don't have parent
	n1, n2 := split(n)
	if n1.parent == nil {
		t.root = getNode(false, 0)
		t.root.childs = append(t.root.childs, n1)
		n1.parent = t.root
	}

	n2.parent = n1.parent
	insertInd := findChildInd(n1.parent, n1.keys[0]) + 1
	insertAt[*node](&n.parent.childs, insertInd, n2)
	// n1.parent.childs = append(n1.parent.childs, n2) : bug (insertion position is important always insert after the n1 index)

	// find position where to insert the value
	i = findItemInd(n.parent, key)
	key = n2.keys[0]
	if !n.isLeaf {
		// in case of non leaf node we don't need to copy the key
		// we can simply move it
		removeAt[int](&n2.keys, 0)
	}
	// recursion
	t.insert(n1.parent, i, key, -1)
}

// balance strategy
// if possible borrow from neighbour
//   case1 : borrow from left
//   case2 : borrow from right
// if not possible then merge two nodes

// borrow will take
// neighbour node,
// neighbour node index element to be borrow,
// where to put the the borrowed elemtn in current node
func (n *node) borrow(neighbour *node, nInd int, ind int) {
	key := removeAt[int](&neighbour.keys, nInd)
	insertAt[int](&n.keys, ind, key)
	if neighbour.isLeaf {
		value := removeAt[int](&neighbour.values, nInd)
		insertAt[int](&n.values, ind, value)
	} else {
		child := removeAt[*node](&neighbour.childs, nInd)
		insertAt[*node](&n.childs, ind, child)
	}
}

func (n *node) borrowKey(neighbour *node, nInd int, ind int) {
	key := removeAt[int](&neighbour.keys, nInd)
	insertAt[int](&n.keys, ind, key)
}

// getSmallest will return the smallest node in the give node tree
func getSmallest(n *node) int {
	tmp := n
	for !tmp.isLeaf {
		tmp = tmp.childs[0]
	}
	return tmp.keys[0]
}

// merging will only happen for the root nodes
// for deletion in the intermediate node only replacement happen
func (n *node) balanceV2(key int, minKeys int) {
	p := n.parent
	if len(n.keys) >= minKeys || p == nil {
		return
	}

	minKey := key
	if len(n.keys) > 0 {
		minKey = n.keys[0]
	}

	myNodeInd := findChildInd(n.parent, minKey)

	// BORROW LOGIC
	l := getLeftNode(n, myNodeInd)
	if l != nil && len(l.keys) > minKeys {
		if n.isLeaf {
			// direct borrow and change the parent of current node
			n.borrow(l, len(l.keys)-1, 0)
			n.parent.keys[myNodeInd-1] = n.keys[0]
		} else {
			// borrow from the parent and copy the child through borrow
			// TODO : check is this logic is correct
			n.borrowKey(n.parent, myNodeInd-1, 0)
			n.parent.borrowKey(l, len(l.keys)-1, myNodeInd-1)
			// removing the child logic
			lChild := removeAt[*node](&l.childs, len(l.childs)-1)
			lChild.parent = n
			insertAt[*node](&n.childs, 0, lChild)
		}

		return
	}

	r := getRightNode(n, myNodeInd)
	if r != nil && len(r.keys) > minKeys {
		if n.isLeaf {
			n.borrow(r, 0, len(n.keys))
			r.parent.keys[myNodeInd] = r.keys[0]
		} else {
			n.borrowKey(n.parent, myNodeInd, len(n.keys))
			n.parent.borrowKey(r, 0, myNodeInd)

			// removig child logic
			rChild := removeAt[*node](&r.childs, 0)
			rChild.parent = n
			insertAt[*node](&n.childs, len(n.childs), rChild)
		}

		return
	}

	// merging process
	var left, right *node
	// this right key declaration is just to handle zero key in the current node
	// for which balancing is happening

	rightNodeInd := myNodeInd
	left, right = l, n
	if left == nil {
		left, right = n, r
		rightNodeInd += 1
	}

	if left.isLeaf {
		left.keys = append(left.keys, right.keys...)
		left.values = append(left.values, right.values...)

		removeAt[int](&n.parent.keys, rightNodeInd-1)   // removed the key from parent
		removeAt[*node](&n.parent.childs, rightNodeInd) // removed the child from parent
		left.next = right.next                          // pointer change

		// cleanup of the right node
		right.next = nil
		truncate[int](&right.keys, 0)
		truncate[int](&right.values, 0)
		if len(left.parent.keys) == 0 {
			left.parent.balanceV2(key, minKey)
			return
		}
		left.parent.balanceV2(left.parent.keys[0], minKey)
		return
	}
	parent := left.parent
	// if len(parent.childs)+len(parent.keys)-1 > getMaxKeys() {
	parentKey := removeAt[int](&parent.keys, rightNodeInd-1)

	left.keys = append(left.keys, parentKey)
	left.keys = append(left.keys, right.keys...)
	left.childs = append(left.childs, right.childs...)

	reAssignParent(left, len(right.childs))
	truncate[int](&right.keys, 0)
	truncate[*node](&right.childs, 0)
	removeAt[*node](&parent.childs, rightNodeInd)
	right.parent = nil
	if len(parent.keys) == 0 {
		parent.balanceV2(key, minKey)
		return
	}
	parent.balanceV2(parent.keys[0], minKey)
}

func (t *tree) delete(key int) {

	nn, ii := findLeafNode(t.root, key)
	if len(nn.keys) <= ii || nn.keys[ii] != key {
		panic(fmt.Sprint("key not found : ", key))
	}

	// remove key and value
	removeAt[int](&nn.keys, ii)
	if nn.isLeaf {
		removeAt[int](&nn.values, ii)
	}
	nn.balanceV2(key, t.getMinKeys())

	nn, ii = findNode(t.root, key)
	// remove key and value
	if len(nn.keys) <= ii || nn.keys[ii] != key {
		return
	}
	nn.keys[ii] = getSmallest(nn.childs[ii+1])
}

func (t *tree) Put(key int, value int) {
	if t.root == nil {
		t.root = getNode(true, 0)
	}
	node, i := findLeafNode(t.root, key)
	t.insert(node, i, key, value)
}

// TODO : should return error along with the zero values
func (t *tree) Get(key int) int {
	if t.root == nil {
		return -1
	}

	node, ind := findLeafNode(t.root, key)
	if ind >= len(node.keys) || node.keys[ind] != key {
		// key not found
		return -1
	}

	return node.keys[ind]
}

func (t *tree) Delete(key int) {
	if t.root == nil {
		return
	}

	t.delete(key)
	if t.root.isLeaf && len(t.root.keys) == 0 {
		t.root = nil
		return
	}

	if len(t.root.childs) == 1 {
		t.root = t.root.childs[0]
		t.root.parent = nil
	}
}
