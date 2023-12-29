package tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func Test_tree_delete_randomised(t1 *testing.T) {
	randomeInsert := 30
	// make sure that you will set a small value in case of printing
	printOnTerminal := false
	for degree := 3; degree < 100; degree++ {
		for numKeys := 1; numKeys < 1000; numKeys += degree {
			tempValues := createRandArray(numKeys)
			for i := 0; i < randomeInsert; i++ {
				insert := getRandomSlice(tempValues)
				deleteElm := getRandomSlice(tempValues)

				if printOnTerminal {
					fmt.Println(insert, deleteElm)
					fmt.Printf("insert : %d delete : %d\n", len(insert), len(deleteElm))
					fmt.Println()
				}

				t := newTree(insert, degree, printOnTerminal)
				deleteItems(t, deleteElm, false)

				if printOnTerminal {
					fmt.Printf("------------   done : %d  ------------------\n\n", i+1)
				}
			}
		}
	}
	fmt.Println("done mission completed")
}

func Test_tree_Put(t1 *testing.T) {
	t := tree[int, int]{}
	keys := []int{5, 15, 25, 35, 45}
	for i, key := range keys {
		t.Put(key, i)
	}
	t = tree[int, int]{}
	keys = []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for i, key := range keys {
		t.Put(key, i)
	}
	fmt.Println(t.root)
}

func Test_tree_Delete(t1 *testing.T) {
	t := tree[int, int]{}
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
	degree := 3
	insert := "20 50 24 5 42"
	delete := "5 50 20 42 24"
	insrt := stringToSlice(insert)
	t := newTree(insrt, degree, true)
	fmt.Printf("minChilds : %d\n\n", t.getMinKeys())
	fmt.Printf("minChilds : %d\n\n", t.getMaxKeys())
	eld := stringToSlice(delete)
	fmt.Printf("insert : %d delete : %d\n", len(insrt), len(eld))
	fmt.Println(eld)
	deleteItems(t, eld, true)
}

func createRandArray(size int) []int {
	arr := make([]int, size)
	m := make(map[int]struct{})
	for len(m) != size {
		key := rand.Intn(size*10) + 1
		m[key] = struct{}{}
	}

	i := 0
	for key := range m {
		arr[i] = key
		i++
	}
	return arr
}

// [20 8 9 29 16 4 79 80 93 1] [1 8 29 93 16 80 20 9 79 4]

func getRandomSlice[T1 Comparable](a []T1) []T1 {
	temp := make([]T1, len(a))
	copy(temp, a)
	var randomArray []T1
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

func newTree[T1 Comparable](insertElm []T1, degree int, printOnTerminal bool) *tree[T1, int] {
	// forcefully value is int
	t := &tree[T1, int]{
		degree: degree,
	}
	for i, key := range insertElm {
		t.Put(key, i)
	}
	if printOnTerminal {
		display(t)
		fmt.Printf("---------------tree---------------\n\n")
	}
	return t
}

func deleteItems[T1 Comparable, T2 any](t *tree[T1, T2], deleteElm []T1, dispaly bool) {
	for _, key := range deleteElm {
		t.Delete(key)
		if dispaly {
			fmt.Println("---> ", key)
			display(t)
			if t.root != nil {
				fmt.Printf("------------------------------\n")
			}
		}
		_, err := t.Get(key)
		if err == nil {
			// this show that key is present
			panic("key is not deleted")
		}

	}
	if t.root != nil {
		panic("root is not empty")
	}
}
