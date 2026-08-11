
Yep. We finally found the actual project 😭.

No fake complexity. No “Month 3: implement consensus 🤓.” No feature exists just because it teaches a buzzword.

The product has one promise:

> **Give Grug an application. Grug deploys it and keeps it running.**

You start with one process on one computer. Over 6–9 months, you keep making that promise stronger. Every new feature should arise because the previous version has a concrete limitation.

This is the roadmap I'd actually follow.

---

# Grug Platform — 6–9 month engineering roadmap

At the end, a developer should be able to take:

```text
hello-api/
├── Dockerfile
├── go.mod
└── main.go
```

and do:

```bash
grug deploy

grug status
grug logs

grug config set ENV=production

grug scale 3

grug rollback

grug destroy
```

Grug handles:

```text
             Developer
                 │
           grug deploy
                 │
                 ▼
          GRUG CONTROL PLANE
                 │
        ┌────────┼────────┐
        ▼        ▼        ▼
      node-1   node-2   node-3
        │        │        │
      my-api   my-api   my-api
        └────────┼────────┘
                 │
           Load Balancer
                 │
                User
```

Eventually you'll deploy **Grug itself onto Kubernetes**, which gives you the beautiful comparison:

> Here's my shitty implementation of scheduling/reconciliation/service discovery. Here's how Kubernetes solves the same problem.

That's the project.

---

# Stage 1 — Make Grug deploy one application

### ~Weeks 1–3

Forget distribution.

Your first problem:

> I have a Dockerized application. I want Grug to run it.

Build:

```bash
grug deploy ./my-api
grug status
grug logs
grug stop
grug start
grug destroy
```

For V1, `grug deploy ./my-api` should:

1. read a small `grug.yaml`;
2. build the Dockerfile;
3. create a container;
4. start it;
5. record the deployment;
6. print where it's running.

Example configuration:

```yaml
name: hello-api
port: 8080
health: /healthz
```

User experience:

```text
$ grug deploy ./hello-api

Building hello-api...
Starting container...

✓ hello-api deployed
✓ http://localhost:9001
```

Then:

```text
$ grug status

APP         STATUS     PORT
hello-api   running    9001
```

And:

```text
$ grug logs hello-api

2026-09-01T13:21:03 server listening :8080
2026-09-01T13:21:08 GET /healthz 200
```

Store deployment information in SQLite initially.

Don't start with Postgres because you don't have a distributed system yet.

### Make it fail

Deploy an application that:

* exits immediately;
* has a broken Dockerfile;
* listens on the wrong port;
* returns 500 from `/healthz`.

Grug should distinguish these failures.

### What you're learning

This forces you underneath Docker:

```text
process
  ↓
Linux process lifecycle
  ↓
signals
  ↓
ports/sockets
  ↓
container
  ↓
namespaces/cgroups
  ↓
Docker networking
```

Learn particularly well:

* processes and signals;
* PID 1;
* graceful shutdown;
* TCP sockets;
* ports;
* container lifecycle;
* environment variables;
* exit codes;
* stdout/stderr.

Use `ps`, `ss`, `lsof`, `curl`, `docker inspect`.

Don't merely call Docker and move on. Understand what Grug is asking the OS/container runtime to accomplish.

### Reading

Continue **Understanding Distributed Systems**, but don't force distributed chapters yet.

Continue **A Philosophy of Software Design** slowly while designing Grug's Go packages.

Your design question becomes:

> Does the rest of Grug know *how* Docker works, or is that knowledge contained behind one useful abstraction?

That's Ousterhout territory.

---

# Stage 2 — Deployments become real deployments

### ~Weeks 4–6

There's an obvious problem with V1.

You change your application and run:

```bash
grug deploy
```

What happens to the old one?

Now implement versions.

```bash
grug deploy
grug releases
grug rollback
```

Example:

```text
$ grug releases

VERSION    STATUS
a81f92     current
39cd21     previous
18af82
```

Deploy version B.

Grug should:

```text
build B
   ↓
start B
   ↓
wait for /healthz
   ↓
B healthy?
   ↓ yes
switch traffic A → B
   ↓
stop A
```

