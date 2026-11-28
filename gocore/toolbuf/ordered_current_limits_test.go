package toolbuf

import (
	"errors"
	"testing"

	"github.com/NggKhaiz/FCC/gocore/openaichat"
	"github.com/NggKhaiz/FCC/gocore/retention"
)

func TestOrderedBufferDefaultLimitsAndReuse(t *testing.T) {
	buffer := NewOrderedBuffer[int, string](nil)
	if err := buffer.Append("boundary", DefaultPendingBytes); err != nil {
		t.Fatal(err)
	}
	var limit *retention.PendingLimitError
	if err := buffer.Append("overflow", 1); !errors.As(err, &limit) {
		t.Fatalf("expected byte-limit error, got %v", err)
	}
	value, ok := buffer.Next()
	if !ok || value != "boundary" {
		t.Fatalf("next = %q, %t", value, ok)
	}
	if err := buffer.Reserve(7, "retry", 1); err != nil {
		t.Fatal(err)
	}
	if err := buffer.Complete(7); err != nil {
		t.Fatal(err)
	}
	if value, ok := buffer.Next(); !ok || value != "retry" {
		t.Fatalf("next = %q, %t", value, ok)
	}
	if buffer.Pending() || buffer.budget.Bytes() != 0 || buffer.budget.Records() != 0 {
		t.Fatal("buffer did not release its charges")
	}
}

func TestOrderedBufferReturnsEmptyValuesWithoutLosingPresence(t *testing.T) {
	buffer := NewOrderedBuffer[string, *string](nil)
	if err := buffer.Append(nil, 0); err != nil {
		t.Fatal(err)
	}
	if value, ok := buffer.Next(); value != nil || !ok {
		t.Fatalf("next = %v, %t; want nil, true", value, ok)
	}
	if value, ok := buffer.Next(); value != nil || ok {
		t.Fatalf("empty next = %v, %t; want nil, false", value, ok)
	}
}

func TestOrderedBufferSharesBudgetWithIdentityResolver(t *testing.T) {
	budget := retention.NewPendingBudget(40, 3)
	resolver := openaichat.NewToolCallIdentityResolverWithBudget(budget)
	if _, _, err := resolver.ResolveCurrent(1, "abc"); err != nil {
		t.Fatal(err)
	}
	buffer := NewOrderedBuffer[int, string](budget)
	if err := buffer.Reserve(1, "payload", 5); err != nil {
		t.Fatal(err)
	}
	var limit *retention.PendingLimitError
	if err := buffer.Grow(1, 1, 0); !errors.As(err, &limit) {
		t.Fatalf("expected shared byte-limit error, got %v", err)
	}
	buffer.Clear()
	if budget.Bytes() != 35 || budget.Records() != 1 {
		t.Fatalf("clear released other owner's charges: %d/%d", budget.Bytes(), budget.Records())
	}
	if len(resolver.Snapshot()) != 1 {
		t.Fatal("clear changed resolver identity state")
	}
}

func TestOrderedBufferReleasesDrainedReferences(t *testing.T) {
	buffer := NewOrderedBuffer[int, *int](nil)
	a, b := 1, 2
	if err := buffer.Append(&a, 1); err != nil {
		t.Fatal(err)
	}
	if err := buffer.Reserve(1, &b, 1); err != nil {
		t.Fatal(err)
	}
	backing := buffer.queue[:cap(buffer.queue)]
	if value, ok := buffer.Next(); !ok || value != &a {
		t.Fatal("expected first value")
	}
	if backing[0] != nil {
		t.Fatal("drained slot retained in backing array")
	}
	buffer.Clear()
	if backing[1] != nil || len(buffer.pending) != 0 || buffer.queue != nil {
		t.Fatal("clear retained queued references")
	}
}
