# Talent Profile Intelligence

Talent Profile Intelligence is a backend microservice of the TalentIQ
platform. It manages core talent profile information through REST APIs
and PostgreSQL.

---

## 1. Current Status

### Completed

- TPI-001: Service Foundation
- TPI-002: PostgreSQL Setup
- TPI-003: Talent Profile Model, Migration & Logging
- TPI-004: DTO & Service Layer
- TPI-005: Repository & REST API
- TPI-006: Postman API Testing
- TPI-007: Service Unit Testing

### Next

- TPI-008: Repository Integration Testing

---

## 2. Technology Stack

- Go
- Gin
- PostgreSQL
- pgx
- Docker
- Postman
- Go Testing

Planned later:

- Kafka
- Redis
- Swagger/OpenAPI
- Prometheus
- OpenTelemetry
- Kubernetes

---

## 3. Architecture

```text
Client
  |
  v
Gin Router
  |
  v
Handler
  |
  v
DTO
  |
  v
Service
  |
  v
Repository
  |
  v
PostgreSQL
```

### Layer Responsibilities

| Layer | Responsibility |
|---|---|
| Handler | HTTP request/response |
| DTO | API request/response structure |
| Service | Business logic & validation |
| Repository | PostgreSQL & SQL queries |
| PostgreSQL | Data storage |

The repository uses **SQL + pgx** instead of an ORM.

---

## 4. PostgreSQL Setup

### Create PostgreSQL Container

```powershell
docker run --name talentiq-postgres `
  -e POSTGRES_USER=postgres `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=talent_db `
  -p 5432:5432 `
  -d postgres:17
```

If the container already exists:

```powershell
docker start talentiq-postgres
```

Check status:

```powershell
docker ps
```

### Database Details

```text
Host: localhost
Port: 5432
Database: talent_db
User: postgres
```

---

## 5. Dependencies

Install pgx:

```powershell
go get github.com/jackc/pgx/v5
```

Synchronize dependencies:

```powershell
go mod tidy
```

---

## 6. Database Migration

Migration file:

```text
migrations/001_create_talent_profiles.sql
```

### Windows PowerShell

```powershell
Get-Content migrations/001_create_talent_profiles.sql |
docker exec -i talentiq-postgres psql -U postgres -d talent_db
```

Use this command to apply the SQL migration to the PostgreSQL
container.

To access PostgreSQL manually:

```powershell
docker exec -it talentiq-postgres psql -U postgres -d talent_db
```

---

## 7. Run the Service

From the service directory:

```powershell
go run .
```

Health check:

```text
GET /health
```

---

## 8. REST APIs

Base path:

```text
/api/v1
```

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/api/v1/talent-profiles` | Create profile |
| GET | `/api/v1/talent-profiles` | List profiles |
| GET | `/api/v1/talent-profiles/:id` | Get profile |
| PUT | `/api/v1/talent-profiles/:id` | Update profile |
| DELETE | `/api/v1/talent-profiles/:id` | Deactivate profile |

### Pagination

```text
GET /api/v1/talent-profiles?page=1&limit=20
```

---

## 9. Testing

### Run Service Unit Tests

```powershell
go test ./internal/service/...
```

**Use:** Tests Service business logic using the mock repository.

### Run All Tests

```powershell
go test ./...
```

**Use:** Verifies that all project tests pass before committing.

### Verbose Tests

```powershell
go test -v ./...
```

**Use:** Shows detailed test execution when debugging failures.

### Check Coverage

```powershell
go test ./internal/service/... -cover
```

**Use:** Shows Service layer test coverage.

### Generate Coverage Report

```powershell
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
start coverage.html
```

**Use:** Opens a visual report showing which code is covered by tests.

---

## 10. Code Verification

Format code:

```powershell
gofmt -w .
```

Build the project:

```powershell
go build ./...
```

Recommended before committing:

```powershell
gofmt -w .
go test ./...
go build ./...
```

---

## 11. Git Workflow

Update local `develop`:

```powershell
git checkout develop
git pull origin develop
```

Create feature branch:

```powershell
git checkout -b feature/<feature-name>
```

Check changes:

```powershell
git status
git diff
```

Commit:

```powershell
git add .
git commit -m "type(scope): description"
```

Push:

```powershell
git push origin feature/<feature-name>
```

---

## 12. Development Rules

- Keep HTTP logic in Handler.
- Keep business logic in Service.
- Keep SQL/database logic in Repository.
- Use DTOs for API requests/responses.
- Add tests for new business logic.
- Do not commit `.env` or secrets.
- Run tests before pushing changes.
- Update documentation when API or architecture changes.

---

## 13. Lesson Tracking

| ID | Lesson | Status |
|---|---|---|
| TPI-001 | Service Foundation | Completed |
| TPI-002 | PostgreSQL Setup | Completed |
| TPI-003 | Model, Migration & Logging | Completed |
| TPI-004 | DTO & Service Layer | Completed |
| TPI-005 | Repository & REST API | Completed |
| TPI-006 | Postman Testing | Completed |
| TPI-007 | Service Unit Testing | Completed |
| TPI-008 | Repository Integration Testing | Next |

---

## 14. Current Development Flow

```text
Talent Profile API
       |
       v
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

Testing:

Service --> Mock Repository       (TPI-007)
Repository --> PostgreSQL         (TPI-008)
```

The next implementation step is **TPI-008: Repository Integration
Testing**, where the actual SQL repository will be tested against
PostgreSQL.