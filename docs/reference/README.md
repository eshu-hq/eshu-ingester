# Reference

## Exit codes (`cmd/ingester`)

| Code | Meaning |
| --- | --- |
| `0` | Accepted — the item went into a fresh one-item batch |
| `1` | Intake failure — the batch couldn't be created, or the item couldn't be added |
| `2` | Bad usage — anything other than exactly one non-empty argument |

## Library errors (`internal/intake`)

| Error | When |
| --- | --- |
| `ErrEmptyID` | An item has no identifier |
| `ErrEmptyBody` | An item has no payload |
| `ErrInvalidConfig` | `maxItems` or `maxAge` is not positive, or the clock is unset |
| `ErrBatchFull` | `Add` against a batch already holding `maxItems` |

Key API: `Item` (+ `Validate`), `NewBatcher(maxItems, maxAge, now)`,
`Add`, `FlushDue`, `Drain`, `Len`, `Oldest`. Source:
[internal/intake/intake.go](../../internal/intake/intake.go). A
`Batcher` is not safe for concurrent use.

## Configuration surface

| Knob | Value today |
| --- | --- |
| Flags | none |
| Environment variables | none (`INGESTER_TEST_*` exist only in tests) |
| Command batch limit | 1 item, constant in `defaultBatch` |
| Command age limit | 1 hour, constant in `defaultBatch` |
| Item body | fixed placeholder bytes `{}` |
| Ports, files, network | none |

## Module facts

| | |
| --- | --- |
| Module | `github.com/eshu-hq/eshu-ingester` |
| Go | 1.26.10 (see [go.mod](../../go.mod)) |
| Third-party dependencies | none — standard library only |
