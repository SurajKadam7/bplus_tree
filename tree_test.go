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
	keys := []int{5, 15, 25, 35, 45, 55, 40, 30, 20}
	for i, key := range keys {
		t.Put(key, i)
	}

	delete := []int{40, 5, 45, 35, 25, 55, 20, 30, 15}
	for _, key := range delete {
		t.Delete(key)
	}
	fmt.Println(t.root)
}

func Test_separate(t1 *testing.T) {
	insert := "40 25 20 30 35 5 15 55 45"
	delete := "15 45 5 35 40 20 25 30 55"
	insrt := stringToSlice(insert)
	t := newTree(insrt)
	eld := stringToSlice(delete)
	fmt.Println(eld)
	deleteItems(t, eld)
}

func Test_tree_delete(t1 *testing.T) {
	tempValues := []int{5, 15, 25, 35, 45, 55, 40, 30, 20}

	for i := 0; i < 10; i++ {
		insert := getRandomSlice(tempValues)
		deleteElm := getRandomSlice(tempValues)
		fmt.Println(insert, deleteElm)
		fmt.Println()
		t := newTree(insert)
		deleteItems(t, deleteElm)
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

func deleteItems(t *tree, deleteElm []int) {
	for _, key := range deleteElm {
		t.Delete(key)
		fmt.Println("---> ", key)
		display(t)
		if t.root != nil {
			fmt.Printf("------------------------------\n")
		}
	}
}
