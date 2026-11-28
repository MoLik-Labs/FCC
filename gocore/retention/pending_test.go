package retention

import (
	"errors"
	"math"
	"testing"
)

func TestPendingBudgetAtomicLimitsAndRelease(t *testing.T) {
	budget := NewPendingBudget(5, 2)
	if err := budget.Add(3, 1); err != nil {
		t.Fatal(err)
	}
	if err := budget.Add(2, 1); err != nil {
		t.Fatal(err)
	}
	for _, charge := range [][2]int{{1, 0}, {0, 1}, {1, 1}} {
		var limit *PendingLimitError
		if err := budget.Add(charge[0], charge[1]); !errors.As(err, &limit) {
			t.Fatalf("Add(%v): expected pending limit error, got %v", charge, err)
		}
		if budget.Bytes() != 5 || budget.Records() != 2 {
			t.Fatal("failed add changed counters")
		}
	}
	budget.Release(3, 1)
	if budget.Bytes() != 2 || budget.Records() != 1 {
		t.Fatal("release did not subtract counters")
	}
	if err := budget.Add(3, 1); err != nil {
		t.Fatal(err)
	}
}

func TestPendingBudgetRejectsUnrepresentablePositiveCharge(t *testing.T) {
	budget := NewPendingBudget(math.MaxInt, math.MaxInt)
	if err := budget.Add(math.MaxInt, 1); err != nil {
		t.Fatal(err)
	}
	if err := budget.Add(1, 0); err == nil {
		t.Fatal("overflow must not wrap")
	}
	if budget.Bytes() != math.MaxInt || budget.Records() != 1 {
		t.Fatal("failed add changed counters")
	}
}
