# Lint backlog recount — issue #3

Refs #3, Part of #1. Recount run on `2b776c1` (main after #10 and #11)
with the committed `.golangci.yml`, uncapped:

`golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 ./...`
reports `0 issues.` (exit 0).

The config carries no baselines and no new-code-only mode: it enables
the org set (govet, staticcheck, ineffassign, unused, errcheck) with
the gofumpt formatter, and its `exclusions.rules` list is empty, so no
exclusion needs a reason or a linked issue. The two `exclusions.presets`
entries (`std-error-handling`, `common-false-positives`) are the shared
org set itself, not per-line suppressions. Nothing to fix; the backlog
is zero at the first-Go-code census.
