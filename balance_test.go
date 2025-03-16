package balance_test

import (
	"errors"
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

	t.Run("adding duplicate entry", func(t *testing.T) {
		bl := balance.NewBalance()
		err := bl.Add("c", 1)
		if !errors.Is(err, nil) {
			t.Error("Wrong error received", err.Error())
		}

		err = bl.Add("c", 1)
		if !errors.Is(err, balance.ErrDuplicateID) {
			t.Error("Wrong error received", err.Error())
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

	// Test remove item.
	t.Run("remove item", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 1)
		bl.Add("b", 1)
		bl.Add("c", 1)

		err := bl.Remove("b")
		if err != nil {
			t.Error("Expected no error, got", err)
		}

		ids := bl.ItemIDs()
		expected := map[string]bool{"a": true, "c": true}
		for _, id := range ids {
			if !expected[id] {
				t.Error("Unexpected ID in list", id)
			}
		}

		// Ensure removed item isn't returned by a Get.
		for i := 0; i < 100; i++ {
			if bl.Get() == "b" {
				t.Error("Removed item 'b' still returned by Get")
			}
		}
	})

	// Test remove non-existent item.
	t.Run("remove non-existent item", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("a", 1)
		err := bl.Remove("x")
		if !errors.Is(err, balance.ErrIDNotFound) {
			t.Error("Expected ErrIDNotFound, got", err)
		}
	})

	// Test list items ids.
	t.Run("list items", func(t *testing.T) {
		bl := balance.NewBalance()
		bl.Add("x", 3)
		bl.Add("y", 2)

		ids := bl.ItemIDs()
		expected := map[string]bool{"x": true, "y": true}
		for _, id := range ids {
			if !expected[id] {
				t.Error("Unexpected ID in list", id)
			}
		}

		if len(ids) != 2 {
			t.Error("Expected 2 items, got", len(ids))
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
