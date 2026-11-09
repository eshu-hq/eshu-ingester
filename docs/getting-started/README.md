# Getting started with eshu-ingester

You'll need Go 1.26.10 (the `go` line in [go.mod](../../go.mod)). The
module has no third-party dependencies — standard library only.

## Build and run once

Build the binary first, because `go run` would mask the exit codes and
those are the whole point:

```sh
go build -o /tmp/ingester ./cmd/ingester
/tmp/ingester my-item-1   # accepted: exit code 0
echo $?
/tmp/ingester             # bad usage: exit code 2
```

That's the entire interface today: run the binary with one item ID,
read the exit code. The binary prints nothing.

## What just happened

The run built a one-item batcher in memory, added your argument as an
item ID (with placeholder body bytes `{}`), and exited. No files, no
network, no stored state. [Concepts](../concepts/README.md) explains
the contract exactly as the code defines it.

## Where to go next

- Container run and CI: [Run it](../guides/run-it.md).
- Working on the batcher itself: [Work on the batcher](../guides/work-on-the-batcher.md).
- Honest scope — what this is not: [Architecture](../architecture/README.md).