If B never becomes healthy:

```text
Deployment failed.
Keeping a81f92 active.
```

Now:

```bash
grug rollback
```

should restore the previous working version.

### Add configuration

```bash
grug config set hello-api LOG_LEVEL=debug
grug config list hello-api
grug config unset hello-api LOG_LEVEL
```

Changing config and redeploying should inject the new configuration.

Secrets must not appear in ordinary output.

### What you're learning

Now you're learning:

* deployment state machines;
* health checks;
* readiness;
* immutable releases;
* configuration;
* rollback;
* failure atomicity;
* backward compatibility;
* zero/minimal-downtime deployment.

And you start developing an extremely useful instinct:

> A deployment isn't “start new program.” It's a state transition that can fail halfway through.

That thought follows you straight into Kubernetes and production CI/CD.

---

# Stage 3 — Grug keeps applications alive

### ~Weeks 7–9

Current Grug deploys applications.

But:

```bash
docker kill hello-api
```

and it's dead forever.

That's a pretty shit platform 😂.

Change the promise:

> If I tell Grug an application should be running, Grug should keep it running.

Introduce:

```bash
grug scale hello-api 3
```

Store:

```text
desired replicas = 3
```

Then continuously observe:

```text
desired = 3
actual  = 2

      ↓

start another instance
```

Your platform now has a small **reconciliation loop**.

Don't copy Kubernetes controllers.

Try designing it yourself first.

`grug status` becomes:

```text
$ grug status

APP         VERSION   DESIRED   RUNNING   STATUS
hello-api   a81f92    3         3         healthy
```

Kill one container manually.

```bash
docker kill ...
```

Do nothing else.

Grug should eventually return to:

```text
3/3 healthy
```

Now deliberately make the application crash every ten seconds.

You will discover another problem:

```text
crash
restart
crash
restart
crash
restart
AAAAAAAAAAAA
```

Figure out what policy Grug should have.

### What you're learning

This is your first major platform concept:

**desired state vs observed state.**

You're also learning:

* control loops;
* reconciliation;
* supervision;
* health vs existence;
* crash loops;
* retry/backoff;
* convergence;
* eventual consistency.

And NOW when you later see:

```yaml
spec:
  replicas: 3
```

you understand something much deeper than YAML.

---

# Stage 4 — Cross the distributed-systems boundary

### ~Weeks 10–14

Until now:

```text
Grug
 │
 ├── app
 ├── app
 └── app

ONE MACHINE
```

Great.

Now impose the requirement:

> One machine isn't enough anymore.

Create two programs:

```bash
grug server
grug agent
```

Architecture:

```text
               grug server
             CONTROL PLANE
                  │
          ┌───────┴───────┐
          │               │
          ▼               ▼
      grug agent       grug agent
        node A           node B
          │               │
        apps            apps
```

`grug agent` runs on every machine and reports:

```text
node ID
CPU capacity
memory capacity
running applications
heartbeat
```

Now:

```bash
grug nodes
```

returns:

```text
NODE       CPU       MEMORY      STATUS
node-a     4 cores   8 GB        healthy
node-b     8 cores   16 GB       healthy
```

And:

```bash
grug deploy hello-api
```

requires the control plane to choose a node.

Start stupid:

> node with most available memory wins.

Congratulations.

You wrote a tiny scheduler.

### Now murder node A.

Grug eventually notices its heartbeat disappeared.

Suppose:

```text
desired replicas = 3

node A: DEAD
node B: 1 replica

actual replicas = 1
```

Grug should schedule replacements onto healthy capacity.

THIS is the moment you've genuinely built a distributed system.

### What you're learning

Now UDS becomes your bible.

You'll encounter:

* partial failure;
* RPC;
* serialization;
* heartbeats;
* leases;
* failure detection;
* timeouts;
* retries;
* idempotency;
* split-brain-ish situations;
* stale state;
* clocks;
* distributed ownership.

And constantly ask:

> How do I know node A is dead rather than merely unreachable?

There is no magic `isComputerDead()` syscall.

Welcome to distributed systems. 😭

---

