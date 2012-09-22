package main

import (
	"fmt"
	"github.com/emnl/goods/linkedlist"
)

func main() {

	list := linkedlist.New()
	list.AddFirst(10)
	list.AddLast(20)

	list2 := linkedlist.New()
	list2.AddFirst(30)
	list2.AddLast(40)

	list.Conc(list2)

	fmt.Println("The list:")

	for item := range list.Iter() {
		fmt.Println(item)
	}

	fmt.Println("Contains 20?", list.Contains(20))
	fmt.Println("Is empty?", list.Empty())
	fmt.Println("First + Last =", (list.First().(int) + list.Last().(int)))

	// ...

}
