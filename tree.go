package tree

import (
	"sort"
)

const (
	// degree represents max number of chids in one node
	degree int = 3 //minKey = 1, maxKey = 2
)

func getMinChilds() int {
	return (degree + 1) / 2
}

func getMaxChilds() int {
	return degree
}

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

func (n *node) removeChild() {

}

func (n *node) del(key int) {

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
func (n *node) balanceV2(key int) {
	p := n.parent
	if len(n.keys) >= getMinKeys() || p == nil {
		return
	}

	myNodeInd := findChildInd(n.parent, key)

	// BORROW LOGIC
	l := getLeftNode(n, myNodeInd)
	if l != nil && len(l.keys) > getMinKeys() {
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
	if r != nil && len(r.keys) > getMinKeys() {
		if n.isLeaf {
			n.borrow(r, 0, len(n.keys))
			r.parent.keys[myNodeInd] = r.keys[0]
		} else {
			n.borrowKey(n.parent, 0, len(n.keys))
			n.parent.borrowKey(r, 0, 0)

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
			left.parent.balanceV2(key)
			return
		}
		left.parent.balanceV2(left.parent.keys[0])
		return
	}

	c1 := 0
	c2 := 0
	parent := left.parent
	cp := make([]*node, len(parent.childs))
	copy(cp, parent.childs)
	truncate[*node](&parent.childs, 0)

	for _, nn := range cp {
		for _, key := range nn.keys {
			insertAt[int](&parent.keys, c1, key)
			c1++
		}
		truncate[int](&nn.keys, 0)
		for _, child := range nn.childs {
			insertAt[*node](&parent.childs, c2, child)
			child.parent = parent
			c2++
		}
		truncate[*node](&nn.childs, 0)
		nn.childs = nil
		nn.parent = nil
		// to keep the key in the parent
		c1++
	}
	if len(parent.keys) == 0 {
		parent.balanceV2(key)
		return
	}
	parent.balanceV2(parent.keys[0])
}

func delete(key int, kps []keyPosition) {
	leftInd := len(kps) - 1
	nn, ii := kps[leftInd].n, kps[leftInd].ind

	// remove key and value
	removeAt[int](&nn.keys, ii)
	if nn.isLeaf {
		removeAt[int](&nn.values, ii)
	}
	nn.balanceV2(key)

	if len(kps) == 1 {
		return
	}

	nn, ii = kps[0].n, kps[0].ind

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
	if t.root.isLeaf && len(t.root.keys) == 0 {
		t.root = nil
		return
	}

	if len(t.root.childs) == 1 {
		t.root = t.root.childs[0]
		t.root.parent = nil
	}
}

// func (n *node) rebalance() {

// 	// Root node has special handling.
// 	if n.parent == nil {
// 		// If root node is a branch and only has one node then collapse it.
// 		if !n.isLeaf && len(n.keys) == 1 {
// 			// Move root's child up.
// 			// find the first child
// 			cInd := findChildInd(n, n.childs[0].keys[0])
// 			child := n.childs[cInd]
// 			n.isLeaf = child.isLeaf
// 			n.keys = child.keys[:]
// 			n.childs = child.childs

// 			// Reparent all child nodes being moved.
// 			for _, inode := range n.childs {
// 				if inode.childs != nil {
// 					child.parent = n
// 				}
// 			}

// 			// Remove old child.
// 			child.parent = nil
// 		}

// 		return
// 	}

// 	// If node has no keys then just remove it.
// 	if len(n.childs) == 0 {
// 		n.parent.del(n.key)
// 		n.parent.removeChild(n)
// 		delete(n.bucket.nodes, n.pgid)
// 		n.free()
// 		n.parent.rebalance()
// 		return
// 	}

// 	// Merge with right sibling if idx == 0, otherwise left sibling.
// 	var leftNode, rightNode *node
// 	var useNextSibling = n.parent.childIndex(n) == 0
// 	if useNextSibling {
// 		leftNode = n
// 		rightNode = n.nextSibling()
// 	} else {
// 		leftNode = n.prevSibling()
// 		rightNode = n
// 	}

// 	// If both nodes are too small then merge them.
// 	// Reparent all child nodes being moved.
// 	for _, inode := range rightNode.keys {
// 		if child, ok := n.bucket.nodes[inode.Pgid()]; ok {
// 			child.parent.removeChild(child)
// 			// first it has changed its parent and then added own self as a child in parent
// 			child.parent = leftNode
// 			child.parent.children = append(child.parent.children, child)
// 		}
// 	}

// 	// Copy over keys from right node to left node and remove right node.
// 	leftNode.keys = append(leftNode.keys, rightNode.keys...)
// 	n.parent.del(rightNode.key)
// 	n.parent.removeChild(rightNode)
// 	delete(n.bucket.nodes, rightNode.pgid)

// 	// Either this node or the sibling node was deleted from the parent so rebalance it.
// 	n.parent.rebalance()
// }
