package tree

import (
	"errors"
	"sort"
)

type Comparable interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int32 | int64 | float32 | float64
}

type BPTree[T1 Comparable, T2 any] interface {
	Get(key T1) (value T2, err error)
	Put(key T1, value T2) (err error)
	Delete(key T1) (err error)
}

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

// truncate will truncate the slice till the given index
// index = 0 means truncate all
func truncate[T any](s *[]T, index int) {
	var toClear []T
	*s, toClear = (*s)[:index], (*s)[index:]
	var zero T
	for i := 0; i < len(toClear); i++ {
		toClear[i] = zero
	}
}

type node[T1 Comparable, T2 any] struct {
	isLeaf bool
	keys   []T1
	values []T2
	next   *node[T1, T2]
	parent *node[T1, T2]
	childs []*node[T1, T2]
}

type tree[T1 Comparable, T2 any] struct {
	root   *node[T1, T2]
	degree int
}

// New will take degree and return Tree
// minimum value of degree possible is 3
func New[T1 Comparable, T2 any](degree int) BPTree[T1, T2] {
	if degree < 3 {
		degree = 3
	}
	return &tree[T1, T2]{
		degree: degree,
	}
}

func (t *tree[T1, T2]) getMinKeys() int { return (t.degree+1)/2 - 1 }

func (t *tree[T1, T2]) getMaxKeys() int { return t.degree - 1 }

func (t *tree[T1, T2]) getMinChilds() int { return (t.degree + 1) / 2 }

func (t *tree[T1, T2]) getMaxChilds() int { return t.degree }

func (n *node[T1, T2]) removeChild() {}

func (n *node[T1, T2]) del(key int) {}

// findNode will return first node where the key may present from given node
func findNode[T1 Comparable, T2 any](n *node[T1, T2], key T1) (*node[T1, T2], int) {
	for {
		ind := findItemInd(n, key)
		if n.isLeaf || (ind < len(n.keys) && n.keys[ind] == key) {
			return n, ind
		}

		n = n.childs[ind]
	}
}

// findLeafNode will return the leaf node and the index where the key might present
// it will be callers responsibility to check whether the key is present on given index
// or not
func findLeafNode[T1 Comparable, T2 any](n *node[T1, T2], key T1) (*node[T1, T2], int) {
	nn, ii := findNode[T1, T2](n, key)
	if nn.isLeaf {
		return nn, ii
	}
	// finding the leaf node in second attempt as there are only two chances of having key in
	// any node
	nn, ii = findNode(nn.childs[ii+1], key)
	return nn, ii
}

type keyPosition[T1 Comparable, T2 any] struct {
	n   *node[T1, T2]
	ind int
}

// findAllNodes will return all the nodes where the given key may present
func findAllNodes[T1 Comparable, T2 any](n *node[T1, T2], key T1) []keyPosition[T1, T2] {
	kp := make([]keyPosition[T1, T2], 0, 2)
	nn, ii := findNode(n, key)
	kp = append(kp, keyPosition[T1, T2]{n: nn, ind: ii})
	if nn.isLeaf {
		return kp
	}
	// finding the leaf node in second attempt as there are only two chances of having key in
	// any node
	nn, ii = findNode(nn.childs[ii+1], key)
	kp = append(kp, keyPosition[T1, T2]{n: nn, ind: ii})
	return kp
}

// need to find first grater number
func findItemInd[T1 Comparable, T2 any](n *node[T1, T2], key T1) int {
	ind := sort.Search(len(n.keys), func(i int) bool {
		return !(n.keys[i] < key)
	})
	return ind
}

