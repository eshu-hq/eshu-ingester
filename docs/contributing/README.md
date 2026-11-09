# Contributing to eshu-ingester

The canonical guide is [CONTRIBUTING.md](../../CONTRIBUTING.md) — this
page is the quick version, and CONTRIBUTING.md wins on any disagreement.

## Setup

Go 1.26.10 (the `go` line in [go.mod](../../go.mod)). No third-party
dependencies to fetch. For a quick focused run:

```sh
go test ./... -count=1
```

## Before you push

```sh
make ci
```

That covers vet, lint, format, a tidy diff, the 500-line file cap, the
tests, and the race suite with its coverage floor.

## Ground rules

- CI job names are frozen because branch protection requires
  `ci-complete`. Docs-only PRs run no lane and stay mergeable through it.
- The coverage floor may only rise. Never lower it, exclude packages
  silently, or delete assertions to move the number.
- PR lanes never touch live systems: tests run against `testdata/`
  fixtures. Keep base images digest-pinned.
- Shell scripts pass `shellcheck`.
- The [CI notes](../ci) explain how each gate was set up and proved —
  read the relevant one before changing a gate.
