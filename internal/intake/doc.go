// SPDX-License-Identifier: MIT
// Copyright (c) 2026 eshu-hq

// Package intake holds the ingester's intake orchestration: validating
// intake items and accumulating them into batches that flush when they
// reach a configured size or age.
package intake
