package openaichat

import (
	"testing"

	"github.com/NggKhaiz/FCC/gocore/retention"
)

func BenchmarkCurrentIdentityBudget(b *testing.B) {
	for _, withBudget := range []bool{false, true} {
		name := "unbudgeted"
		if withBudget {
			name = "budgeted"
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var budget *retention.PendingBudget
				if withBudget {
					budget = retention.NewPendingBudget(1024, 16)
				}
				resolver := NewToolCallIdentityResolverWithBudget(budget)
				for index := 0; index < 8; index++ {
					if _, _, err := resolver.ResolveCurrent(index, ""); err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}
