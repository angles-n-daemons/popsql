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
	list := data.NewSkiplist()
	start := time.Now()
	for i := 0; i < 50; i++ {
		val := rng.Intn(200)
		list.Put(val)
	}
	list.DebugPrint(5)
	for i := 0; i < 20; i++ {
		val := rng.Intn(200)
		list.Get(val)
	}
	duration := time.Since(start).Milliseconds()
	fmt.Printf("skiplist took %d ms\n", duration)
}
