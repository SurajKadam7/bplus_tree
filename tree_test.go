package tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func Test_tree_Put(t1 *testing.T) {
	t := tree{}
	keys := []int{5, 15, 25, 35, 45}
	for i, key := range keys {
		t.Put(key, i)
	}
	t = tree{}
	keys = []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for i, key := range keys {
		t.Put(key, i)
	}
	fmt.Println(t.root)
}

func Test_tree_Delete(t1 *testing.T) {
	t := tree{}
	keys := []int{1, 7, 20, 25, 31, 42, 17, 28, 21, 19, 10, 4}
	for i, key := range keys {
		t.Put(key, i)
	}

	display(&t)
	fmt.Println("--------------------------------------")
	delete := []int{21, 31, 20, 10, 7, 25, 42, 4, 1, 19, 17, 28}
	for _, key := range delete {
		fmt.Printf("key -> %d\n", key)
		t.Delete(key)
		display(&t)
		fmt.Println("--------------------------------------")

	}
	fmt.Println(t.root)
}

func Test_separate(t1 *testing.T) {
	insert := "190 497 69 237 482 988 60 644 609 768"
	delete := "60 988 644 482 237 768 69 609 190 497"
	insrt := stringToSlice(insert)
	t := newTree(insrt)
	eld := stringToSlice(delete)
	fmt.Println(eld)
	deleteItems(t, eld, true)
}

func createRandArray(size int) []int {
	arr := make([]int, size)
	m := make(map[int]struct{})
	for len(m) != size {
		key := rand.Intn(1000) + 1
		m[key] = struct{}{}
	}

	i := 0
	for key := range m {
		arr[i] = key
		i++
	}
	return arr
}
// [962 28 887 73 902 861 20 282 430 491 585 821 979 855 140 900 672 286 573 993] [672 430 20 900 962 821 28 902 585 573 855 73 861 887 140 491 979 993 282 286]
// [190 497 69 237 482 988 60 644 609 768] [60 988 644 482 237 768 69 609 190 497]
func Test_tree_delete(t1 *testing.T) {
	numKeys := 10
	tempValues := createRandArray(numKeys)

	for i := 0; i < 50; i++ {
		insert := getRandomSlice(tempValues)
		deleteElm := getRandomSlice(tempValues)
		fmt.Println(insert, deleteElm)
		fmt.Println()
		t := newTree(insert)
		deleteItems(t, deleteElm, false)
		fmt.Printf("------------   done  ------------------\n\n")
	}
}

func getRandomSlice(a []int) []int {
	temp := make([]int, len(a))
	copy(temp, a)
	var randomArray []int
	for i := 0; i < len(a); i++ { // Change 5 to the desired length of the new array
		randomIndex := rand.Intn(len(temp))
		randomArray = append(randomArray, temp[randomIndex])
		temp[randomIndex], temp[len(temp)-1] = temp[len(temp)-1], temp[randomIndex]
		temp = temp[:len(temp)-1]
	}
	return randomArray
}

func stringToSlice(s string) []int {
	a := strings.Split(s, " ")
	res := []int{}
	for _, key := range a {
		n, _ := strconv.Atoi(key)
		res = append(res, n)
	}
	return res
}

func newTree(insertElm []int) *tree {
	t := &tree{}
	for i, key := range insertElm {
		t.Put(key, i)
	}
	display(t)
	fmt.Printf("---------------tree---------------\n\n")
	return t
}

func deleteItems(t *tree, deleteElm []int, dispaly bool) {
	for _, key := range deleteElm {
		t.Delete(key)
		if dispaly {
			fmt.Println("---> ", key)
			display(t)
			if t.root != nil {
				fmt.Printf("------------------------------\n")
			}
		}
	}
	if t.root != nil {
		panic("root is not empty")
	}
}
