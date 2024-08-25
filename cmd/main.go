package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/angles-n-daemons/popsql/internal/data"
)

func main() {
	// test skiplist first
	rng := rand.New(rand.NewSource(1))
	list := data.NewSkipList()
	start := time.Now()
	for i := 0; i < 100000; i++ {
		val := rng.Intn(100000)
		list.Insert(val)
	}
	for i := 0; i < 100000; i++ {
		val := rng.Intn(100000)
		list.Search(val)
	}
	duration := time.Since(start).Milliseconds()
	fmt.Printf("skiplist took %d ms\n", duration)

	// repeat test with linked list
	rng = rand.New(rand.NewSource(1))
	llist := data.LinkedList{}
	start = time.Now()
	for i := 0; i < 100000; i++ {
		val := rng.Intn(100000)
		llist.Insert(val)
	}
	for i := 0; i < 100000; i++ {
		val := rng.Intn(100000)
		llist.Search(val)
	}
	duration = time.Since(start).Milliseconds()
	fmt.Printf("linkedlist took %d ms\n", duration)
}
