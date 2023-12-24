package tree

import (
	"fmt"
	"strings"
)

func printNode(node *node) string {
	r := make([]string, 0)
	for _, val := range node.keys {
		r = append(r, fmt.Sprintf("%v ", val))
	}
	return strings.Join(r, "")
}

func display(t *tree) {
	if t.root == nil {
		return
	}

	n := t.root
	stack := make([][]*node, 0)

	stack = append(stack, []*node{n})
	bigStack := make([][]*node, 0)

	for len(stack) > 0 {
		curr := stack[0]
		bigStack = append(bigStack, curr)
		stack = stack[1:]
		currNode := make([]*node, 0)
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
