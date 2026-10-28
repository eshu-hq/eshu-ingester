# Coverage baseline and floor — issue #4

Refs #4, Part of #1. Baseline measured on #10 (`3b7161b`, merged as
`fd8ac0f`, package `internal/intake` only):

- Mode: `-covermode=atomic`, via the race lane
  (`go test ./... -race -count=1 -covermode=atomic -coverprofile=coverage.out`).
- Total: 100.0% of statements (`go tool cover -func=coverage.out`).
- Per function: `Validate`, `NewBatcher`, `Add`, `FlushDue`, `Drain`,
  `Len`, `Oldest` — each 100%.
- `COVERAGE_FLOOR` set at `100.0` in both `.github/workflows/ci.yml`
  and `Makefile`. It may only rise.
- Finish line is 90%: already met at the baseline, so no test-writing
  follow-up issue was filed and the epic is not closed below 90%.

Proof on this branch: `make race` re-runs the atomic race tests and
fails when total coverage drops below the new floor.
