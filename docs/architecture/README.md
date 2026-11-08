# Architecture

Let's be blunt about what "architecture" means here today: two small
packages and a direction.

## What exists today

- **`internal/intake`** — item validation and an in-memory batcher
  (group by count or age, drain in arrival order). This is the seed
  worth reusing.
- **`cmd/ingester`** — a thin one-shot wrapper. One argument in, one
  exit code out, then the process ends.

That's the whole system. No server, no health endpoint, no input
stream, no persistence, no scheduler, no caller. Nothing in this tree
runs the binary on a schedule or from another service, and no
production Eshu system sends work here.

## Where this is heading

The intended place in Eshu is intake, before fact processing. A future
caller would run the binary (or its successor) and act on the exit
code. That's direction, per the repo's stated purpose in
[CONTRIBUTING.md](../../CONTRIBUTING.md) — not wiring that exists
today, in this tree or anywhere else we can point you to. The overview
diagram (dashed arrow and all) is in
[Diagrams](../diagrams/README.md#where-the-ingester-sits).

## Dependencies

None, in either direction. [go.mod](../../go.mod) declares no
requirements and the code uses only the Go standard library. Eshu core
(`eshu-hq/eshu`) is where fact processing lives today and where this
component's tracking issues live. Nothing calls this repo yet, and it
calls nothing back.
