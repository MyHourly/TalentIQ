# TalentIQ

## Enterprise Talent Intelligence & Skill Graph Platform

TalentIQ is the team's office project for building an enterprise Talent Intelligence and Skill Graph platform. The platform is intended to centralize technical talent information and support talent discovery, skill intelligence, semantic search, skill-gap analysis, technology trends, and talent recommendations.

> **Working note:** This repository is the actual project repository. The earlier React frontend prototype was created only for product/UI understanding and is not automatically the production project.

---

## Current Team Module Allocation

| Team Member | Assigned Module | Feature Branch |
|---|---|---|
| Manjunatha R | Talent Profile | `feature/talent-profile` |
| Md | Certificate and Assessment | `feature/certificate-and-assessment` |
| Naveen | Skills | `feature/skills` |
| Anaiza Java | Skill Graph | `feature/skill-graph` |
| **Hari Babu** | **AI / Semantic Talent Search** | `feature/ai-semantic-talent-search` |
| Ashutosh Kumar | Skill Gap Analysis | `feature/skill-gap-analysis` |
| Kalyan Reddy | Technology Trends | `feature/technology-trends` |
| Charan G | Talent Recommendations | `feature/talent-recommendations` |
| Yashodhar G K | Dashboard | `feature/dashboard` |

---

## Technology Direction

### Backend

- Golang
- Gin
- PostgreSQL
- Redis
- Apache Kafka
- RabbitMQ
- Vector Database

### Frontend

- React
- TypeScript

### Platform / Observability

- Docker
- Kubernetes
- Prometheus
- Grafana
- OpenTelemetry
- OpenAPI / Swagger
- CI/CD

Final technology choices that are not yet confirmed should be agreed by the team/manager before implementation.

---

## High-Level Architecture

```text
Existing Company Systems
          |
          v
    Integration Layer
          |
          v
       Apache Kafka
          |
          v
 Talent / Skill / Project Services
          |
          v
   Intelligence / Analytics
       /          \
      v            v
 PostgreSQL     Vector DB
       \            /
        v          v
       Search & Recommendation
                |
                v
             Gin APIs
                |
                v
          React Dashboard
```

### Kafka vs RabbitMQ

**Kafka** is used for business/event streaming between services.

**RabbitMQ** is used for reliable asynchronous/background processing such as resume/document processing, embedding generation, skill extraction, recommendation jobs, retries and dead-letter handling.

---

# Git Branch Strategy

We use a simple three-level workflow:

```text
main
  |
  |-- stable / approved code
  |
  v
develop
  |
  |-- integration branch
  |
  +-- feature/talent-profile
  +-- feature/certificate-and-assessment
  +-- feature/skills
  +-- feature/skill-graph
  +-- feature/ai-semantic-talent-search
  +-- feature/skill-gap-analysis
  +-- feature/technology-trends
  +-- feature/talent-recommendations
  +-- feature/dashboard
```

### `main`

- Stable/approved code.
- No normal direct pushes.
- Changes reach `main` through reviewed pull requests.

### `develop`

- Team integration branch.
- Completed module work is merged here for integration and testing.

### `feature/*`

- Individual module development branches.
- Each member primarily works on their assigned module branch.

---

# Standard Development Workflow

## 1. Clone the repository

```bash
git clone https://github.com/<OWNER_OR_ORGANIZATION>/TalentIQ.git
cd TalentIQ
```

## 2. Fetch the latest branches

```bash
git fetch --all
```

## 3. Switch to your module branch

Example for Hari Babu:

```bash
git checkout feature/ai-semantic-talent-search
```

## 4. Update your branch before starting work

```bash
git pull origin feature/ai-semantic-talent-search
```

## 5. Make and test your changes

Keep changes focused on your module.

## 6. Review your changes

```bash
git status
git diff
```

## 7. Stage changes

```bash
git add .
```

## 8. Commit

```bash
git commit -m "feat: add semantic talent search"
```

## 9. Push

```bash
git push origin feature/ai-semantic-talent-search
```

## 10. Open a Pull Request

Create a PR from your feature branch into `develop`.

---

# Pull Request Flow

```text
feature/<module>
       |
       v
 Pull Request
       |
       v
   develop
       |
       v
Integration Testing
       |
       v
Review / Approval
       |
       v
     main
```

