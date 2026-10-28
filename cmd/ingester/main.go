// SPDX-License-Identifier: MIT
// Copyright (c) 2026 eshu-hq

// Command ingester is the ingester service entrypoint (skeleton): it
// admits one intake item per invocation and reports the outcome as the
// process exit code. Flag parsing and long-running intake arrive with
// later issues; the batch constructor stays injectable so tests can
// drive every failure branch without I/O.
package main

import (
	"os"
	"time"

	"github.com/eshu-hq/eshu-ingester/internal/intake"
)

// batchFunc builds the batch an invocation fills once.
type batchFunc func(maxItems int) (*intake.Batcher, error)

// defaultBatch is the production constructor: single-item batches that
// age out after an hour if they are never filled.
func defaultBatch(maxItems int) (*intake.Batcher, error) {
	return intake.NewBatcher(maxItems, time.Hour, time.Now)
}

// run admits the single item ID in args and reports the outcome: 2 for
// usage errors, 1 for intake failures, 0 when the item is accepted.
func run(args []string, newBatch batchFunc) int {
	if len(args) != 1 || args[0] == "" {
		return 2
	}

	b, err := newBatch(1)
	if err != nil {
		return 1
	}

	if err := b.Add(intake.Item{ID: args[0], Body: []byte("{}")}); err != nil {
		return 1
	}

	return 0
}

// main exits with run's code. It is a thin wrapper on purpose: every
// decision lives in run, where tests cover it.
func main() {
	os.Exit(run(os.Args[1:], defaultBatch))
}
