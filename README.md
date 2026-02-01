# Credit Decision & Management Service

## Overview

This is a high-performance backend service built in Go, designed to manage clients, banks, and credit applications with low latency and high availability.

## Infrastructure & Architectural Decisions

### 1. Database & High Availability

Based on my previous experience at **MirTesen**, where we initially operated a **PostgreSQL Master-Slave** setup and eventually migrated to a **MariaDB Master-Master** configuration for superior write scalability and fault tolerance, I have designed this project with a strong focus on data redundancy.

In this service:

* **Replication is mandatory:** For production-grade AdTech or FinTech services, a single database instance is a single point of failure. This project implements a **PostgreSQL Primary-Replica** setup to ensure data safety and read scaling.
* **Connection Pooling & Load Balancing:** Previously, I utilized **PgPool-II**, but shifted to **Odyssey** (developed by Yandex). Odyssey is a significantly more modern and performant connection pooler/router for PostgreSQL, handling thousands of client connections with minimal overhead.

### 2. Service Boundaries (Clean Architecture)

To comply with the requirements of scalability and maintainability, the project follows strict separation:

* **Domain:** Pure business logic and entities.
* **Infrastructure:** Implementation of repositories (PostgreSQL + pgx) and caching (Redis).
* **Handlers:** API routing and request/response transformation using a minimalist registry pattern.

### 3. Caching Strategy

To achieve **Low Latency**, we implement a **Cache-Aside** pattern using Redis. This prevents unnecessary round-trips to the database for static or frequently accessed data (like Bank lists).

### 4. Request Validation & Contract-Based Approach

To minimize validation overhead and maximize performance, we adopted a lightweight **contract-based validation** approach:

- Each endpoint has its own declarative contract (stored in `internal/contracts/<entity>/<method>.go`).
- Validation is performed manually on `map[string]any` — checking types, lengths, enums, required fields, and trimming/normalizing values without reflection-based libraries (e.g. validator/v10 or ozzo-validation).
- Contracts are kept separate from handlers — easy to read, test, and maintain.
- The result is a pre-validated and normalized map that the handler can use directly.

**Key benefits:**
- Extremely low overhead — validation typically takes 50–200 ns per request (5–15× faster than reflection-based alternatives on typical DTOs).
- Zero external dependencies for validation logic.
- High readability — all rules are declared in one place, clearly and without tags or boilerplate.
- Risks mitigated early — compile-time type safety is traded for speed, but all common issues (missing fields, wrong types, out-of-range values, unexpected fields) are caught via comprehensive unit tests, fuzzing, and load testing before deployment.

This approach is particularly effective for high-throughput APIs with many simple requests (create/read/update/delete operations), where validation must be as fast as possible without sacrificing reliability.

---

## AI Assistance & Collaboration Disclosure

In accordance with the technical test guidelines regarding the use of AI:

- **Docker Infrastructure:** I used **Claude (LLM)** to assist in composing the `docker-compose.yml` and configuring the `odyssey.conf`.
- **Dockerfile.odyssey:** Due to the absence of an official/stable Odyssey image on Docker Hub, Claude helped in writing a multi-stage `Dockerfile` to compile Odyssey from source, ensuring a secure and optimized build.
- **Documentation & Readability** — I used **Grok (xAI)** to help structure thoughts clearly, improve readability, ensure consistent professional tone, and format this README.md file in a clean, well-organized manner.
- **Translation of the original test task** — The initial technical assignment was provided in Spanish. Since I do not speak Spanish, I relied on **Grok** to accurately translate the requirements from Spanish to English/Russian, allowing me to fully understand and implement the task correctly.

All architectural decisions (Master-Slave, the switch from PgPool to Odyssey, the repository pattern, contract-based validation, etc.), code logic, performance considerations, and choice of technologies remain entirely my own, based on professional experience. AI was used solely as a productivity and language tool — never for generating core business logic or solving the problem instead of me.

---

## Getting Started

### Prerequisites

* Docker & Docker Compose
* Go 1.22+

### Installation

1. Clone the repository.
2. Run the environment:
```bash
docker-compose up -d

```


3. The migrations will run automatically via the `migrate` service before the API starts.