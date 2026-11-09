# eshu-ingester docs

Straight up: this repo is an early skeleton, not a running service. It
holds a small batching library (`internal/intake`) and a one-shot
command (`cmd/ingester`) that accepts a single item ID and exits.
Nothing here listens on a port, reads a queue, or stores anything, and
no production Eshu system sends work to it. These docs describe what
the seed code proves — please don't read them as a service manual.

## Start here

| Section | What you'll find |
| --- | --- |
| [Getting started](getting-started/README.md) | Build the binary and run it once |
| [Guides](guides/README.md) | Run it locally or in a container; work on the batcher |
| [Concepts](concepts/README.md) | The one-shot contract, the batching rules, and where the code stops |
| [Reference](reference/README.md) | Exit codes, errors, and the (short) config surface |
| [Architecture](architecture/README.md) | What exists today vs. the planned direction |
| [Diagrams](diagrams/README.md) | Every diagram in one place, each with a plain-text version |
| [Decisions](decisions/README.md) | Architecture decision records and how to add one |
| [Contributing](contributing/README.md) | `make ci`, CI lanes, and the coverage-floor rule |

## Existing docs (kept where they are)

| Doc | What's in it |
| --- | --- |
| [CI notes](ci) | Dated notes on how each CI gate in this repo was set up and proved |
| [Contributing](../CONTRIBUTING.md) | Local gates, CI lanes, and the upstream tracking issues (canonical) |

Back to the [repo README](../README.md).
