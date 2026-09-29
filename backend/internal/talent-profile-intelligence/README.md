# Talent Profile Intelligence

Talent Profile Intelligence is a backend microservice in the TalentIQ platform.

The service is responsible for managing the core talent profile information
used by the TalentIQ platform.

## Current Status

Implementation is currently in progress.

Completed lessons:

- TPI-001: Service Foundation
- TPI-002: PostgreSQL Database Setup
- TPI-003: Talent Profile Model, Migration and Structured Logging

## Technology Stack

- Go
- Gin
- PostgreSQL
- pgx
- Docker

Planned technologies:

- Kafka
- Redis
- OpenAPI / Swagger
- Prometheus
- OpenTelemetry
- Docker/Kubernetes

## Current Architecture

```text
HTTP Client
    |
    v
Gin Router
    |
    v
Talent Profile Intelligence
    |
    +------> PostgreSQL
    |
    +------> Application Logs

## Repository Layer

The repository layer is responsible for direct communication with PostgreSQL.

The Talent Profile repository provides:

- Create talent profile
- Get talent profile by ID
- List talent profiles
- Update talent profile
- Soft delete talent profile

The application follows this dependency flow:

```text
Handler
   |
   v
Service
   |
   v
Repository
   |
   v
PostgreSQL

## REST API

The Talent Profile Intelligence service exposes versioned REST APIs
under `/api/v1`.

### Talent Profile Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/talent-profiles` | Create a talent profile |
| GET | `/api/v1/talent-profiles` | List talent profiles |
| GET | `/api/v1/talent-profiles/:id` | Get a profile |
| PUT | `/api/v1/talent-profiles/:id` | Update a profile |
| DELETE | `/api/v1/talent-profiles/:id` | Deactivate a profile |

### Pagination

The list endpoint supports:

```text
GET /api/v1/talent-profiles?page=1&limit=20