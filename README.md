# worker-service

A production-ready Go worker that consumes events from **Redis Streams**, stores
job state in **PostgreSQL**, and executes blockchain transactions via the
**DID Root SDK** + **Hara Core Blockchain Library**.

```
Redis Stream ──XREADGROUP──► Worker Pool ──► EventService ──► BlockchainAdapter
                                  │                                    │
                                  ▼                                    ▼
                            PostgreSQL (jobs)              did-root-sdk / hara-core
```

---

## Folder Structure

```
.
├── cmd/
│   ├── worker/        # Main entrypoint
│   └── dlq-reader/    # Ops tool: tails the dead-letter queue
│
├── internal/
│   ├── config/        # Env-based configuration + validation
│   ├── domain/        # Core types: Event, Job, Payloads, Errors
│   ├── infra/
│   │   ├── db/        # PostgreSQL connection + JobRepository impl + migrations
│   │   └── redis/     # Redis client, consumer group bootstrap, DLQ push
│   ├── mocks/         # Test doubles for BlockchainService + JobRepository
│   ├── repository/    # JobRepository interface (port)
│   ├── sdk/           # BlockchainAdapter — ONLY place SDK imports appear
│   ├── service/       # EventService: idempotency → DB → blockchain orchestration
│   └── worker/        # Handler (parse/validate/retry) + Pool (concurrency loop) + HTTP server
│
├── pkg/               # Shared utilities: retry backoff, Prometheus metrics
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Quick Start

### 1. Configure

```bash
cp .env.example .env
# Edit .env — set REDIS_URL, DB_URL, RPC_URLS, PRIVATE_KEY at minimum
```

### 2. Start dependencies

```bash
make docker-up    # starts Redis + PostgreSQL via docker-compose
```

### 3. Run the worker

```bash
make run
```

### 4. Inspect the DLQ

```bash
make run-dlq
```

---

## Environment Variables

| Variable            | Required | Default            | Description                                      |
|---------------------|----------|--------------------|--------------------------------------------------|
| `REDIS_URL`         | ✅       | —                  | Redis connection URL                             |
| `STREAM_NAME`       | ✅       | —                  | Redis stream to consume                          |
| `GROUP_NAME`        | ✅       | —                  | Consumer group name                              |
| `DB_URL`            | ✅       | —                  | PostgreSQL DSN                                   |
| `RPC_URLS`          | ✅       | —                  | Comma-separated blockchain RPC endpoints         |
| `PRIVATE_KEY`       | ✅       | —                  | Hex-encoded wallet private key                   |
| `CONSUMER_NAME`     | ❌       | `worker-<hostname>`| Unique consumer identity (set per replica)       |
| `WORKER_CONCURRENCY`| ❌       | `10`               | Max concurrent goroutines                        |
| `MAX_RETRY`         | ❌       | `3`                | Max blockchain retry attempts per event          |
| `RETRY_BASE_DELAY`  | ❌       | `1s`               | Base delay for exponential backoff               |
| `BATCH_SIZE`        | ❌       | `10`               | Events per XREADGROUP call                       |
| `POLL_INTERVAL`     | ❌       | `100ms`            | Block timeout for XREADGROUP                     |
| `SHUTDOWN_TIMEOUT`  | ❌       | `30s`              | Max time to drain in-flight jobs on SIGTERM      |
| `HNS_NAME`          | ❌       | —                  | HNS contract name (mutually exclusive with ABI)  |
| `ABI_PATH`          | ❌       | —                  | Path to contract ABI JSON file                   |
| `DLQ_SUFFIX`        | ❌       | `:dlq`             | Appended to STREAM_NAME to form DLQ stream name  |
| `SERVER_PORT`       | ❌       | `8080`             | Port for /healthz and /metrics                   |

---

## Processing Flow

For each Redis stream message the worker:

1. **Parses** the raw stream entry into `domain.Event`
2. **Validates** the event (id present, type recognised, payload non-empty)
3. **Idempotency check** — queries `jobs` table by `event_id`; skips if already `success`
4. **Creates** a `pending` job row in PostgreSQL
5. **Dispatches** to the correct `BlockchainService` method with retry + exponential backoff
6. **Updates** the job row to `success` (with tx hashes) or `failed` (with error)
7. **ACKs** the message if successful; pushes to DLQ and ACKs if retries are exhausted

---

## Observability

### Health check

```
GET http://localhost:8080/healthz
→ 200 {"status":"ok"}
```

### Prometheus metrics

```
GET http://localhost:8080/metrics
```

Key metrics:

| Metric                                  | Type      | Description                              |
|-----------------------------------------|-----------|------------------------------------------|
| `worker_events_received_total`          | Counter   | Events read from Redis                   |
| `worker_events_processed_total{status}` | Counter   | Events by outcome: success/failed/skipped|
| `worker_events_retried_total`           | Counter   | Total retry attempts                     |
| `worker_events_dlq_total`               | Counter   | Events routed to DLQ                     |
| `worker_event_process_duration_seconds` | Histogram | End-to-end processing latency            |

---

## Horizontal Scaling

Each replica **must** have a unique `CONSUMER_NAME`. With the default
`worker-<hostname>` value, Kubernetes pods get unique names automatically.

All replicas join the **same consumer group** — Redis Streams guarantees
each message is delivered to exactly one consumer in the group.

```yaml
# k8s Deployment excerpt
env:
  - name: CONSUMER_NAME
    valueFrom:
      fieldRef:
        fieldPath: metadata.name   # e.g. worker-service-7d9f8b-xkqp2
```

---

## Running Tests

```bash
make test       # all tests with race detector
make cover      # coverage HTML report → coverage.html
```

---

## Architecture Decisions

**SDK isolation**: `internal/sdk/blockchain_adapter.go` is the *only* file that
imports `hara-core-blockchain-lib` and `did-root-sdk`. Everything else depends
on the `service.BlockchainService` interface, making the SDK trivially swappable
and all business logic independently testable with `mocks.MockBlockchainService`.

**Idempotency via `event_id` UNIQUE constraint**: even under concurrent processing
across replicas, the database constraint prevents double-processing at the storage
level. The application-level check is an optimisation to avoid unnecessary
blockchain calls.

**ACK-after-success only**: messages stay in the Redis PEL (Pending Entries List)
until processing is confirmed successful. After exhausting retries the event is
written to the DLQ stream and then ACKed — keeping the main stream clean while
preserving the failed payload for manual inspection or replay.

**Graceful shutdown**: `signal.NotifyContext` propagates cancellation through the
entire call tree. The pool's goroutine `WaitGroup` ensures no in-flight job is
abandoned mid-transaction.
