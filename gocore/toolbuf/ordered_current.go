package toolbuf

import (
	"errors"

	"github.com/NggKhaiz/FCC/gocore/retention"
)

// ErrDuplicateKey and ErrMissingKey identify source-order buffer misuse. The
// Python equivalent raises ValueError and KeyError, respectively.
var (
	ErrDuplicateKey = errors.New("Tool slot is already reserved")
	ErrMissingKey   = errors.New("tool slot is not pending")
)

const (
	// Current Python defaults from core/tool_call_buffer.py.
	DefaultPendingBytes   = 16 * 1024 * 1024
	DefaultPendingRecords = 65_536
)

type orderedSlot[T any] struct {
	value     T
	ready     bool
	discarded bool
	size      int
	records   int
}

// OrderedBuffer holds retained source outputs until their predecessors are
// ready. It is request-local and is not safe for concurrent use, like the
// Python reference. Keys must be comparable and charges representable as int.
type OrderedBuffer[K comparable, T any] struct {
	budget  *retention.PendingBudget
	queue   []*orderedSlot[T]
	pending map[K]*orderedSlot[T]
}

// NewOrderedBuffer shares the supplied budget, or creates one with Python's
// default limits. A shared budget may also be used by the source assembler.
func NewOrderedBuffer[K comparable, T any](budget *retention.PendingBudget) *OrderedBuffer[K, T] {
	if budget == nil {
		budget = retention.NewPendingBudget(DefaultPendingBytes, DefaultPendingRecords)
	}
	return &OrderedBuffer[K, T]{budget: budget, pending: make(map[K]*orderedSlot[T])}
}

// Pending reports whether any reserved slot remains unresolved.
func (b *OrderedBuffer[K, T]) Pending() bool { return len(b.pending) != 0 }

// Reserve appends an unresolved slot; failed charges do not alter the queue.
func (b *OrderedBuffer[K, T]) Reserve(key K, value T, size int) error {
	if _, exists := b.pending[key]; exists {
		return ErrDuplicateKey
	}
	if err := b.budget.Add(size, 1); err != nil {
		return err
	}
	slot := &orderedSlot[T]{value: value, size: size, records: 1}
	b.pending[key] = slot
	b.queue = append(b.queue, slot)
	return nil
}

// Append adds a ready value behind any previously reserved slots.
func (b *OrderedBuffer[K, T]) Append(value T, size int) error {
	if err := b.budget.Add(size, 1); err != nil {
		return err
	}
	b.queue = append(b.queue, &orderedSlot[T]{value: value, ready: true, size: size, records: 1})
	return nil
}

// Grow adds retention charges to an unresolved slot after a successful budget
// check. The missing-key check precedes charging, matching Python.
func (b *OrderedBuffer[K, T]) Grow(key K, size, records int) error {
	slot, exists := b.pending[key]
	if !exists {
		return ErrMissingKey
	}
	if err := b.budget.Add(size, records); err != nil {
		return err
	}
	slot.size += size
	slot.records += records
	return nil
}

// Value returns the unresolved slot's value without changing its state.
func (b *OrderedBuffer[K, T]) Value(key K) (T, error) {
	if slot, exists := b.pending[key]; exists {
		return slot.value, nil
	}
	var zero T
	return zero, ErrMissingKey
}

// Complete makes a reserved slot eligible for ordered release.
func (b *OrderedBuffer[K, T]) Complete(key K) error {
	slot, exists := b.pending[key]
	if !exists {
		return ErrMissingKey
	}
	delete(b.pending, key)
	slot.ready = true
	return nil
}

// Discard makes an unresolved slot ready without publishing its value.
func (b *OrderedBuffer[K, T]) Discard(key K) error {
	slot, exists := b.pending[key]
	if !exists {
		return ErrMissingKey
	}
	delete(b.pending, key)
	slot.ready, slot.discarded = true, true
	return nil
}

// DiscardPending drops only unresolved slots, leaving completed values intact.
func (b *OrderedBuffer[K, T]) DiscardPending() {
	for key, slot := range b.pending {
		delete(b.pending, key)
		slot.ready, slot.discarded = true, true
	}
}

// Next advances through the ready prefix, releasing each charge before the
// corresponding value is returned. Discarded slots are released but skipped.
// Like one step of Python's drain generator, it stops at an unresolved slot.
func (b *OrderedBuffer[K, T]) Next() (T, bool) {
	for len(b.queue) > 0 && b.queue[0].ready {
		slot := b.queue[0]
		b.queue[0] = nil
		b.queue = b.queue[1:]
		if len(b.queue) == 0 {
			b.queue = nil
		}
		b.budget.Release(slot.size, slot.records)
		if !slot.discarded {
			return slot.value, true
		}
	}
	var zero T
	return zero, false
}

// Clear releases only this buffer's charges, preserving other users of a
// shared budget, and drops references to every queued payload.
func (b *OrderedBuffer[K, T]) Clear() {
	for i, slot := range b.queue {
		b.budget.Release(slot.size, slot.records)
		b.queue[i] = nil
	}
	b.queue = nil
	clear(b.pending)
}
