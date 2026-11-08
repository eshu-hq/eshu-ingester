# Concepts

Eshu turns source data into facts. Before facts can be processed,
something has to accept incoming items and hand them on in sane
groups. That accepting-and-grouping step is **intake**. This repo
exists so intake can live as its own component instead of inside Eshu
core. The upstream issue references are in
[CONTRIBUTING.md](../../CONTRIBUTING.md).

What's here so far is the smallest honest version of that step: item
validation and an in-memory batcher, with tests and CI gates around
them. The long-running part — reading real input over time — is later
work. The package comment in
[cmd/ingester/main.go](../../cmd/ingester/main.go) says so in those words.

## One run, start to finish

A run checks its arguments, builds a one-item batcher, adds the
argument as an item, and exits. Every outcome maps to an exit code and
nothing else happens: no output, no files, no network calls. Diagram:
[One run, start to finish](../diagrams/README.md#one-run-start-to-finish).

## The contract, exactly as the code defines it

From `run` in [cmd/ingester/main.go](../../cmd/ingester/main.go):

| Input | Result |
| --- | --- |
| Exactly one non-empty argument (the item ID) | The item goes into a fresh one-item batch; exit `0` |
| Zero arguments, an empty argument, or more than one argument | Exit `2` (usage error) |
| The batch can't be created, or the item can't be added | Exit `1` (intake failure) |

There are no flags, no environment variables, and no stdin or stdout
protocol. The binary prints nothing. The item body is the fixed
placeholder bytes `{}`, so the argument is an ID only and no payload
is read from anywhere. Nothing is written to disk or sent anywhere.
The batch exists only in memory for the life of the process, which
ends straight away. (The `INGESTER_TEST_*` variables in the tree exist
only inside the test suite.)

## How the batcher decides to flush

An item only joins the batch after it validates, and only if there's
room. The batch becomes due when it fills up or when its oldest item
gets too old, whichever happens first, and `Drain` then hands the
items back in the order they arrived. Diagram:
[How the batcher decides to flush](../diagrams/README.md#how-the-batcher-decides-to-flush).

Wiring the batcher to a real input source is the future work this
skeleton is shaped for. It is not wiring that exists today.
