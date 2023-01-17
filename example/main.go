package main

import (
	"fmt"

	balance "github.com/mr-karan/balance"
)

func main() {
	// Create a new load balancer.
	b := balance.NewBalance()

	// Add items to the load balancer.
	b.Add("a", 5)
	b.Add("b", 3)
	b.Add("c", 2)

	for i := 0; i < 10; i++ {
		item := b.Get()
		fmt.Printf("%s ", item)
	}
}
