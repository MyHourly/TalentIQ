# Talent Profile Intelligence Service

## 1. Overview

The Talent Profile Intelligence service is a backend microservice
within the TalentIQ platform.

The service manages the core profile information of talents and
provides the foundation for other TalentIQ capabilities that require
structured talent profile information.

The service is being implemented using Go and Gin and uses PostgreSQL
for structured profile data.

## 2. Current Responsibilities

The service is responsible for the core Talent Profile domain.

The current profile information includes:

- Employee code
- First name
- Last name
- Email
- Phone
- Designation
- Department
- Location
- Professional summary
- Total professional experience
- Profile status
- Audit timestamps

Additional profile-related capabilities will be introduced
incrementally.

## 3. Technology Stack

Current:

- Go
- Gin
- PostgreSQL
- pgx
- Docker
- Go `log/slog`

Planned:

- Kafka
- Redis
- OpenAPI / Swagger
- Prometheus
- OpenTelemetry
- Docker/Kubernetes

## 4. Architecture

The current service follows a layered microservice structure.

```text
HTTP Client
     |
     v
Gin Router
     |
     v
Handler
     |
     v
Service Layer
     |
     v
Repository
     |
     v
PostgreSQL