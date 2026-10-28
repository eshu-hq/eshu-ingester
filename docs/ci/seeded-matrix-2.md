# Seeded RED/GREEN matrix — issue #2 (Go landing)

Refs #2, Part of #1. Proves each named lint-lane step fails on its
seeded fault and passes once the fault is removed. All probes ran
locally with the exact lane commands from `.github/workflows/ci.yml`
(Go 1.26.10 toolchain, golangci-lint v2.12.2), against head `553f955`
before this artifact was added. Each fault was reverted after its RED
read; the tree was `git status` clean between faults.

| # | Seeded fault | Named lane step | RED | GREEN after removal |
| - | - | - | - | - |
| 1 | `fmt.Sprintf("%d", "not-a-number")` in a scratch file | `go vet ./...` | exit 1: `format %d has arg "not-a-number" of wrong type string` | exit 0 |
| 2 | unexported unused func in a scratch file | `golangci-lint run ./...` | exit 1: `func seedUnusedHelper is unused (unused)` | exit 0, `0 issues.` |
| 3 | scratch test racing two goroutines on one counter | `go test -race` | exit 1: `WARNING: DATA RACE` | exit 0, `ok ... 1.006s` |
| 4 | `x:=1` in a scratch file | `golangci-lint fmt --diff ./...` | exit 1: diff on `x:=1` | exit 0 |
| 5 | `go.mod` module path rewritten to `example.com/seed-drift` | `go mod tidy -diff` | exit 1: tidy cannot resolve the drifted module path | exit 0 after `git checkout -- go.mod` |

Whole-lane proof on the un-faulted tree (`make ci`, same toolchain):

- `go vet ./...`: clean
- `golangci-lint run ./...`: `0 issues.`
- `golangci-lint fmt --diff ./...`: clean
- `go mod tidy -diff`: clean
- 500-line cap: OK (3 files checked)
- `go test ./... -count=1`: ok
- `go test -race -covermode=atomic`: ok, total statement coverage
  100.0% (floor 0%), per-function 100% across `intake.go`
  (`Validate`, `NewBatcher`, `Add`, `FlushDue`, `Drain`, `Len`,
  `Oldest`).

The coverage baseline (100.0%) is handed to #4, which sets
`COVERAGE_FLOOR` at the measured value.