# Stage 5 — Networking becomes a platform problem

### ~Weeks 15–17

You now have:

```text
node A → hello-api-1
node B → hello-api-2
node C → hello-api-3
```

Cool.

How does the user reach them?

Add a Grug proxy/load balancer.

```text
                    :80
                     │
                 Grug Proxy
               /           \
              ▼             ▼
          node A:9123   node B:8431
             API           API
```

Applications get names:

```text
hello-api.grug.local
```

Traffic only goes to healthy instances.

Deploy new instances and routing updates.

Kill one and it disappears from routing.

Now you've earned:

* reverse proxies;
* L4 vs L7;
* DNS;
* service discovery;
* load balancing;
* connection draining;
* health-based routing;
* TLS termination.

Use `dig`, `curl`, `ss`, `tcpdump` aggressively.

At this point you should be able to draw a packet from browser → application.

---

# Stage 6 — AWS SAA

### ~Weeks 18–21

Now move the distributed Grug deployment onto AWS.

Don't redesign it.

Map what you already understand onto cloud primitives.

```text
Your concept        AWS

machine             EC2
network             VPC
firewall            Security Group
load balancer       ALB
DNS                 Route53
database            RDS
object storage      S3
identity            IAM
metrics             CloudWatch
```

Deploy:

```text
                    Route53
                       │
                      ALB
                       │
            ┌──────────┴──────────┐
            │                     │
           EC2                   EC2
        grug agent            grug agent
            │                     │
            └────── apps ─────────┘

                 RDS / S3
```

Build the VPC manually.

Understand:

* CIDRs;
* subnets;
* route tables;
* IGW/NAT;
* security groups;
* IAM roles;
* availability zones.

Then reproduce your infrastructure using Terraform.

**Only now Terraform enters.**

You understand what it's creating.

### Cert

This is your **AWS SAA window**.

Study whatever exam-only gaps Grug didn't cover and take the exam.

Cloud Practitioner can be knocked out much earlier as the vocabulary cert.

---

# Stage 7 — Make Grug observable

### ~Weeks 22–24

Your system is now sufficiently complicated that:

> “idk bro it isn't working”

is unacceptable.

Instrument the control plane and agents.

Metrics:

```text
deployment_success_total
deployment_failure_total

reconciliation_duration
reconciliation_errors

node_count
node_unreachable

desired_instances
running_instances

http_request_duration
http_errors
```

Structured logs should carry things like:

```text
deployment_id
application
node_id
instance_id
request_id
```

Add OpenTelemetry traces where they genuinely help.

Define a few actual SLIs/SLOs.

Then build:

```bash
grug doctor
```

Example:

```text
$ grug doctor hello-api

Control plane    OK
node-a           OK
node-b           UNREACHABLE
hello-api        DEGRADED  2/3
routing          OK
database         OK

Likely boundary:
control-plane → node-b
```

Don't make `doctor` magical AI.

Make it encode **your operational model of the system**.

That's way more educational.

### Reading

Now use the **Google SRE Workbook** alongside Grug.

Learn:

* SLIs;
* SLOs;
* error budgets;
* alerting;
* monitoring;
* incident response;
* capacity.

---

# Stage 8 — Deploy Grug onto Kubernetes

### ~Weeks 25–28

Now comes the payoff.

You built:

```text
Grug scheduler
Grug agents
Grug reconciliation
Grug health checks
Grug service discovery
Grug deployments
```

Then Kubernetes walks into the room like:

> cute.

😭

Containerize Grug itself and deploy it onto a cluster.

Hand-write:

```text
Deployment
Service
ConfigMap
Secret
ServiceAccount
RBAC
NetworkPolicy
PDB
HPA
Ingress/Gateway
PVC where appropriate
```

Now explicitly compare concepts:

```text
YOUR GRUG              KUBERNETES

application            workload
instance               Pod
node                    Node
desired instances       replicas
reconciler              controller
scheduler               kube-scheduler
agent                   kubelet-ish
health check            probes
routing                 Service
external routing        Ingress/Gateway
config                   ConfigMap
```

Don't conclude that your implementation is equivalent internally. It's deliberately tiny.

