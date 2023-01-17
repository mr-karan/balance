package balance_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/mr-karan/balance"
)

func BenchmarkBalance(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())

	for n := 10; n <= 100000; n *= 10 {
		b.Run("items-"+strconv.Itoa(n), func(b *testing.B) {
			bl := balance.NewBalance()
			items := generateItems(n)
			for i, w := range items {
				bl.Add(i, w)
			}

			b.ResetTimer()
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					_ = bl.Get()
				}
			})
		})
	}

}

func generateItems(n int) map[string]int {
	items := make(map[string]int)
	for i := 0; i < n; i++ {
		items["server-"+strconv.Itoa(i)] = rand.Intn(100) + 50
	}
	return items
}
