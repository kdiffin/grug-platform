| Milestone            | Outcome                           | Estimated effort |
| -------------------- | --------------------------------- | ---------------: |
| 0 — Walking skeleton | Runnable end-to-end foundation    |        1 evening |
| 1 — Registration     | First usable capability           |     1–2 evenings |
| 2 — Job queue        | Concurrency and state transitions |     2–3 evenings |




Alright, fresh version. The key is to make Phase 1 a tiny vertical slice of a platform—not a sad YAML generator wearing a fake moustache.

# Phase 1 outcome

A developer opens an HTMX page, enters:

* application name
* container image
* container port

Your Go platform then:

1. accepts the deployment request;
2. puts it onto an internal job queue;
3. returns immediately with `queued`;
4. processes deployments using background workers;
5. creates a Kubernetes namespace and Deployment;
6. records deployment status;
7. updates the browser until the application becomes ready.

The demo target can be something predictable:

```text
name: hello
image: nginx:alpine
port: 80
```

The visible lifecycle:

```text
queued → deploying → ready
                  ↘ failed
```

That is already a legitimate baby internal developer platform. Tiny Coolify has spawned. It knows nothing and fears no YAML.

# Keep Phase 1 brutally scoped

Build:

* one Go binary;
* one HTML page;
* HTMX status updates;
* one Kubernetes cluster, preferably `k3d`;
* one deployment workflow;
* bounded concurrent workers;
* persistent deployment records;
* generated Kubernetes manifests;
* `kubectl` execution through a small adapter.

Do not build yet:

* authentication;
* GitLab API integration;
* branch environments;
* Helm;
* Kubernetes operators;
* secrets management;
* logs UI;
* databases offered as a service;
* multiple clusters;
* rollback automation;
* WebSockets;
* microservices;
* “enterprise-grade plugin architecture,” the traditional final boss of unfinished side projects.

# Repository layout

Start with structure that reflects business capabilities rather than technical layers:

```text
grug-platform/
├── cmd/
│   └── platform/
│       └── main.go
│
├── internal/
│   ├── applications/
│   │   ├── application.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── deployments/
│   │   ├── deployment.go
│   │   ├── service.go
│   │   ├── queue.go
│   │   └── worker.go
│   │
│   ├── cluster/
│   │   ├── cluster.go
│   │   ├── manifests.go
│   │   └── kubectl.go
│   │
│   ├── web/
│   │   ├── server.go
│   │   ├── handlers.go
│   │   └── templates/
│   │       ├── index.html
│   │       └── deployment.html
│   │
│   └── storage/
│       └── jsonstore.go
│
├── deployments/
│   └── demo-app/
│       └── Dockerfile
│
├── compose.yaml
├── Makefile
├── go.mod
└── README.md
```

The dependency direction should be:

```text
web → applications/deployments → cluster
                      ↓
                   storage
```

Important rule: `deployments` should not know about HTTP or HTML. The web module calls it just like a future CLI, GitLab webhook, or API would.

# Core domain models

Keep them boring:

```go
type Application struct {
	ID        string
	Name      string
	Image     string
	Port      int
	CreatedAt time.Time
}

type Deployment struct {
	ID            string
	ApplicationID string
	Namespace     string
	Status        Status
	Error         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Status string

const (
	StatusQueued    Status = "queued"
	StatusDeploying Status = "deploying"
	StatusReady     Status = "ready"
	StatusFailed    Status = "failed"
)
```

Then define the boundaries before implementations:

```go
type DeploymentRepository interface {
	Save(context.Context, Deployment) error
	FindByID(context.Context, string) (Deployment, error)
}

type Cluster interface {
	Deploy(context.Context, DeploymentSpec) error
	Ready(context.Context, DeploymentSpec) (bool, error)
}
```

Interfaces belong near the code that consumes them. Don’t make an `interfaces/` junk drawer—it becomes the architectural equivalent of that chair in your room covered in clothes.

# Concurrency design

This is where Phase 1 becomes educational instead of merely CRUD.

## Request path

The HTTP handler should not run `kubectl` directly:

