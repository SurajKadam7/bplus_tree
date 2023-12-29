package tree

import (
	"fmt"
	"log"
)

func main() {
	t := New[int, float64](3)

	err := t.Put(10, 20.0)
	if err != nil {
		log.Println(err)
	}

	value, err := t.Get(10)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("key : %d value : %f \n", 10, value)

	err = t.Delete(10)
	if err != nil {
		log.Println(err)
	}
}
