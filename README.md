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

---

## AI Assistance & Collaboration Disclosure

In accordance with the technical test guidelines regarding the use of AI:

* **Docker Infrastructure:** I used **Claude (LLM)** to assist in composing the `docker-compose.yml` and configuring the `odyssey.conf`.
* **Dockerfile.odyssey:** Due to the absence of an official/stable Odyssey image on Docker Hub, Claude helped in writing a multi-stage `Dockerfile` to compile Odyssey from source, ensuring a secure and optimized build.
* **Decision Making:** All architectural choices (Master-Slave, the switch from PgPool to Odyssey, and the repository pattern) were made by me based on professional experience. AI was used as a productivity tool to accelerate boilerplate generation.

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