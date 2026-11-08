# Work on the batcher

The reusable part is [internal/intake](../../internal/intake/intake.go).
`cmd/ingester` is a thin wrapper on top, so most seed work lands in the
library and its tests.

## The pieces

- An `Item` is an ID plus an opaque body. `Validate` rejects an empty
  ID (`ErrEmptyID`) or an empty body (`ErrEmptyBody`).
- A `Batcher` holds validated items until the batch is due to flush.
  Bad limits are rejected up front with `ErrInvalidConfig`: `maxItems`
  and `maxAge` must both be positive and the clock must be set.
- Adding to a full batch fails with `ErrBatchFull`.
- `FlushDue` reports whether the batch is full or its oldest item has
  reached `maxAge`. `Drain` hands items back in arrival order and
  resets the batch. `Len` and `Oldest` report current state.

Two things to respect while you work:

- The clock is injected, so keep tests deterministic — pass a fake clock.
- A `Batcher` is documented as not safe for concurrent use. Don't share
  one across goroutines.

## How the command uses it

`cmd/ingester` builds a batcher with a limit of 1 item and a one-hour
age limit, adds the single argument, and maps the outcome to the exit
code. The limits are constants in `defaultBatch`, not settings you can
pass in. Because the limit is 1, the age and flush logic are exercised
by the library tests rather than by the binary itself. The batch
constructor is injectable, which lets tests drive every failure branch
without any I/O.

## Verify your change

```sh
go test ./... -count=1   # quick focused run
make ci                  # the full gate before you push
```

See [Contributing](../contributing/README.md) for the lane details.