The value is:

> You now understand the **problems** those abstractions exist to solve.

### Build your own cluster

Use `kubeadm`.

Break:

```text
kubelet
CoreDNS
CNI
Service selectors
RBAC
scheduling
storage
certificates
control-plane components
```

### Cert

Take **CKA** here.

Your cert preparation becomes mostly:

> make your existing understanding fast and exam-compatible.

Rather than:

> memorize 600 kubectl commands.

---

# Stage 9 — Security / CKS

### ~Weeks 29–34

Final boss.

Assume:

> An attacker gets shell access inside a Grug application.

What can they reach?

Try it.

Then progressively restrict the system:

```text
non-root
↓
drop capabilities
↓
read-only filesystem
↓
seccomp
↓
ServiceAccounts
↓
minimal RBAC
↓
NetworkPolicies
↓
secret isolation
↓
image scanning
↓
supply-chain controls
↓
runtime monitoring
```

Every security mechanism must answer:

> **What attack does this prevent or limit?**

Not:

> CKS says add this YAML.

Threat-model the whole platform.

```text
Internet
   │
   ▼
proxy ───── trust boundary
   │
control plane
   │
   ├── database
   │
   └── agents
         │
         ▼
    user workloads
```

Compromise each component mentally and determine its blast radius.

### Cert

Take **CKS**.

Then:

## STOP COLLECTING CERTS.

Unless your job gives you a reason for another one.

---

# The whole 6–9 months

Your progression is incredibly clean:

```text
MONTH 1
Deploy one container
        │
        ▼
processes / Linux / networking

MONTH 2
Versions + health + rollback
        │
        ▼
deployment engineering

MONTH 2–3
Keep N instances alive
        │
        ▼
reconciliation / desired state

MONTH 3–4
Multiple machines
        │
        ▼
distributed systems / scheduling

MONTH 4
Route traffic
        │
        ▼
networking / service discovery

MONTH 4–5
Move to AWS
        │
        ├──────────────► AWS SAA
        ▼
cloud architecture

MONTH 5–6
Observe + break it
        │
        ▼
SRE / debugging

MONTH 6–7
Run on Kubernetes
        │
        ├──────────────► CKA
        ▼
orchestration

MONTH 7–9
Attack + harden
        │
        ├──────────────► CKS
        ▼
security
```

There are also some things I explicitly **do not want you adding** unless Grug creates the requirement:

Kafka. Redis. Microservices. Service mesh. ArgoCD. Operators. Vault. Prometheus operator stacks on day three. Five databases because they're cool.

The entire philosophy is:

```text
problem
   ↓
experience problem
   ↓
understand problem
   ↓
implement simplest solution
   ↓
discover limitations
   ↓
study production solution
```

Not:

```text
learn technology
   ↓
invent justification for technology
```

That's the **lazy senior** part.

---

And your reading stays hilariously small.

**Understanding Distributed Systems** alongside stages 3–5.

**A Philosophy of Software Design** throughout the Go implementation.

**DDIA 2e selectively** once persistence, replication, consistency, queues, or data architecture produce questions.

**Google SRE Workbook** once Grug becomes observable.

**Official Kubernetes docs** during CKA.

**Official AWS docs + SAA practice material** during AWS.

No 37-course Udemy pilgrimage.

Your project is the curriculum.

And here's the benchmark I'd care about at the end.

I don't care whether you can say:

> “I implemented a scheduler.”

I want to be able to give you:

> “Users intermittently get 502s after deploying v47, but all three instances appear healthy.”

And watch your brain automatically go:

```text
What changed?
      ↓
What path does a request take?
      ↓
DNS → LB → routing → instance → application
      ↓
Which boundaries could produce 502?
      ↓
What observation distinguishes them?
      ↓
measure
      ↓
hypothesis
      ↓
experiment
```

No random restarting.

No shotgun debugging.

No asking AI to fix it.

No `kubectl describe` divination ritual.

Just an engineer building a model of reality and narrowing uncertainty.

**That's the actual roadmap.**

Grug is merely the machine we're going to use to force you to practice it. 🗿🔥
