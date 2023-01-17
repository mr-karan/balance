package balance_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/mr-karan/balance"
)

func TestBalance(t *testing.T) {
	// Test Init.
	t.Run("init", func(t *testing.T) {
		bl := balance.NewBalance()
		if bl.Get() != "" {
			t.Error("Expected empty string")
		}
	})

	// Test round robin.
	t.Run("round robin", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 1)
		bl.Add("b", 1)
		bl.Add("c", 1)
		result := make(map[string]int)
		for i := 0; i < 999; i++ {
			result[bl.Get()]++
		}

		if result["a"] != 333 || result["b"] != 333 || result["c"] != 333 {
			t.Error("Wrong counts", result)
		}
	})

	// Test weighted.
	t.Run("weighted custom split", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 2)
		bl.Add("b", 1)
		bl.Add("c", 1)
		result := make(map[string]int)
		for i := 0; i < 1000; i++ {
			result[bl.Get()]++
		}

		if result["a"] != 500 || result["b"] != 250 || result["c"] != 250 {
			t.Error("Wrong counts", result)
		}
	})

	t.Run("weighted another custom split", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 5)
		bl.Add("b", 3)
		bl.Add("c", 2)
		result := make(map[string]int)
		for i := 0; i < 1000; i++ {
			result[bl.Get()]++
		}

		if result["a"] != 500 || result["b"] != 300 || result["c"] != 200 {
			t.Error("Wrong counts", result)
		}
	})

	// Test with one item as zero weight.
	t.Run("weighted with zero", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 0)
		bl.Add("b", 1)
		bl.Add("c", 1)
		result := make(map[string]int)
		for i := 0; i < 1000; i++ {
			result[bl.Get()]++
		}

		if result["a"] != 0 || result["b"] != 500 || result["c"] != 500 {
			t.Error("Wrong counts", result)
		}
	})
}

func TestBalance_Concurrent(t *testing.T) {
	t.Run("concurrent", func(t *testing.T) {
		var (
			a, b, c int64
		)
		bl := balance.NewBalance()
		bl.Add("a", 1)
		bl.Add("b", 1)
		bl.Add("c", 1)

		var wg sync.WaitGroup

		for i := 0; i < 999; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				switch bl.Get() {
				case "a":
					atomic.AddInt64(&a, 1)
				case "b":
					atomic.AddInt64(&b, 1)
				case "c":
					atomic.AddInt64(&c, 1)
				default:
					t.Error("Wrong item")
				}
			}()
		}

		wg.Wait()

		if a != 333 || b != 333 || c != 333 {
			t.Error("Wrong counts", a, b, c)
		}
	})
}
