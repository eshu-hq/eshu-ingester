# Contributing

This repo is a skeleton: the tree holds only `README.md` today. The
ingester holds no cutover authority until its reviewed cutover
(Refs `eshu-hq/eshu-ingester#1`, part of `eshu-hq/eshu#6707`).

## Local gates

`make ci` runs the same `lint`, `test`, and `race` lane commands CI runs,
once Go code exists. Run it before every push. Until `go.mod` lands each
target passes vacuous; afterwards `go.mod` is the version of record and CI
installs exactly that. `make lint` also runs `actionlint` for workflow
changes.

Shell scripts pass `shellcheck` (see the `shell` lane).

## CI lanes

Job names are frozen because branch protection requires `ci-complete`:

| Lane | Runs when |
| --- | --- |
| `lint` | Go sources or module changes, workflow changes, or anything unclassified |
| `test` | Above, plus fixtures and config the tests read |
| `race` | Same inputs as `test`, plus a coverage artifact and floor check |
| `shell` | Shell script changes |
| `image` | `Dockerfile` / `.dockerignore`, Go sources, or workflow changes (build only, no push) |
| `govulncheck` | Dependency changes, pushes to `main`, weekly schedule, manual dispatch |
| `gosec` | Go, dependency, or workflow changes; pushes to `main`; weekly schedule. Advisory; SARIF to the Security tab |
| `ci-complete` | Always; fails on any needed lane failure, passes on legitimate skips |

Docs-only PRs run no lane and stay mergeable through `ci-complete`.
A manual `workflow_dispatch` run records every lane; the weekly schedule
runs `govulncheck` and `gosec` only.

## Coverage floor

The `race` job fails below `COVERAGE_FLOOR` in `ci.yml`. The seed floor is
`0`: when the first Go code lands, that child measures the baseline,
sets the floor there, and files the test-writing work to reach the 90%
target. The floor may only rise. Never lower it, exclude packages silently,
or delete assertions to move the number.

## Fixture-only CI

PR lanes never touch live systems: tests run against `testdata/`
fixtures. Keep base images digest-pinned when a `Dockerfile` lands.
