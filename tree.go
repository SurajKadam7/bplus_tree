package tree

import "sort"

const (
	// degree represents max number of chids in one node
	degree int = 3 //minKey = 1, maxKey = 2
)

func getMinKeys() int {
	return (degree+1)/2 - 1
}

func getMaxKeys() int {
	return degree - 1
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

// truncate truncates this instance at index so that it contains only the
// first index items. index must be less than or equal to length.
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
	root *node
}

// 20 30 50
// 10 | 20 | 30 | 40 50

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
	nn, ii = findNode(nn, key)
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
	i := sort.Search(len(p.childs), func(i int) bool {
		ind := len(p.childs[i].keys) - 1
		if ind < 0 {
			// this means that the key is removed from this child node
			return true
		}
		return !(p.childs[i].keys[ind] < key)
	})
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

// TODO it may go into the heap memory handle pointer returns
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

	mx := getMaxKeys()
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

// balance will take a key just to know the index of the node in parent child list
// and it will balance the tree if possible
func (n *node) balance(key int) {
	if n.parent == nil {
		return
	}
	if len(n.keys) >= getMinKeys() {
		// no need to balance the node
		return
	}

	// need to find the current node child index in parent
	myNodeInd := findChildInd(n.parent, key)

	l := getLeftNode(n, myNodeInd)
	if l != nil && len(l.keys) > getMinKeys() {
		// we always borrow the last element from the left neighbour
		n.borrow(l, len(l.keys)-1, 0)
		n.parent.keys[myNodeInd-1] = n.keys[0]
		return
	}

	r := getRightNode(n, myNodeInd)
	if r != nil && len(r.keys) > getMinKeys() {
		n.borrow(r, 0, len(n.keys))
		// this is because always n.child == n.keys + 1
		r.parent.keys[myNodeInd] = r.keys[0]
		return
	}

	// merge logic
	// always merge into the left node as this will be append operation
	// for coping data

	var n1, n2 *node
	if l != nil {
		n1, n2 = l, n
	} else {
		n1, n2 = n, r
		myNodeInd += 1
	}

	n2Childs := len(n2.childs)
	// merging and truncating nodes
	n1.keys = append(n1.keys, n2.keys...)
	if n1.isLeaf {
		n1.values = append(n1.values, n2.values...)
		truncate[int](&n2.values, len(n2.values))
	} else {
		n1.childs = append(n1.childs, n2.childs...)
		truncate[*node](&n2.childs, len(n2.childs))
	}

	// re-assigning next pointer
	n1.next = n2.next
	n2.next = nil

	// remove the current node from child of parent
	removeAt[*node](&n.parent.childs, myNodeInd)
	reAssignParent(n1, n2Childs)
	// // remove the n2 node first value from the parent
	// if len(n2.keys) > 0 && n.parent.keys[myNodeInd-1] == n2.keys[0] {
	// 	// if key present it will be on n2 (child index - 1) in parent
	// 	removeAt[int](&n.parent.keys, n2.keys[0])
	// 	truncate[int](&n2.keys, len(n2.keys))
	// } else if len(n2.keys) > 0 {
	// 	panic("keys are not arranged correctly")
	// }

	n1.balance(key)
}

func delete(key int, kps []keyPosition) {
	for i := len(kps) - 1; i >= 0; i-- {
		// starting from the leaf node
		kp := kps[i]
		nn, ii := kp.n, kp.ind
		// remove keys
		removeAt[int](&nn.keys, ii)
		if nn.isLeaf {
			removeAt[int](&nn.values, ii)
		}

		if len(nn.childs) > len(nn.keys)+1 {
			// mapping of the parent element to child is always with the offset
			// of 1 there for used ii + 1
			newKey := nn.childs[ii+1].keys[0]
			if !nn.childs[ii].isLeaf {
				removeAt[int](&nn.childs[ii+1].keys, 0)
			}
			// inserted from where the key is missing in node from child keys
			insertAt[int](&nn.keys, ii, newKey)
			continue
		}
		nn.balance(key)
	}
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

	kps := findAllNodes(t.root, key)
	if len(kps) == 0 {
		return
	}

	// checking if the key is present in the tree or not
	if len(kps) == 1 {
		nn, ii := kps[0].n, kps[0].ind
		if len(nn.keys) <= ii || nn.keys[ii] != key {
			return
		}
	}

	delete(key, kps)
	// changing the tree root
	switch {
	case len(t.root.childs) == 1 && len(t.root.keys)+len(t.root.childs[0].keys) <= getMaxKeys():
		t.root.childs[0].keys = append(t.root.childs[0].keys, t.root.keys...)
		// removing child address from t.root
		cNode := removeAt[*node](&t.root.childs, 0)
		t.root = cNode
		t.root.parent = nil
	case len(t.root.keys) == 0 && len(t.root.childs) == 0:
		t.root = nil
	}
}
