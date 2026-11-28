package retention

import "math"

// PendingLimitError corresponds to Python's ToolBufferLimitError. The
// transport-level failure metadata is outside this accounting primitive.
type PendingLimitError struct{}

func (*PendingLimitError) Error() string {
	return "Buffered response exceeded the pending output limit."
}

// PendingBudget accounts for source tool-buffer bytes and records. Unlike the
// delivered-evidence Budget, failed additions leave the counters unchanged.
// Counters use machine integers; valid callers must keep values representable.
type PendingBudget struct {
	maxBytes, maxRecords int
	bytes, records       int
}

// NewPendingBudget creates a budget with explicit limits, including zero.
func NewPendingBudget(maxBytes, maxRecords int) *PendingBudget {
	return &PendingBudget{maxBytes: maxBytes, maxRecords: maxRecords}
}

// Add charges bytes and records atomically, rejecting an addition above either
// limit. The overflow check prevents large positive charges from wrapping.
func (b *PendingBudget) Add(size, records int) error {
	if (size > 0 && b.bytes > math.MaxInt-size) ||
		(records > 0 && b.records > math.MaxInt-records) ||
		(size < 0 && b.bytes < math.MinInt-size) ||
		(records < 0 && b.records < math.MinInt-records) {
		return &PendingLimitError{}
	}
	if b.bytes+size > b.maxBytes || b.records+records > b.maxRecords {
		return &PendingLimitError{}
	}
	b.bytes += size
	b.records += records
	return nil
}

// Release subtracts charges after retained output is discarded or published.
// Like Python's release, it does not validate that the resulting tally is
// non-negative. Callers must keep the resulting machine integers representable.
func (b *PendingBudget) Release(size, records int) {
	b.bytes -= size
	b.records -= records
}

// Bytes and Records report currently charged quantities.
func (b *PendingBudget) Bytes() int   { return b.bytes }
func (b *PendingBudget) Records() int { return b.records }
