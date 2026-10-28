# Supply-chain and image activation proof — issue #6

Refs #6, Part of #1. The seed wired govulncheck, gosec SARIF,
Dependabot (gomod, github-actions), and the image lane behind tree
guards. This PR lands the first Dockerfile and Go service entrypoint,
activating them.

## What landed

- `cmd/ingester`: service entrypoint skeleton. `run(args, newBatch)`
  admits one item per invocation (exit 2 usage, 1 intake failure, 0
  accepted); `main` is a thin `os.Exit` wrapper. The constructor is
  injected so every failure branch is test-covered without I/O.
  Statement coverage holds at 100.0% in both packages (floor 100.0).
- `Dockerfile`: multi-stage build (`golang:1.26.10-alpine` build,
  `distroless/static-debian12:nonroot` runtime, `CGO_ENABLED=0`,
  `-trimpath`). Both bases digest-pinned; digests re-resolved
  2026-10-28 with `docker buildx imagetools inspect`.
- `.dockerignore`: keeps the build context to `go.mod` plus sources.

## Proof

| Claim | Command / source | Result |
|---|---|---|
| Image builds, no push | `docker build -t ci-ingester:ci .` | success, wall time ~4s |
| Bases digest-pinned | lane's `FROM` grep snippet | 2/2 pinned |
| Image runs | `docker run --rm ci-ingester:ci probe-item-1` / bare | exit 0 / exit 2 (usage) |
| govulncheck clean | `govulncheck ./...` (v1.9.0, stdlib-only tree) | `No vulnerabilities found.` |
| gosec clean | `gosec -fmt sarif -out gosec.sarif ./...` (v2.30.0) | 0 findings |
| Fork-PR safe | fork PR #14 (Go change from `linuxdynasty:`) run 40782064652 | gosec pass on the read-only token; `gosec-sarif` skipped (no `security-events:write` on `pull_request`); closed unmerged |
| Dependabot active | `.github/dependabot.yml` (gomod + github-actions, weekly) | live `dependabot/github_actions/*` branch; earlier gomod dynamic-run failure (`No go.mod files found`, run 40778064813) is diagnosed as the pre-#10 skeleton tree, resolved by the `go.mod` landing |
| Full local floor | `make ci` (Go 1.26.10) | exit 0, coverage 100.0% at floor 100.0 |
| Sensitivity | `len(args) != 1` mutated to `!= 2` | `cmd/ingester` tests FAIL; reverted to GREEN |

SARIF upload to the Security tab happens on pushes (the
`gosec-sarif` job runs on push/schedule/dispatch only); the merge push
of this PR exercises that path.

## Incidental fix: `make lint` masked lane failures

While proving this PR, the local floor was found green on a tree with
a real `unused` finding (`var version` in `cmd/ingester/main.go`,
since removed — version stamping returns with a later `-version`
flag). Root cause: the Makefile `lint` recipe chained the lane
commands with `&&` but ended the `else` branch with an unconditional
`echo OK`, so any `vet`/`golangci`/`fmt`/`tidy` failure was swallowed
and the target exited 0 (CI runs each step separately and was never
affected; earlier `make ci` evidence on lint-clean trees stays valid).
The cap block is now the final `&&` operand, so a lane failure fails
the target. Seeded proof: with a scratch unused variable present,
`make lint` exits 2 on the `unused` finding; clean tree exits 0.
