package intake

import (
	"errors"
	"time"
)

// Sentinel errors returned by item validation and batch operations.
var (
	// ErrEmptyID reports an intake item with no identifier.
	ErrEmptyID = errors.New("intake: empty item ID")
	// ErrEmptyBody reports an intake item with no payload.
	ErrEmptyBody = errors.New("intake: empty item body")
	// ErrInvalidConfig reports a Batcher configuration that cannot admit
	// or age out any item.
	ErrInvalidConfig = errors.New("intake: invalid batcher config")
	// ErrBatchFull reports an Add against a batcher already holding the
	// configured maximum number of items.
	ErrBatchFull = errors.New("intake: batch full")
)

// Item is a single unit of intake: an opaque payload keyed by an ID that
// is unique within its batch.
type Item struct {
	ID   string
	Body []byte
}

// Validate reports whether the item carries the fields a batch requires.
func (it Item) Validate() error {
	if it.ID == "" {
		return ErrEmptyID
	}

	if len(it.Body) == 0 {
		return ErrEmptyBody
	}

	return nil
}

// Batcher accumulates validated items and reports when the batch is due
// to flush: once it holds maxItems items, or once the oldest item reaches
// maxAge. The clock is injected so flush timing stays deterministic in
// tests. A Batcher is not safe for concurrent use.
type Batcher struct {
	maxItems int
	maxAge   time.Duration
	now      func() time.Time
	items    []Item
	oldest   time.Time
}

// NewBatcher returns a Batcher that flushes after maxItems items or when
// the oldest item reaches maxAge, whichever comes first.
func NewBatcher(maxItems int, maxAge time.Duration, now func() time.Time) (*Batcher, error) {
	if maxItems <= 0 || maxAge <= 0 || now == nil {
		return nil, ErrInvalidConfig
	}

	return &Batcher{maxItems: maxItems, maxAge: maxAge, now: now}, nil
}

// Add validates the item and appends it to the pending batch. It reports
// the validation failure, or ErrBatchFull when the batch already holds
// maxItems items.
func (b *Batcher) Add(it Item) error {
	if err := it.Validate(); err != nil {
		return err
	}

	if len(b.items) >= b.maxItems {
		return ErrBatchFull
	}

	if len(b.items) == 0 {
		b.oldest = b.now()
	}

	b.items = append(b.items, it)

	return nil
}

// FlushDue reports whether the pending batch should flush now: it holds
// at least one item and is either full or aged past maxAge.
func (b *Batcher) FlushDue() bool {
	if len(b.items) == 0 {
		return false
	}

	return len(b.items) >= b.maxItems || !b.now().Before(b.oldest.Add(b.maxAge))
}

// Drain returns the pending items in arrival order and resets the batch.
func (b *Batcher) Drain() []Item {
	drained := b.items
	b.items = nil

	return drained
}

// Len reports the number of pending items.
func (b *Batcher) Len() int {
	return len(b.items)
}

// Oldest reports when the oldest pending item arrived. The second return
// is false when the batch is empty.
func (b *Batcher) Oldest() (time.Time, bool) {
	if len(b.items) == 0 {
		return time.Time{}, false
	}

	return b.oldest, true
}
