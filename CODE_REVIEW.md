# Code Review

This review compares the current implementation with the design principles described in `/README.md`: Domain-Driven Design (DDD), CQRS, Event Sourcing, clear bounded contexts, and a backend that acts as an auditor.

## Overall Assessment

The codebase already has a solid architectural foundation and clearly aims in the right direction. The current implementation **partially fits** the README design goals: the domain model, command handlers, and context boundaries are present, but some of the critical infrastructure needed to fully realize Event Sourcing and CQRS is still missing.

## Pros

### 1. Clear bounded contexts
- `internal/wallet`, `internal/journey`, and `internal/kernel` are separated cleanly.
- This matches the README goal of modular DDD boundaries.

### 2. Event-driven aggregates
- `wallet.Card` and `journey.JourneyLog` both mutate state by applying domain events.
- Rehydration (`Rehydrate`) and `BaseAggregate` support the Event Sourcing direction described in the README.

### 3. Command/query separation is visible
- Commands and command handlers are split from query objects and query handlers.
- `CardProjection` and `CardProjectionRepository` show a CQRS-style read model intent.

### 4. Good dependency direction between contexts
- `journey` depends on a `PaymentPort` instead of directly depending on wallet internals.
- This is a good application of dependency inversion and keeps bounded contexts decoupled.

## Cons and Recommended Solutions

### 1. Event Sourcing is only partially implemented
**Observation**
- `kernel.EventStore` only exposes `Exists`, `Load`, and `Save`.
- The current events only expose `EventName()` and do not carry sequence/version metadata, persistence metadata, or ordering guarantees.

**Why this matters**
- The README explicitly requires strict ordering, auditability, and historical reconstruction.
- Without persisted event metadata and optimistic concurrency, those guarantees are not fully enforced.

**Recommended solution**
- Introduce a persisted event envelope with fields such as aggregate ID, aggregate version, event name, occurred-at timestamp, and payload.
- Make `Save` enforce optimistic concurrency with an expected aggregate version.

### 2. The CQRS read side is incomplete
**Observation**
- `CardProjection` and `CardProjectionRepository` exist, but there is no projector or event consumer that updates the projection from domain events.

**Why this matters**
- CQRS depends on a reliable read model pipeline.
- Without projection updating, query data can drift from the write model or remain unimplemented.

**Recommended solution**
- Add a projection updater that listens to wallet events and updates `CardProjectionRepository`.
- Keep projection updates explicit so the read path remains fast and auditable.

### 3. The backend is still modeled as a synchronous command processor
**Observation**
- The API layer directly invokes command handlers synchronously.

**Why this matters**
- The README describes the backend primarily as an asynchronous auditor, not the latency-critical decision maker.
- The current shape is useful for development, but it does not yet reflect the target edge/auditor flow.

**Recommended solution**
- Treat the current handlers as a development adapter, then introduce an ingestion path for gate events that can be processed asynchronously.
- Keep the backend as the system of record while moving edge decisions outside the central request path.

### 4. Application contracts are inconsistent around payment authorization
**Observation**
- `journey.PaymentPort` defines `AuthorizePayment(mediaId string, referenceId string, amount int) error`.
- `wallet.PaymentAdapter` currently exposes `AuthorizePayment(ctx context.Context, cardNo string, referenceId string, amount int) error`.

**Why this matters**
- This makes real wiring between both contexts harder and shows that the application boundary is not fully stabilized yet.
- It also drops clarity around context propagation.

**Recommended solution**
- Standardize the port contract on both sides.
- Prefer including `context.Context` in the port interface so application-layer operations can propagate deadlines and cancellation consistently.

### 5. Important architectural behavior is not protected by tests
**Observation**
- The repository currently has no Go test files.

**Why this matters**
- The most important risks in DDD/CQRS/Event Sourcing systems are behavioral: replay correctness, invariant enforcement, idempotency, and projection consistency.

**Recommended solution**
- Add unit tests for aggregate behavior, rehydration, duplicate payment handling, insufficient balance, and journey start/end rules.
- Add application tests for the command handlers and port interactions.

## Final Conclusion

The codebase **does fit the README design principles at a structural level**, especially for DDD boundaries, event-based aggregates, and CQRS intent. However, it **does not yet fully satisfy the operational promises** made in the README around auditability, strict ordering, asynchronous reconciliation, and mature read-model handling.

In short:
- **Good foundation:** DDD structure, CQRS intent, domain events, ports/adapters.
- **Main gaps:** missing event-store guarantees, incomplete projection flow, synchronous backend flow, inconsistent application contracts, and missing tests.