// findChildInd takes two arg parent node and any key from child node
// and it will return the index of child node
func findChildInd[T1 Comparable, T2 any](p *node[T1, T2], key T1) int {
	isEmpty := false
	emptyInd := 0

	i := sort.Search(len(p.childs), func(i int) bool {
		ind := len(p.childs[i].keys) - 1
		if ind < 0 {
			// this means that the key is removed from this child node
			isEmpty = true
			emptyInd = i
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
// and retun the left node if present else nil
func getLeftNode[T1 Comparable, T2 any](n *node[T1, T2], i int) *node[T1, T2] {
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
// and retun the right node if present else nil
func getRightNode[T1 Comparable, T2 any](n *node[T1, T2], i int) *node[T1, T2] {
	p := n.parent
	if p == nil {
		return nil
	}
	if i+1 >= len(p.childs) {
		return nil
	}
	return n.parent.childs[i+1]
}

// getNode will take two ars and return new node
func getNode[T1 Comparable, T2 any](isLeaf bool, size int) *node[T1, T2] {
	n := &node[T1, T2]{
		isLeaf: isLeaf,
		keys:   make([]T1, size),
	}
	if isLeaf {
		n.values = make([]T2, size)
		return n
	}
	n.childs = make([]*node[T1, T2], size)
	return n
}

// reAssignPraent will reassign parent for last nChilds
func reAssignParent[T1 Comparable, T2 any](n *node[T1, T2], nChilds int) {
	for i := len(n.childs) - 1; nChilds > 0; nChilds-- {
		n.childs[i].parent = n
		i--
	}
}

// TODO it may go into the heap memory handlepointer returns
func split[T1 Comparable, T2 any](n *node[T1, T2]) (n1 *node[T1, T2], n2 *node[T1, T2]) {
	l := len(n.keys)
	mid := l / 2
	n1 = n
	n2 = getNode[T1, T2](n.isLeaf, l-mid)
	// copy keys
	copy(n2.keys, n1.keys[mid:])
	truncate(&n1.keys, mid)

	if n.isLeaf {
		// copy values && adjust next pointer
		copy(n2.values, n1.values[mid:])
		truncate(&n1.values, mid)
		n1.next, n2.next = n2, n1.next
		return n1, n2
	}

	// copy childs
	copy(n2.childs, n1.childs[mid+1:])
	truncate(&n1.childs, mid+1)
	// reassigning parent of all the n2 childs
	reAssignParent(n2, len(n2.childs))
	return n1, n2
}

func (t *tree[T1, T2]) insert(n *node[T1, T2], i int, key T1, value T2) {

	if i < len(n.keys) && n.keys[i] == key {
		// just update the value for the key
		n.values[i] = value
		return
	}

	mx := t.getMaxKeys()
	insertAt(&n.keys, i, key)
	if n.isLeaf {
		insertAt(&n.values, i, value)
	}

	if mx >= len(n.keys) {
		return
	}
	// need to grow the parent
	// case 1 : n1 has parent
	// case 2 : n1 don't have parent
	n1, n2 := split(n)
	if n1.parent == nil {
		t.root = getNode[T1, T2](false, 0)
		t.root.childs = append(t.root.childs, n1)
		n1.parent = t.root
	}

	n2.parent = n1.parent
	insertInd := findChildInd(n1.parent, n1.keys[0]) + 1
	insertAt(&n.parent.childs, insertInd, n2)

	// find position where to insert the value
	i = findItemInd(n.parent, key)
	key = n2.keys[0]
	if !n.isLeaf {
		// in case of non leaf node we don't need to copy the key
		// we can simply move it
		removeAt(&n2.keys, 0)
	}

	// value will be only inserted if the node is leaf node for non-leaf node only key
	// get inserted
	t.insert(n1.parent, i, key, value)
}

// borrow will take
// neighbour node,
// neighbour node index element to be borrow,
// where to put the the borrowed element in current node
func (n *node[T1, T2]) borrow(neighbour *node[T1, T2], nInd int, ind int) {
	key := removeAt(&neighbour.keys, nInd)
	insertAt(&n.keys, ind, key)
	if neighbour.isLeaf {
		value := removeAt(&neighbour.values, nInd)
		insertAt(&n.values, ind, value)
	} else {
		child := removeAt(&neighbour.childs, nInd)
		insertAt(&n.childs, ind, child)
	}
}

// borrowKey is same like borrow but it just borrow the key
func (n *node[T1, T2]) borrowKey(neighbour *node[T1, T2], nInd int, ind int) {
	key := removeAt(&neighbour.keys, nInd)
	insertAt(&n.keys, ind, key)
}

// getSmallest will return the smallest node in the give node tree
func getSmallest[T1 Comparable, T2 any](n *node[T1, T2]) T1 {
	tmp := n
	for !tmp.isLeaf {
		tmp = tmp.childs[0]
	}
	return tmp.keys[0]
}

// balance strategy
// case 1 :  if possible borrow from neighbour
// 	caseA : borrow from left
//	   caseB : borrow from right
// case 2 : if not possible then merge two nodes

func (n *node[T1, T2]) balanceV2(key T1, minKeys int) {
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
			n.borrowKey(n.parent, myNodeInd-1, 0)
			n.parent.borrowKey(l, len(l.keys)-1, myNodeInd-1)
			// removing the child logic
			lChild := removeAt(&l.childs, len(l.childs)-1)
			lChild.parent = n
			insertAt(&n.childs, 0, lChild)
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
			rChild := removeAt(&r.childs, 0)
			rChild.parent = n
			insertAt(&n.childs, len(n.childs), rChild)
		}

		return
	}

	// MERGE LOGIC
	var left, right *node[T1, T2]

	rightNodeInd := myNodeInd
	left, right = l, n
	if left == nil {
		left, right = n, r
		rightNodeInd += 1
	}

	if left.isLeaf {
		left.keys = append(left.keys, right.keys...)
		left.values = append(left.values, right.values...)

		removeAt(&n.parent.keys, rightNodeInd-1) // removed the key from parent
		removeAt(&n.parent.childs, rightNodeInd) // removed the child from parent
		left.next = right.next                   // pointer change

		// cleanup of the right node
		right.next = nil
		truncate(&right.keys, 0)
		truncate(&right.values, 0)
		if len(left.parent.keys) == 0 {
			left.parent.balanceV2(key, minKeys)
			return
		}
		left.parent.balanceV2(left.parent.keys[0], minKeys)
		return
	}

	parent := left.parent
	parentKey := removeAt(&parent.keys, rightNodeInd-1)

	left.keys = append(left.keys, parentKey)
	left.keys = append(left.keys, right.keys...)
	left.childs = append(left.childs, right.childs...)

	reAssignParent(left, len(right.childs))
	truncate(&right.keys, 0)
	truncate(&right.childs, 0)
	removeAt(&parent.childs, rightNodeInd)
	right.parent = nil
	if len(parent.keys) == 0 {
		parent.balanceV2(key, minKeys)
		return
	}
	parent.balanceV2(parent.keys[0], minKeys)
}

func (t *tree[T1, T2]) delete(key T1) (err error) {

	nn, ii := findLeafNode(t.root, key)
	if len(nn.keys) <= ii || nn.keys[ii] != key {
		return errors.New("key is not present")
		// BELOW LINE IS FOR TESTING
		// panic(fmt.Sprint("key not found : ", key))
	}

	// remove key and value
	removeAt(&nn.keys, ii)
	if nn.isLeaf {
		removeAt(&nn.values, ii)
	}
	nn.balanceV2(key, t.getMinKeys())

	nn, ii = findNode[T1, T2](t.root, key)
	// remove key and value
	if len(nn.keys) <= ii || nn.keys[ii] != key {
		return
	}
	nn.keys[ii] = getSmallest(nn.childs[ii+1])
	return nil
}

// Put add given key and value into the tree and return error if any
func (t *tree[T1, T2]) Put(key T1, value T2) (err error) {
	if t.root == nil {
		t.root = getNode[T1, T2](true, 0)
	}
	node, i := findLeafNode(t.root, key)
	t.insert(node, i, key, value)
	return err
}

// Get will take a key and will return value and error if any
func (t *tree[T1, T2]) Get(key T1) (value T2, err error) {

	if t.root == nil {
		return value, errors.New("root is empty")
	}

	node, ind := findLeafNode[T1, T2](t.root, key)
	if ind >= len(node.keys) || node.keys[ind] != key {
		// key not found
		return value, errors.New("key is not present")
	}

	return node.values[ind], err
}

// Delete will take key and return error if key is not present
func (t *tree[T1, T2]) Delete(key T1) (err error) {
	if t.root == nil {
		return
	}

	err = t.delete(key)
	if err != nil {
		return err
	}

	if t.root.isLeaf && len(t.root.keys) == 0 {
		t.root = nil
		return
	}

	if len(t.root.childs) == 1 {
		t.root = t.root.childs[0]
		t.root.parent = nil
	}
	return err
}