Do not merge individual feature work directly to `main` during normal development.

---

# Commit Message Convention

Use a consistent format:

```text
<type>: <short description>
```

Examples:

```text
feat: add semantic talent search
feat: add skill graph visualization
fix: handle empty search results
fix: correct skill gap calculation
refactor: simplify search service
test: add semantic search tests
docs: update API documentation
chore: update dependencies
```

Recommended types:

- `feat` — new functionality
- `fix` — bug fix
- `refactor` — code restructuring without behavior change
- `test` — tests
- `docs` — documentation
- `chore` — maintenance/tooling

---

# Team Working Rules

1. Work only in your assigned module branch unless the team agrees otherwise.
2. Pull the latest changes before beginning work.
3. Do not commit passwords, API keys, tokens or `.env` files.
4. Keep commits focused and understandable.
5. Test your changes before opening a Pull Request.
6. Update documentation when an API, event, data model or workflow changes.
7. Coordinate API contracts with the corresponding frontend/backend member before integration.
8. For Kafka/RabbitMQ changes, document the event/queue name, payload, producer/consumer and failure/retry behavior.
9. Avoid unrelated changes in another member's module.
10. Do not force-push shared branches unless the team lead explicitly asks you to.

---

# Frontend / Backend Contract

For every module, frontend and backend members should agree on:

- Endpoint
- HTTP method
- Request schema
- Response schema
- IDs and field names
- Pagination
- Sorting/filtering
- Status values
- Error response format

For event-driven functionality, also agree on:

- Kafka event/topic name
- Producer
- Consumer(s)
- Event payload
- Delivery assumptions
- Idempotency
- Retry/DLQ behavior

---

# Module Summary

## Talent Profile

Centralized talent profile, experience, projects, certifications, assessments and skills.

## Certificate and Assessment

Certification records, assessment history/results and related profile intelligence.

## Skills

Skill catalog, categories, proficiency, relationships and skill intelligence.

## Skill Graph

Relationships between people, skills, technologies, projects, certifications, domains and experience.

## AI / Semantic Talent Search

Natural-language talent discovery using semantic similarity, vector retrieval, filters and ranking.

## Skill Gap Analysis

Current capability vs expected demand and identification of organizational gaps.

## Technology Trends

Technology demand, growth/decline, emerging technologies and trend analysis.

## Talent Recommendations

Candidate recommendations based on skills, experience, projects, certifications, assessments and semantic similarity.

## Dashboard

Unified project overview, summary metrics, charts, recent activity and navigation into the core modules.

---

# Hari Babu — AI / Semantic Talent Search

The initial responsibility for the Semantic Talent Search module is:

### Frontend

- Natural-language search interface
- Search filters
- Search result list/cards
- Relevance/similarity presentation
- Navigation to talent profiles

### Backend

- Search API
- Embedding generation integration
- Vector DB retrieval
- Semantic similarity search
- Search/ranking logic
- Indexing pipeline integration

Expected conceptual flow:

```text
User Query
    |
    v
Gin Search API
    |
    v
Search Service
   / \
  v   v
Vector DB  PostgreSQL
   |
   v
Semantic Retrieval
   |
   v
Ranking / Filtering
   |
   v
Relevant Talent Profiles
```

The final Vector DB and embedding model are to be confirmed by the project team.

---

# Repository Structure (Target)

```text
TalentIQ/
├── README.md
├── CONTRIBUTING.md
├── .gitignore
├── docs/
│   ├── architecture/
│   ├── api/
│   └── modules/
├── frontend/
├── backend/
├── infrastructure/
│   ├── docker/
│   ├── kafka/
│   ├── rabbitmq/
│   └── kubernetes/
└── .github/
    └── workflows/
```

This is a target structure. Do not create empty directories just for appearance; add them when the corresponding work begins.

---

# Project Governance

- `main` is the protected stable branch.
- `develop` is the integration branch.
- Feature branches are owned by module developers.
- Pull Requests are the normal route for integrating changes.
- The team should agree who performs final approvals and who owns the eventual `main` branch release/merge responsibility.
- The project manager's confirmed module assignments override this working document if they change.

---

# Important Note

This is an office project. Prioritize correctness, maintainability, security, testing, observability, documentation and clear engineering decisions over rapidly generating large amounts of code.
