# grug-platform plan

Build a small platform to relearn software development by hand.

AI may explain concepts, review code, or help diagnose errors. It should not write
the implementation or complete TODOs.

## Shape

One Go binary with two clients:

```text
grug app ...
grug deploy ...
grug deployment ...
grug tui
```

Both clients call the same application layer. The TUI does not execute the CLI,
and neither client owns validation, persistence, queueing, or Kubernetes logic.

```text
cmd/grug/          wiring and process lifecycle
internal/app/      workflows and validation
internal/cli/      arguments and terminal output
internal/tui/      Bubble Tea UI
internal/cluster/  fake and Kubernetes adapters
internal/storage/  persistence
```

Keep `internal/app` as one deep module until splitting it solves a real problem.
Use the standard library for the CLI and Bubble Tea for the TUI.

## Milestones

### 0 — CLI skeleton

Estimated time: 1 evening.

Create the `grug` command with help, subcommand routing, useful errors, exit codes,
configuration, and graceful shutdown.

Done when valid, invalid, help, and cancelled commands behave predictably.

### 1 — Applications

Estimated time: 1–2 evenings.

Add application registration and listing with in-memory storage. Validate names,
images, ports, and duplicates in the application layer.

Done when `app add` and `app list` work without containing business logic.

### 2 — Deployments

Estimated time: 2–3 evenings.

Add a bounded queue, fixed workers, cancellation, and these states:

```text
queued -> deploying -> ready
                     -> failed
```

Use a fake cluster first. Add commands to request and inspect deployments.

Done when concurrent jobs show correct transitions and backpressure.

### 3 — Persistence

Estimated time: 1–2 evenings.

Store records as JSON using atomic file replacement. Keep storage behind the
application boundary.

Done when data survives restarts and `go test -race ./...` passes.

### 4 — Kubernetes

Estimated time: 2–3 evenings.

Add a small `kubectl` adapter for namespaces, Deployments, Services, readiness,
timeouts, cancellation, and useful errors.

Done when a deployment succeeds in a disposable `k3d` cluster.

### 5 — TUI

Estimated time: 2–4 evenings.

Build a thin Bubble Tea client over the same application API. Support Vim-style
movement, contextual help, loading and error states, refresh, and clean quit.

Done when registration, deployment, listing, and inspection produce the same
state through both clients.

### 6 — Polish

Estimated time: 1–2 evenings.

Add structured logs, configurable worker and queue limits, timeouts, and clear
failure reporting. Do not add a web server or frontend build system.

## Tests

- Unit-test validation, state transitions, backpressure, and cancellation.
- Test CLI parsing, output, and exit codes without starting the TUI.
- Test TUI updates with messages and fake application dependencies.
- Test storage with temporary directories.
- Run race tests before each milestone is considered done.

Prefer small files, few dependencies, and code written only when the current
milestone needs it.
