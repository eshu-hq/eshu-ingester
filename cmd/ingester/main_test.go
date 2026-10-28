// SPDX-License-Identifier: MIT
// Copyright (c) 2026 eshu-hq

package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/eshu-hq/eshu-ingester/internal/intake"
)

func TestRunUsage(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{nil, {}, {""}, {"a", "b"}} {
		if got := run(args, defaultBatch); got != 2 {
			t.Fatalf("run(%q) = %d, want 2", args, got)
		}
	}
}

func TestRunBatchFailure(t *testing.T) {
	t.Parallel()

	broken := func(int) (*intake.Batcher, error) {
		return nil, errors.New("seed: no batcher")
	}

	if got := run([]string{"a"}, broken); got != 1 {
		t.Fatalf("run(constructor error) = %d, want 1", got)
	}

	full, err := intake.NewBatcher(1, time.Second, time.Now)
	if err != nil {
		t.Fatalf("NewBatcher() error = %v", err)
	}

	if err := full.Add(intake.Item{ID: "taken", Body: []byte("{}")}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	prefilled := func(int) (*intake.Batcher, error) { return full, nil }

	if got := run([]string{"a"}, prefilled); got != 1 {
		t.Fatalf("run(full batch) = %d, want 1", got)
	}
}

func TestRunAcceptsItem(t *testing.T) {
	t.Parallel()

	if got := run([]string{"a"}, defaultBatch); got != 0 {
		t.Fatalf("run(valid) = %d, want 0", got)
	}
}

// TestMainExitCodes covers main itself by re-executing the test binary
// with INGESTER_TEST_MAIN set: the child rewrites os.Args, calls main,
// and exits with run's code, which the parent asserts.
func TestMainExitCodes(t *testing.T) {
	if os.Getenv("INGESTER_TEST_MAIN") == "1" {
		os.Args = []string{"ingester"}
		if raw := os.Getenv("INGESTER_TEST_ARGS"); raw != "" {
			os.Args = append(os.Args, strings.Split(raw, "\x00")...)
		}

		main()
		os.Exit(99)
	}

	for _, tt := range []struct {
		name string
		args []string
		want int
	}{
		{name: "no args", args: nil, want: 2},
		{name: "accepted", args: []string{"a"}, want: 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := exec.Command(os.Args[0], "-test.run=TestMainExitCodes$")
			cmd.Env = append(os.Environ(),
				"INGESTER_TEST_MAIN=1",
				"INGESTER_TEST_ARGS="+strings.Join(tt.args, "\x00"))

			code := 0

			if err := cmd.Run(); err != nil {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("cmd.Run() error = %v", err)
				}

				code = exitErr.ExitCode()
			}

			if code != tt.want {
				t.Fatalf("main(%q) exit = %d, want %d", tt.args, code, tt.want)
			}
		})
	}
}
