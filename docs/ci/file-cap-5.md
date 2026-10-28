# 500-line file cap census and gate proof — issue #5

Refs #5, Part of #1. Census on `2b776c1` (main after #10 and #11):

| File | Lines |
|---|---|
| `internal/intake/doc.go` | 7 |
| `internal/intake/intake.go` | 116 |
| `internal/intake/intake_test.go` | 159 |

No file exceeds 500 lines, so nothing was split.

Gate proof: a scratch 501-line file (`const` block, otherwise
vet/lint/gofumpt-clean — vet, `golangci-lint run`, and `fmt --diff`
all still passed with it present) was staged so the lane's
`git ls-files '*.go'` census could see it, exactly as a committed PR
file. The cap step from the lint lane then reported:

`internal/intake/seedfault_501.go (501 lines)` and failed (RED).

After removing the scratch file, the same step reports
`500-line cap: OK (3 files checked)` (GREEN). The scratch file was
never committed; the tree is clean.
