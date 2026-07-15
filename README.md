# go-ic-wallet
A Simple reconstruction of japans IC Card System. It is just for me, for learning GOlang with EventSourceing, CQRS and DDD

# Description of the Project
During a recent trip to Japan, I was fascinated by the incredible speed and efficiency of their IC Card system (such as Suica and Pasmo) used for public transportation.

Curious about the underlying mechanics of how such a high-throughput, low-latency system operates under the hood, I decided to build a simulation of this system as a personal, hands-on project.

# Architectural Thoughts
To choose the right technology stack and design patterns, I carefully analyzed the real-world operational requirements of a transit barrier network.

The application is split into two main components: a Frontend (representing the physical ticket gates) and a Backend. Currently, my focus is entirely on developing the backend architecture.

My architectural approach is built around Edge Computing and decentralized decision-making:

**The Challenge (Latency):** The physical ticket gates must decide whether to open within a fraction of a second (ideally under 100–200 milliseconds) to maintain a continuous flow of passengers. A synchronous, blocking network request to a central backend would be far too slow and prone to network latency.

**The Solution (Optimistic Edge Decisions):** The gates function as intelligent edge devices. They read the data stored directly on the physical IC card, perform local validation and tariff calculations, and make an optimistic decision to open the barrier autonomously.

**The Backend as an Auditor:** The backend is removed from the critical path of the physical entry/exit process. Instead, it serves as the Auditor and Single Source of Truth. It asynchronously ingests transaction streams from the gates, reconciles account balances, detects anomalies (like card cloning), and distributes global hot-fixes or blocklists back to the edge.

To support this decoupled, high-throughput architecture, the backend requires a specialized data flow that ensures absolute auditability, historical accuracy, and high read performance.

**Data Consistency & History:** Event Sourcing
To allow the backend to audit the system effectively, several strict requirements must be met:

*Strict Ordering:* Events must be processed in the exact order they occurred to calculate correct fares and prevent state anomalies.

*Audit Trail:* We need a foolproof, immutable history of when and why a specific state change occurred.

To solve this, I chose Event Sourcing. Instead of just storing the current state of a card (e.g., its current balance), every single interaction is stored as an immutable sequence of events. The current state of any card can be reconstructed at any point in time by replaying these events.

**Performance & Queryability:** CQRS
While Event Sourcing is perfect for writing and auditing data, it introduces a challenge: reconstructing a card's current state (rehydrating) from hundreds of events is computationally expensive, and querying or filtering complex data (like listing active trips or card balances) directly from an Event Store is highly inefficient.

To overcome this, I implemented the CQRS (Command Query Responsibility Segregation) pattern:

The Write Path (Commands): Handles incoming state changes and appends them to the Event Store.

The Read Path (Queries): Uses optimized, denormalized read models tailored for fast fetching, filtering, and reporting.

Domain-Driven Design (DDD)
To keep the business logic clean, maintainable, and aligned with real-world transit rules, the backend is designed using Domain-Driven Design (DDD) principles. The system is divided into distinct Bounded Contexts to enforce clear boundaries and keep the codebase highly modular.

**In Short**
I am leveraging DDD, CQRS, and Event Sourcing to build a highly resilient, audit-safe, and performant backend that elegantly solves the real-world complexities of a modern transit network.

# Future Roadmap & What’s Next
Currently, the Journey and Wallet Bounded Contexts are actively being developed. I have completed the core Domain and Application layers for both contexts, ensuring the business logic and core use cases are fully established.

Moving forward, the immediate next steps and future enhancements include:

Infrastructure & API Layers: Implementing the database adapters, event store integration, and the API endpoints (gRPC/REST) to make the backend fully operational.

Additional Bounded Contexts: Expanding the ecosystem to support real-world IC card use cases beyond transit, such as a Retail/Point of Sale (POS) context to simulate paying at convenience stores (Konbinis) or vending machines.

# Disclaimer
This project is a personal journey to master the Go (Golang) programming language and explore advanced backend architectures.

To maximize my productivity, I leveraged generative AI as a collaborative assistant. However, to ensure a deep learning experience, I enforced a strict division of labor:

The Architecture & Design: Entirely self-designed. I defined the bounded contexts, the CQRS flow, the Event Sourcing mechanics, and the overall software design.

The Implementation: I used AI to handle repetitive boilerplate code and to accelerate my understanding of Go-specific idiomatic patterns, syntax, and best practices.