```go
func (s *Service) Request(ctx context.Context, cmd DeployCommand) (Deployment, error) {
	// Validate input.
	// Create queued deployment.
	// Persist it.
	// Submit its ID to the queue.
	// Return immediately.
}
```

Conceptually:

```go
jobs := make(chan DeploymentJob, 16)
```

Start a fixed number of workers:

```go
for i := 0; i < workerCount; i++ {
	go worker.Run(ctx, jobs)
}
```

Each worker:

```text
receive job
    ↓
mark deploying
    ↓
generate manifest
    ↓
apply through cluster adapter
    ↓
wait for readiness
    ↓
mark ready or failed
```

What this teaches:

* goroutine lifecycle;
* bounded concurrency;
* buffered channels;
* backpressure;
* cancellation with `context.Context`;
* graceful shutdown;
* ownership of mutable state;
* avoiding goroutine leaks;
* error propagation.

## Bounded concurrency matters

Do not start one goroutine per deployment forever:

```go
go deploy(job)
```

That allows an unbounded number of expensive `kubectl` processes.

A three-worker pool means:

```text
100 submitted jobs
3 actively deploying
16 buffered
remaining submissions receive backpressure or "queue full"
```

That is real platform-engineering thinking: control resource consumption rather than merely making things asynchronous.

## Graceful shutdown

Your root context should be cancelled on `SIGINT` or `SIGTERM`.

Shutdown order:

1. stop accepting HTTP requests;
2. cancel the application context;
3. workers stop starting new work;
4. active cluster operations receive cancellation;
5. wait for workers using `sync.WaitGroup`;
6. exit.

Do not close the jobs channel from random components. The component that creates and owns the channel decides when to close it.

# Persistence without dragging in a database

For Phase 1, use a JSON-backed repository:

```text
.data/deployments.json
```

But funnel writes through one owner goroutine:

```text
handlers/workers → storage commands channel → storage goroutine → JSON file
```

This gives you an actor-style concurrency exercise:

* only one goroutine owns mutable repository state;
* callers send commands;
* commands carry response channels;
* reads and writes cannot race.

Conceptually:

```go
type saveCommand struct {
	deployment Deployment
	result     chan error
}

type findCommand struct {
	id     string
	result chan findResult
}
```

Later, replace `jsonstore` with PostgreSQL without changing deployment logic. That replacement proves whether your boundary is real.

If the actor store starts swallowing half the project, simplify it to `sync.RWMutex`. The goal is understanding ownership, not reenacting Erlang inside 600 lines of Go.

# Milestone breakdown

## Milestone 0 — Walking skeleton

Goal: prove the binary and browser path work.

Build:

* `go mod init`;
* HTTP server;
* `GET /`;
* embedded HTML templates;
* graceful server shutdown;
* basic Makefile commands.

Done when:

```bash
go run ./cmd/platform
curl localhost:8080
```

returns the page.

Estimated effort: one evening.

## Milestone 1 — Application registration

Goal: accept and display an application.

Build:

* application model;
* validation;
* HTML form;
* `POST /applications`;
* in-memory repository;
* application list rendered through HTMX.

Validation rules:

* Kubernetes-safe name;
* non-empty image;
* port between `1` and `65535`;
* duplicate names rejected.

Done when you can register `nginx:alpine` and see it in the UI.

Do not touch Kubernetes yet.

## Milestone 2 — Deployment state machine

Goal: model deployments before doing real deployments.

Build:

* deployment statuses;
* deployment service;
* job channel;
* two background workers;
* fake cluster adapter that sleeps and sometimes fails;
* status polling with HTMX.

Example fake adapter:

```go
type FakeCluster struct {
	Delay time.Duration
}
```

This milestone is extremely important. It lets you test concurrency without Kubernetes making every bug look like ancient forbidden networking magic.

Done when multiple requests visibly move through:

```text
queued → deploying → ready
```

Run:

```bash
go test -race ./...
```

## Milestone 3 — Persistent repository

Goal: survive process restarts.

Build:

* JSON repository;
* single-owner storage goroutine or carefully locked state;
* atomic file replacement;
* startup loading;
* tests with temporary directories.

