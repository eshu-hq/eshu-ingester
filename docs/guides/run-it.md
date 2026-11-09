# Run it

Short version: don't deploy this as a service, because it isn't one.
This page covers the only two ways to run it today and what you'll observe.

## Run it locally

```sh
go build -o /tmp/ingester ./cmd/ingester
/tmp/ingester my-item-1   # accepted: exit code 0
/tmp/ingester             # bad usage: exit code 2
```

Build first rather than using `go run`. The exit code is the whole
result, and `go run` would mask it.

## Run it in a container

The [Dockerfile](../../Dockerfile) builds a static binary into a
distroless image that runs as a non-root user, with that binary as the
entrypoint. So `docker run` behaves exactly like the local run and
then exits:

```sh
docker build -t eshu-ingester:local .
docker run --rm eshu-ingester:local my-item-1   # exit 0
docker run --rm eshu-ingester:local             # exit 2 (usage)
```

Treat the image as build-and-exit-code proof for CI, not as a
deployable workload. The image behaviour was recorded in
[docs/ci/supply-chain-6.md](../ci/supply-chain-6.md) when the lane was
set up.

## What you'll observe

Exit codes are the contract: `0` accepted, `1` intake failed, `2` bad
usage. The full table is in [Reference](../reference/README.md).

That's all there is to observe. The binary prints nothing, opens no
port, writes no files, and keeps no state between runs. There's no
health endpoint to poll and nothing to keep running — if you were
hoping to watch it with a liveness probe, there is honestly nothing to
watch yet.

## What CI checks

`make ci` is the local version of the CI lanes: `go vet`,
`golangci-lint`, a format check, a `go mod tidy` diff, the 500-line
file cap, the test suite, and the race suite with a coverage floor.
See [Contributing](../contributing/README.md) and the [CI notes](../ci).
