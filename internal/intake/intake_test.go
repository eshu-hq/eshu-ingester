// SPDX-License-Identifier: MIT
// Copyright (c) 2026 eshu-hq

package intake_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eshu-hq/eshu-ingester/internal/intake"
)

func TestItemValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		item intake.Item
		want error
	}{
		{name: "valid", item: intake.Item{ID: "a", Body: []byte("x")}, want: nil},
		{name: "empty id", item: intake.Item{Body: []byte("x")}, want: intake.ErrEmptyID},
		{name: "empty body", item: intake.Item{ID: "a"}, want: intake.ErrEmptyBody},
		{name: "both empty reports id first", item: intake.Item{}, want: intake.ErrEmptyID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.item.Validate(); !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewBatcherRejectsBadConfig(t *testing.T) {
	t.Parallel()

	now := time.Now

	for _, tt := range []struct {
		name     string
		maxItems int
		maxAge   time.Duration
		now      func() time.Time
	}{
		{name: "zero items", maxItems: 0, maxAge: time.Second, now: now},
		{name: "negative items", maxItems: -3, maxAge: time.Second, now: now},
		{name: "zero age", maxItems: 2, maxAge: 0, now: now},
		{name: "nil clock", maxItems: 2, maxAge: time.Second, now: nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := intake.NewBatcher(tt.maxItems, tt.maxAge, tt.now); !errors.Is(err, intake.ErrInvalidConfig) {
				t.Fatalf("NewBatcher() = %v, want %v", err, intake.ErrInvalidConfig)
			}
		})
	}
}

func TestBatcherAddAndDrain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	b, err := intake.NewBatcher(2, time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewBatcher() error = %v", err)
	}

	if got := b.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}

	if _, ok := b.Oldest(); ok {
		t.Fatal("Oldest() ok = true on empty batcher, want false")
	}

	if err := b.Add(intake.Item{}); !errors.Is(err, intake.ErrEmptyID) {
		t.Fatalf("Add(invalid) = %v, want %v", err, intake.ErrEmptyID)
	}

	if err := b.Add(intake.Item{ID: "a", Body: []byte("x")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	oldest, ok := b.Oldest()
	if !ok || !oldest.Equal(now) {
		t.Fatalf("Oldest() = %v, %v; want %v, true", oldest, ok, now)
	}

	if err := b.Add(intake.Item{ID: "b", Body: []byte("y")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if err := b.Add(intake.Item{ID: "c", Body: []byte("z")}); !errors.Is(err, intake.ErrBatchFull) {
		t.Fatalf("Add() over capacity = %v, want %v", err, intake.ErrBatchFull)
	}

	drained := b.Drain()
	if len(drained) != 2 || drained[0].ID != "a" || drained[1].ID != "b" {
		t.Fatalf("Drain() = %+v, want items a and b in order", drained)
	}

	if got := b.Len(); got != 0 {
		t.Fatalf("Len() after Drain = %d, want 0", got)
	}

	if err := b.Add(intake.Item{ID: "d", Body: []byte("w")}); err != nil {
		t.Fatalf("Add() after Drain error = %v", err)
	}
}

func TestBatcherFlushDue(t *testing.T) {
	t.Parallel()

	now := time.Now()
	current := now

	clock := func() time.Time { return current }

	b, err := intake.NewBatcher(3, time.Minute, clock)
	if err != nil {
		t.Fatalf("NewBatcher() error = %v", err)
	}

	if b.FlushDue() {
		t.Fatal("FlushDue() = true on empty batcher, want false")
	}

	if err := b.Add(intake.Item{ID: "a", Body: []byte("x")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if b.FlushDue() {
		t.Fatal("FlushDue() = true for a young partial batch, want false")
	}

	current = now.Add(2 * time.Minute)
	if !b.FlushDue() {
		t.Fatal("FlushDue() = false for an aged batch, want true")
	}

	if err := b.Add(intake.Item{ID: "b", Body: []byte("y")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if err := b.Add(intake.Item{ID: "c", Body: []byte("z")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	current = now
	if !b.FlushDue() {
		t.Fatal("FlushDue() = false for a full batch, want true")
	}
}
