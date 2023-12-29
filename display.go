package tree

import (
	"fmt"
	"strings"
)

func printNode[T1 Comparable, T2 any](node *node[T1, T2]) string {
	r := make([]string, 0)
	for _, val := range node.keys {
		r = append(r, fmt.Sprintf("%v ", val))
	}
	return strings.Join(r, "")
}

func display[T1 Comparable, T2 any](t *tree[T1, T2]) {
	if t.root == nil {
		return
	}

	n := t.root
	stack := make([][]*node[T1, T2], 0)

	stack = append(stack, []*node[T1, T2]{n})
	bigStack := make([][]*node[T1, T2], 0)

	for len(stack) > 0 {
		curr := stack[0]
		bigStack = append(bigStack, curr)
		stack = stack[1:]
		currNode := make([]*node[T1, T2], 0)
		for _, s := range curr {
			if s.childs == nil {
				continue
			}
			currNode = append(currNode, s.childs...)
		}
		if len(currNode) > 0 {
			stack = append(stack, currNode)
		}
	}

	for len(bigStack) > 0 {
		curr := bigStack[0]
		bigStack = bigStack[1:]
		for i, s := range curr {
			fmt.Print(printNode(s))
			if i < len(curr)-1 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
	}
}
