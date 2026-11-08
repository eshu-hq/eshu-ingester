# ADR-0001: Skeleton first, no cutover authority

## Status

Accepted. Recorded in [CONTRIBUTING.md](../../CONTRIBUTING.md): "The
ingester holds no cutover authority until its reviewed cutover (Refs
`eshu-hq/eshu-ingester#1`, part of `eshu-hq/eshu#6707`)."

## Context

Intake needs its own home outside Eshu core, but the real component —
reading input over time, handing batches downstream — is a bigger
build than one step. Landing CI gates, a validated batcher, and a
one-shot command first proves the repo's plumbing without pretending
to be a service.

## Decision

Ship the smallest honest seed: `internal/intake` validation + batching
and a `cmd/ingester` one-shot CLI whose whole contract is its exit
code. No server, queue, or storage. Call it a skeleton everywhere, and
give it no cutover authority until a reviewed cutover says otherwise.

## Consequences

- Docs must never describe deployable behaviour. If a page sounds like
  a service manual, it's wrong.
- The batcher rules get exercised by library tests, since the binary's
  limit is 1 item (see [Concepts](../concepts/README.md)).
- A future ADR must supersede this one before any production wiring.
