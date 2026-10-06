# Talent Profile Intelligence

Talent Profile Intelligence is the currently implemented backend service in TalentIQ. It manages talent profile records through a Go/Gin REST API, persists them in PostgreSQL using pgx, and includes a Kafka publisher and consumer for local development.

## Implemented behavior

- Create, retrieve, list, update, and deactivate talent profiles.
- Validate required profile fields and reject negative experience values.
- Enforce unique employee codes and email addresses in PostgreSQL.
- Paginate profile lists (`page` defaults to `1`; `limit` defaults to `20` and is capped at `100`).
- Publish a `candidate.profile.updated` Kafka event after successful profile updates.
- Consume profile event messages and log their contents. Event-specific processing is not implemented yet.

Profile creation and deactivation do not currently publish Kafka events. If publishing an update event fails, the service logs the failure; the completed database update still returns success.

## Stack and structure

- Go 1.25.5
- Gin HTTP router
- PostgreSQL with pgx and SQL migrations (no ORM)
- Kafka with `segmentio/kafka-go`
- Structured logging with `log/slog`

```text
cmd/api/          Application entry point and route wiring
internal/config/  Environment-based configuration
internal/handler/ HTTP handlers and health endpoints
internal/dto/     Request and response types
internal/service/ Business validation and application logic
internal/repository/ PostgreSQL queries
internal/model/   Domain model
internal/messaging/ Kafka producer, consumer, and event contract
migrations/       PostgreSQL schema
```

## Local setup

Run the following from this service directory. You need Go, Docker, PostgreSQL, and Kafka.

### 1. Start PostgreSQL

The default connection settings are `localhost:5432`, database `talent_db`, user `postgres`, and password `postgres`.

```powershell
docker run --name talentiq-postgres `
  -e POSTGRES_USER=postgres `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=talent_db `
  -p 5432:5432 `
  -d postgres:17
```

If the container already exists, start it with `docker start talentiq-postgres`.

Apply the initial schema:

```powershell
Get-Content migrations/001_create_talent_profiles.sql |
  docker exec -i talentiq-postgres psql -U postgres -d talent_db
```

### 2. Start Kafka

The included Compose file starts a local Kafka broker on `localhost:9092`:

```powershell
docker compose up -d kafka
```

The configured topic is `talent-profile-events`; Kafka is configured to create it automatically for local development.

### 3. Configure and run

Start the API:

```powershell
go run ./cmd/api
```

Configuration is read directly from process environment variables. The application does not load `.env` automatically. The checked-in `.env.example` is a list of variable names/placeholders; use actual values in your shell or deployment environment.

| Variable | Default |
|---|---|
| `APP_NAME` | `talent-profile-intelligence` |
| `APP_ENV` | `development` |
| `SERVER_PORT` | `8080` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `postgres` |
| `DB_PASSWORD` | `postgres` |
| `DB_NAME` | `talent_db` |
| `KAFKA_BROKERS` | `localhost:9092` |
| `KAFKA_TOPIC` | `talent-profile-events` |
| `KAFKA_GROUP_ID` | `talent-profile-intelligence` |

Both PostgreSQL and Kafka clients are initialized during startup. `/health` reports process health; `/ready` checks the PostgreSQL connection.

## HTTP API

All profile routes are under `/api/v1`. UUIDs are used for profile IDs.

| Method | Route | Result |
|---|---|---|
| `GET` | `/health` | Liveness response |
| `GET` | `/ready` | PostgreSQL readiness response |
| `POST` | `/api/v1/talent-profiles` | Create a profile; new profiles are `ACTIVE` |
| `GET` | `/api/v1/talent-profiles` | Return a paginated profile list |
| `GET` | `/api/v1/talent-profiles/:id` | Retrieve one profile |
| `PUT` | `/api/v1/talent-profiles/:id` | Update profile details and status |
| `DELETE` | `/api/v1/talent-profiles/:id` | Deactivate a profile (soft delete) |

List example: `GET /api/v1/talent-profiles?page=2&limit=20`. The response includes `data`, `page`, `limit`, and `total`.

Create request fields:

```json
{
  "employee_code": "EMP-1001",
  "first_name": "Asha",
  "last_name": "Rao",
  "email": "asha.rao@example.com",
  "designation": "Software Engineer",
  "department": "Engineering",
  "total_experience_years": 4.5,
  "phone": "+91-555-0100",
  "location": "Bengaluru",
  "summary": "Backend engineer"
}
```

`employee_code`, `first_name`, `last_name`, `email`, `designation`, and `department` are required by service validation. `total_experience_years` must be non-negative. `phone`, `location`, and `summary` are optional. Update accepts the editable profile fields plus `profile_status`; `employee_code` is not editable.

Success responses have `success`, `message`, and `data` fields. Errors have `success: false` and a `message`. Invalid bodies or IDs return `400`, duplicate employee codes return `409`, missing profiles return `404`, and unexpected failures return `500`.

## Kafka event contract

Successful profile updates publish an event to the configured topic, keyed by profile UUID:

```json
{
  "event_id": "<event UUID>",
  "event_type": "candidate.profile.updated",
  "talent_id": "<profile UUID>",
  "employee_code": "EMP-1001",
  "occurred_at": "<UTC timestamp>"
}
```

The consumer currently unmarshals this contract and logs received events. It skips malformed messages after logging them. Broker settings are local service configuration and have not been moved to centralized TalentIQ infrastructure.

## Development and tests

```powershell
gofmt -w .
go test ./...
go build ./...
```

Unit tests for service behavior are in `internal/service`. Repository integration tests are not currently implemented; `tests/integration` and `tests/unit` are placeholders.

Keep HTTP concerns in handlers, business rules in services, and SQL in repositories. Do not commit `.env` files or secrets. Update this README when service behavior, configuration, or API contracts change.