For atomic persistence:

1. serialize current state;
2. write a temporary file;
3. `fsync`/close it;
4. rename it over the original.

Do not overwrite the main file directly. A crash halfway through a write should not turn your platform’s memory into modern art.

Done when deployments remain visible after restarting the binary.

## Milestone 4 — Real Kubernetes adapter

Goal: replace the fake cluster with `kubectl`.

Build:

* manifest renderer;
* namespace creation;
* Deployment creation;
* Service creation;
* command execution using `exec.CommandContext`;
* readiness checking;
* captured stderr;
* timeout handling.

Keep the boundary narrow:

```go
type Cluster interface {
	Deploy(context.Context, DeploymentSpec) error
}
```

Generated resources:

```text
Namespace: platform-hello
Deployment: hello
Service: hello
```

Done when submitting the form produces real resources:

```bash
kubectl get all -n platform-hello
```

## Milestone 5 — Operational polish

Goal: make failure understandable.

Add:

* `/healthz`;
* `/readyz`;
* structured logs using `log/slog`;
* deployment timeout;
* queue-full response;
* worker count configuration;
* clean shutdown;
* deployment error displayed in UI.

Useful configuration:

```text
HTTP_ADDR=:8080
WORKER_COUNT=2
QUEUE_CAPACITY=16
DEPLOY_TIMEOUT=2m
DATA_PATH=.data/deployments.json
```

Done when you can intentionally submit a broken image and clearly observe:

```text
queued → deploying → failed
reason: deployment readiness timeout
```

# HTMX interaction

You only need three endpoints:

```text
GET  /                          full page
POST /applications             create and enqueue
GET  /deployments/{id}         deployment status fragment
```

After form submission, return something like:

```html
<div
  id="deployment-123"
  hx-get="/deployments/123"
  hx-trigger="every 1s"
  hx-swap="outerHTML"
>
  Status: queued
</div>
```

When the deployment reaches `ready` or `failed`, return the fragment without `hx-trigger`, so polling stops.

No WebSocket kingdom. No frontend state-management bloodline. Just server-rendered HTML doing its honest little job.

# Tests worth writing by hand

Focus on behavior rather than mocking every molecule.

## Service test

Submit a deployment and verify:

```text
initial status = queued
job enters queue
```

## Worker success test

Use the fake cluster and verify:

```text
queued → deploying → ready
```

## Worker failure test

Configure the fake cluster to fail and verify:

```text
queued → deploying → failed
error is stored
```

## Cancellation test

Cancel the context during deployment and verify the worker exits.

## Concurrency test

Submit 50 jobs and assert:

* all reach a terminal state;
* no records disappear;
* `go test -race` reports nothing;
* active cluster calls never exceed worker count.

That final assertion teaches you what “bounded concurrency” actually guarantees.

# Suggested build order

Build this in vertical slices:

1. form creates an in-memory application;
2. form creates a queued deployment;
3. one worker completes fake deployments;
4. multiple workers process concurrently;
5. HTMX displays status transitions;
6. persistence survives restarts;
7. fake cluster becomes real Kubernetes;
8. failure paths and shutdown become clean.

Commit after each working slice. You should always have a runnable platform, even if today’s platform merely sleeps for two seconds and boldly announces success.

# Definition of Phase 1 done

Phase 1 is complete when you can demonstrate:

```text
1. Start local k3d cluster.
2. Start the Go binary.
3. Open the browser.
4. Submit nginx:alpine.
5. Immediately see "queued."
6. Watch it become "deploying."
7. Watch it become "ready."
8. Verify the namespace and workload with kubectl.
9. Restart the platform.
10. See the deployment history still present.
11. Submit several apps and observe bounded parallel processing.
12. Submit a broken image and see a useful failure.
```

After that, Phase 2 can introduce Git repository integration, generated GitLab CI, branch namespaces, container builds, and proper PostgreSQL storage. But Phase 1 should remain about one lesson:

> A platform is an asynchronous control plane that accepts desired state, performs bounded work, and reports observed state.

That mental model will carry much further than merely learning how to call `kubectl` from Go.
