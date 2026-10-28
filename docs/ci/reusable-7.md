# Reusable-workflow move — issue #7

Refs #7, Part of #1. Lane bodies move to the pinned shared workflow
`eshu-hq/.github/.github/workflows/go-collector-ci.yml@v2`
(`4d01ab9`, the current v2 tip, verified against the tag); this repo's
`ci.yml` keeps the paths-filter `changes` job, the `go-ci` call with
this repo's run-gate wiring, and the two caller-side jobs the reusable
cannot own (weekly findings issue, Security-tab upload).

## What changed and what did not

- Inner lane names are unchanged (lint, test, race, shell, image,
  boundary, govulncheck, gosec, actionlint, ci-complete), but check
  runs now nest as `go-ci / <lane>` (same as the sibling
  `eshu-collector-jira` move). The repo ruleset's required context
  moves from `ci-complete` to `go-ci / ci-complete` once; nothing else
  in branch protection changes.
- Run-gate wiring preserves the in-repo semantics exactly: pushes,
  schedules, and dispatches run every lane; PRs select by path class.
  `boundary` stays skipped (empty command); `gosec` stays advisory on
  PRs and fails scheduled runs so the weekly-issue job is reachable.
- The only lane-body deltas inherited from upstream v2: govulncheck
  runs in `-scan package` mode, the 500-line cap counts via `find`
  (no git dependency), and shellcheck runs through the
  action-shellcheck action. The coverage floor input is 100.0.
- Caller-side `scheduled-findings` replaces the old
  `govulncheck-report` download-and-comment job with the org pattern
  (failure-gated issue linking the run). `gosec-sarif` keeps the
  Security-tab upload on non-PR events (artifact `gosec-ingester`).

## Proof

PR #16's CI (run 40783452540, two settled reads): `go-ci / lint`,
`test`, `race`, `image`, `gosec`, `actionlint`, and `go-ci /
ci-complete` green; `boundary`, `govulncheck` (no dep change), `shell`,
and the caller-side jobs legitimately skipped. `go-ci / race` enforced
the 100.0 floor; `go-ci / image` ran the digest-pin command. The
merge-push run 40783999113 repeated the full set green on main, which
is when the ruleset was switched. Fork PRs still run advisory-only
(unchanged caller conditions; proven by probes #14, #18, #19).

## Ruleset

Before merge the active ruleset requires `ci-complete`; after the
merge-push proves `go-ci / ci-complete` green on main, the required
context is switched to the nested name (delete + recreate the same
ruleset), verified with a second settled read.

## Follow-up fix (same issue)

The first merge-push skipped `gosec-sarif`: the move dropped the
`always()` the pre-move workflow carried, and a skipped `changes` job
(a PR-only job) skips downstream jobs whose condition does not
override it. Restored as `always() && event != PR && go-ci succeeded`,
which preserves the old behavior on every event (scheduled-findings
was never affected: `failure()` already overrides).
