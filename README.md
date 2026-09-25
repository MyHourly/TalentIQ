# TalentIQ

### Enterprise Talent Intelligence & Skill Graph Platform

**Repository:** [MyHourly/TalentIQ](https://github.com/MyHourly/TalentIQ)
**Primary development branch:** `develop`
**Stable / approved branch:** `main`

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Project Goals](#2-project-goals)
3. [Current Seven-Module Architecture](#3-current-seven-module-architecture)
4. [Detailed Module Responsibilities](#4-detailed-module-responsibilities)
   - [4.1 Talent Profile Intelligence](#41-talent-profile-intelligence)
   - [4.2 Skill Intelligence](#42-skill-intelligence)
   - [4.3 Skill Graph](#43-skill-graph)
   - [4.4 AI / Semantic Talent Search](#44-ai--semantic-talent-search)
   - [4.5 Skill Gap Analysis](#45-skill-gap-analysis)
   - [4.6 Technology Trend Analysis](#46-technology-trend-analysis)
   - [4.7 Talent Recommendation Engine](#47-talent-recommendation-engine)
5. [Frontend Architecture](#5-frontend-architecture)
6. [Backend Architecture](#6-backend-architecture)
7. [High-Level System Architecture](#7-high-level-system-architecture)
8. [Cross-Module Dependencies](#8-cross-module-dependencies)
9. [Module Boundary Rules](#9-module-boundary-rules)
10. [Frontend / Backend API Contract](#10-frontend--backend-api-contract)
11. [Technology Direction](#11-technology-direction)
12. [Shared Infrastructure](#12-shared-infrastructure)
13. [Standard Development Workflow](#13-standard-development-workflow)
14. [Git Branch Strategy](#14-git-branch-strategy)
15. [Git Setup / Clone Instructions](#15-git-setup--clone-instructions)
16. [Fetching Branches](#16-fetching-branches)
17. [Switching / Creating a Feature Branch](#17-switching--creating-a-feature-branch)
18. [Pulling Latest Changes](#18-pulling-latest-changes)
19. [Commit Workflow](#19-commit-workflow)
20. [Push Workflow](#20-push-workflow)
21. [Pull Request Flow](#21-pull-request-flow)
22. [Pull Request Checklist](#22-pull-request-checklist)
23. [Commit Message Conventions](#23-commit-message-conventions)
24. [Team Working Rules](#24-team-working-rules)
25. [Testing Strategy](#25-testing-strategy)
26. [Definition of Done](#26-definition-of-done)
27. [Error Handling](#27-error-handling)
28. [Logging & Observability](#28-logging--observability)
29. [Security Rules](#29-security-rules)
30. [AI / Semantic Search Engineering Guidelines](#30-ai--semantic-search-engineering-guidelines)
31. [Database Ownership](#31-database-ownership)
32. [Messaging / Events](#32-messaging--events)
33. [API Versioning](#33-api-versioning)
34. [Documentation Rules](#34-documentation-rules)
35. [Governance](#35-governance)
36. [Open Decisions / To Be Confirmed](#36-open-decisions--to-be-confirmed)
37. [Current Repository Structure](#37-current-repository-structure)
38. [Quick Git Reference](#38-quick-git-reference)
39. [Final Engineering Workflow](#39-final-engineering-workflow)
40. [Project Principles](#40-project-principles)

---

## 1. Project Overview

TalentIQ is an **Enterprise Talent Intelligence & Skill Graph Platform**. It is designed to help an organization understand, search, and reason about its technical talent — who has which skills, how those skills relate to one another and to technology trends, where skill gaps exist, and which people are the best fit for a given need.

The platform exists to bring together information that is normally scattered across profiles, experience history, projects, certifications, and skills lists, and turn it into something that can be **queried, searched semantically, and reasoned about**, rather than just stored.

### What TalentIQ addresses

- Talent information (profiles, experience, projects, certifications, assessments) that is hard to search and compare at scale.
- Skills recorded as flat text rather than structured, related data.
- No structured way to see how skills, people, and technologies relate to one another.
- Difficulty finding the right person for a need using natural-language / semantic queries rather than rigid keyword filters.
- No systematic way to see where skill gaps exist against requirements.
- No systematic way to track technology trends against the organization's current skill base.
- No systematic, explainable way to generate talent recommendations.

### Main capabilities

- **Talent Profile Intelligence** — a structured, authoritative record of a person's profile, experience, projects, certifications, and assessments.
- **Skill Intelligence** — structured skill definitions and skill-level intelligence used across the platform.
- **Skill Graph** — a graph of relationships between skills, talent, and technologies.
- **AI / Semantic Talent Search** — natural-language and semantic search over talent data.
- **Skill Gap Analysis** — comparison of required skills against current skill inventory.
- **Technology Trend Analysis** — tracking of technology and skill demand signals.
- **Talent Recommendation Engine** — candidate-fit ranking and recommendations.

### Enterprise talent intelligence concept

"Talent intelligence" in TalentIQ means treating talent data as a connected, queryable asset rather than a collection of static documents. Profile, skill, and relationship data feed into search, gap analysis, trend analysis, and recommendations, instead of each living in isolation.

### Skill graph concept

The skill graph is the structural backbone connecting people, skills, and technologies through relationships, so other modules (search, recommendations) can reason about *how* things relate, not just whether a keyword matches.

### Intelligent search / recommendation concept

Rather than relying purely on exact-match filtering, TalentIQ combines semantic search (via AI / Semantic Talent Search) with the skill graph and skill intelligence to interpret queries, retrieve relevant candidates, and rank them — and the Talent Recommendation Engine builds on that same foundation to proactively surface fit-based recommendations.

---

## 2. Project Goals

TalentIQ's goals, as reflected in the current seven-module structure, are:

| Goal | Description |
|---|---|
| **Talent intelligence** | Maintain a structured, authoritative view of each person's profile, experience, projects, certifications, and assessments. |
| **Skill intelligence** | Maintain structured, consistent skill definitions and skill-level data that other modules can rely on. |
| **Skill relationships** | Model relationships between skills, talent, and technologies through a skill graph. |
| **Semantic talent discovery** | Allow talent to be discovered through semantic / natural-language search rather than rigid keyword matching. |
| **Skill gap analysis** | Identify gaps between required skills and the organization's current skill inventory. |
| **Technology trends** | Track technology and skill demand signals relevant to the organization. |
| **Talent recommendations** | Generate candidate-fit rankings and recommendations grounded in profile, skill, graph, and search data. |
| **Enterprise technical talent understanding** | Bring the above together into a single platform for understanding technical talent at an enterprise level. |

> These goals are derived directly from the seven functional modules described in this document. Goals beyond what the seven modules support are not claimed here.

---

## 3. Current Seven-Module Architecture

TalentIQ's backend and frontend are organized around **exactly seven functional modules**. There are no additional top-level modules; in particular, *Certificate and Assessment*, *Dashboard*, *Talent Profile*, *Skills*, *Technology Trends*, and *Talent Recommendations* are **not** separate top-level modules under the current structure — see [4.1 Talent Profile Intelligence](#41-talent-profile-intelligence) for where certificate/assessment functionality lives.

| # | Module | Backend Owner | Frontend Owner | Responsibility Summary |
|---|---|---|---|---|
| 4.1 | Talent Profile Intelligence | Ashuthosh Kumar | Naveendra | Core person/profile record, experience, projects, certifications, and assessments |
| 4.2 | Skill Intelligence | Vikas | Naveendra | Skill definitions and skill-level intelligence used across the platform |
| 4.3 | Skill Graph | Ranganath Gowda | Anaiza | Relationships between skills, talent, and technologies |
| 4.4 | AI / Semantic Talent Search | Hari Babu | Anaiza | Semantic / natural-language search over talent data |
| 4.5 | Skill Gap Analysis | Manjunath | Yashodhar | Comparison of required skills against current skill inventory |
| 4.6 | Technology Trend Analysis | Indra Kalyan Reddy | Charan | Technology and skill demand trend signals |
| 4.7 | Talent Recommendation Engine | Lalsab | Charan | Candidate-fit ranking and recommendations |

> **Allocation status — working allocation, pending confirmation.** The module allocation above reflects the current working allocation as provided in the project material. It has **not** been finalized and should not be treated as a confirmed assignment until signed off. The source material also flags a possible discrepancy in the total number of people/assignments recorded (for example, whether the correct count is 10 or 11); this is tracked as an open item in [Section 36](#36-open-decisions--to-be-confirmed) rather than resolved here.

---

## 4. Detailed Module Responsibilities

This section documents each of the seven modules in detail: its purpose, what it owns, what it explicitly does not own, its key dependencies, and its frontend/backend scope. Boundaries below follow the ownership and dependency information provided in the project material; responsibilities not supported by that material are not claimed.

### 4.1 Talent Profile Intelligence

**Backend owner:** Ashuthosh Kumar · **Frontend owner:** Naveendra

**Purpose:** Maintain the authoritative, structured record of a person's professional profile.

**Owns:**
- Core person/profile record
- Experience history
- Projects
- Certifications
- Assessments

> Certificate and Assessment functionality is **part of Talent Profile Intelligence** and is **not** a separate top-level 4.x module. Any certificate- or assessment-related data, screens, or endpoints belong under this module.

**Does not own:**
- Skill definitions or skill taxonomy (owned by [Skill Intelligence](#42-skill-intelligence))
- Skill/talent/technology relationships (owned by [Skill Graph](#43-skill-graph))
- Search, ranking, or semantic retrieval (owned by [AI / Semantic Talent Search](#44-ai--semantic-talent-search))
- Gap calculations (owned by [Skill Gap Analysis](#45-skill-gap-analysis))
- Technology trend signals (owned by [Technology Trend Analysis](#46-technology-trend-analysis))
- Recommendation scoring (owned by [Talent Recommendation Engine](#47-talent-recommendation-engine))

**Key dependencies (consumers of this module's data):**
- → [Skill Graph](#43-skill-graph)
- → [AI / Semantic Talent Search](#44-ai--semantic-talent-search)
- → [Talent Recommendation Engine](#47-talent-recommendation-engine)

**Frontend scope:** Profile, experience, projects, certifications, and assessments UI, under `frontend/src/modules/talent-profile-intelligence/`.

**Backend scope:** Profile domain APIs and data under `backend/internal/talent-profile-intelligence/`.

### 4.2 Skill Intelligence

**Backend owner:** Vikas · **Frontend owner:** Naveendra

**Purpose:** Maintain structured skill definitions and skill-level intelligence that other modules rely on.

**Owns:**
- Skill definitions
- Skill taxonomy and related skill intelligence
- Skill data consumed by dependent modules

**Does not own:**
- Person/profile records ([Talent Profile Intelligence](#41-talent-profile-intelligence))
- Skill/talent/technology relationship graph ([Skill Graph](#43-skill-graph))
- Search, gap, trend, or recommendation logic that consumes skill data

**Key dependencies (consumers of this module's data):**
- → [Skill Graph](#43-skill-graph)
- → [AI / Semantic Talent Search](#44-ai--semantic-talent-search)
- → [Skill Gap Analysis](#45-skill-gap-analysis)
- → [Technology Trend Analysis](#46-technology-trend-analysis)
- → [Talent Recommendation Engine](#47-talent-recommendation-engine)

**Frontend scope:** Skill definition and skill-intelligence UI, under `frontend/src/modules/skill-intelligence/`.

**Backend scope:** Skill domain APIs and data under `backend/internal/skill-intelligence/`.

### 4.3 Skill Graph

**Backend owner:** Ranganath Gowda · **Frontend owner:** Anaiza

**Purpose:** Model and expose relationships between skills, talent, and technologies.

**Owns:**
- Skill ↔ talent ↔ technology relationships
- Graph structure and relationship/query behavior

**Does not own:**
- The underlying profile or skill records themselves (owned by Talent Profile Intelligence and Skill Intelligence respectively) — Skill Graph relates them, it does not redefine them
- Search ranking or recommendation scoring

**Key dependencies:**
- Depends on profile identifiers from [Talent Profile Intelligence](#41-talent-profile-intelligence)
- Depends on skill identifiers from [Skill Intelligence](#42-skill-intelligence)
- → Consumed by [AI / Semantic Talent Search](#44-ai--semantic-talent-search)
- → Consumed by [Talent Recommendation Engine](#47-talent-recommendation-engine)

**Frontend scope:** Graph/relationship visualization and interaction, under `frontend/src/modules/skill-graph/`.

**Backend scope:** Graph domain APIs and data under `backend/internal/skill-graph/`.

### 4.4 AI / Semantic Talent Search

**Backend owner:** Hari Babu · **Frontend owner:** Anaiza

**Purpose:** Provide semantic and natural-language search over talent data.

**Owns:**
- Search APIs
- Query interpretation
- Semantic retrieval
- Embedding integration
- Ranking of search results
- Indexing hooks
- Semantic similarity logic

**Does not own:**
- Core profile data ([Talent Profile Intelligence](#41-talent-profile-intelligence))
- Skill definitions ([Skill Intelligence](#42-skill-intelligence))
- Graph relationships ([Skill Graph](#43-skill-graph))
- Final recommendation scoring ([Talent Recommendation Engine](#47-talent-recommendation-engine)) — search surfaces and ranks candidates; it does not itself produce the recommendation output

**Data consumed by this module:**
- Profile data from [Talent Profile Intelligence](#41-talent-profile-intelligence)
- Skill data from [Skill Intelligence](#42-skill-intelligence)
- Relationship data from [Skill Graph](#43-skill-graph)

**Consumed by:**
- → [Talent Recommendation Engine](#47-talent-recommendation-engine)

**Vector DB role:** A vector database is the intended direction for embedding storage/retrieval in support of semantic search. **The specific Vector DB product has not been selected — To Be Confirmed** (see [Section 36](#36-open-decisions--to-be-confirmed)).

**Search boundaries:** This module is responsible for interpreting a query and returning ranked, relevant candidates; it does not own the underlying profile, skill, or graph data it searches over, and it does not perform final recommendation scoring.

**Frontend scope:** Search UI, query input, and results presentation, under `frontend/src/modules/ai-semantic-talent-search/`.

**Backend scope:** Search domain APIs and data under `backend/internal/ai-semantic-talent-search/`.

### 4.5 Skill Gap Analysis

**Backend owner:** Manjunath · **Frontend owner:** Yashodhar

**Purpose:** Compare required skills against the current skill inventory to identify gaps.

**Owns:**
- Skill inventory comparison logic
- Required-skills comparison
- Gap calculation

**Does not own:**
- Skill definitions themselves ([Skill Intelligence](#42-skill-intelligence))
- Graph relationships, search, or recommendation logic

**Key dependencies:**
- Depends on skill data from [Skill Intelligence](#42-skill-intelligence)

**Frontend scope:** Gap analysis UI, under `frontend/src/modules/skill-gap-analysis/`.

**Backend scope:** Gap analysis domain APIs and data under `backend/internal/skill-gap-analysis/`.

### 4.6 Technology Trend Analysis

**Backend owner:** Indra Kalyan Reddy · **Frontend owner:** Charan

**Purpose:** Track technology and skill demand trend signals.

**Owns:**
- Technology trend signals
- Demand signals
- Skill/technology trend data

**Does not own:**
- Skill definitions themselves ([Skill Intelligence](#42-skill-intelligence))
- Recommendation or search logic

**Key dependencies:**
- Depends on skill data from [Skill Intelligence](#42-skill-intelligence)

**Frontend scope:** Trend analysis UI, under `frontend/src/modules/technology-trend-analysis/`.

**Backend scope:** Trend analysis domain APIs and data under `backend/internal/technology-trend-analysis/`.

### 4.7 Talent Recommendation Engine

**Backend owner:** Lalsab · **Frontend owner:** Charan

**Purpose:** Produce candidate-fit ranking and recommendations.

**Owns:**
- Recommendation scoring
- Candidate-fit ranking, combining talent/profile information, skill requirements, semantic similarity, and skill graph context

**Does not own:**
- Core profile data, skill definitions, graph structure, or search ranking algorithms themselves — this module combines their outputs into recommendations rather than owning the underlying data or search mechanics

**Data consumed by this module:**
- Profile data from [Talent Profile Intelligence](#41-talent-profile-intelligence)
- Skill data from [Skill Intelligence](#42-skill-intelligence)
- Relationship data from [Skill Graph](#43-skill-graph)
- Search/similarity output from [AI / Semantic Talent Search](#44-ai--semantic-talent-search)

**Frontend scope:** Recommendation UI, under `frontend/src/modules/talent-recommendation-engine/`.

**Backend scope:** Recommendation domain APIs and data under `backend/internal/talent-recommendation-engine/`.

---

## 5. Frontend Architecture

TalentIQ uses a **single shared React application** rather than separate frontends per module. Each module owns its own pages/components within that application, while common concerns are shared across modules.

- **Module-specific pages/components** — each module's UI lives under its own folder in `frontend/src/modules/`.
- **Shared components** — reusable UI building blocks in `frontend/src/components/`.
- **Shared hooks** — reusable React hooks in `frontend/src/hooks/`.
- **Shared services** — API-consuming service layers in `frontend/src/services/`.
- **Shared types** — shared type definitions in `frontend/src/types/`.
- **Shared utilities** — general-purpose helpers in `frontend/src/utils/`.
- **API consumption** — the frontend consumes backend APIs through the shared services layer; it does not access databases or module internals directly.
- **No duplicated business rules** — business logic that belongs to a module's backend should not be re-implemented in the React layer; the frontend renders and orchestrates, the backend owns the rules.
- **Frontend ownership grouping** — frontend ownership currently groups multiple modules under the same owner (see [Section 3](#3-current-seven-module-architecture)): Naveendra (Talent Profile Intelligence, Skill Intelligence), Anaiza (Skill Graph, AI / Semantic Talent Search), Yashodhar (Skill Gap Analysis), Charan (Technology Trend Analysis, Talent Recommendation Engine).

**Frontend directory structure:**

```
frontend/
└── src/
    ├── components/
    ├── hooks/
    ├── modules/
    │   ├── ai-semantic-talent-search/
    │   ├── skill-gap-analysis/
    │   ├── skill-graph/
    │   ├── skill-intelligence/
    │   ├── talent-profile-intelligence/
    │   ├── talent-recommendation-engine/
    │   └── technology-trend-analysis/
    ├── services/
    ├── types/
    └── utils/
```

## 6. Backend Architecture

The backend is a **modular backend** with seven domain modules living under `internal/`, plus shared layers for the API surface, entry point, configuration, and database access.

- **api/** — API layer exposed to the frontend and other consumers.
- **cmd/** — application entry point(s).
- **config/** — configuration.
- **database/** — database access/shared data layer.
- **internal/** — the seven domain modules, each owning its own logic and data access within its boundary.

Each domain module under `internal/` is owned by a single backend owner (see [Section 3](#3-current-seven-module-architecture)) and is responsible for its own domain logic, consistent with the [Module Boundary Rules](#9-module-boundary-rules).

**Backend directory structure:**

```
backend/
├── api/
├── cmd/
├── config/
├── database/
└── internal/
    ├── ai-semantic-talent-search/
    ├── skill-gap-analysis/
    ├── skill-graph/
    ├── skill-intelligence/
    ├── talent-profile-intelligence/
    ├── talent-recommendation-engine/
    └── technology-trend-analysis/
```

## 7. High-Level System Architecture

```
+---------------------------+
|          Frontend          |
|    (React application)     |
+---------------------------+
              |
              v
+---------------------------+
|        Backend APIs        |
+---------------------------+
              |
              v
+---------------------------+
|  Seven Functional Modules  |
+---------------------------+
              |
              v
+---------------------------+
|   Shared Infrastructure    |
+---------------------------+
```

**Logical relationships between modules** (simplified for illustration — see [Section 8](#8-cross-module-dependencies) for the complete, authoritative dependency list):

```
Talent Profile Intelligence → Skill Graph → AI / Semantic Talent Search → Talent Recommendation Engine
Skill Intelligence → Skill Gap Analysis
Skill Intelligence → Technology Trend Analysis
```

This diagram shows logical flow and module relationships only; it does not represent a confirmed deployment topology (see [Section 11](#11-technology-direction) and [Section 36](#36-open-decisions--to-be-confirmed)).

## 8. Cross-Module Dependencies

The following logical dependencies exist between modules, per the current project structure. A dependency means the downstream module **consumes data/contracts** from the upstream module — it is not, by itself, a license for direct cross-module database access (see [Section 9](#9-module-boundary-rules)).

| Upstream module | Downstream module |
|---|---|
| Talent Profile Intelligence | → Skill Graph |
| Talent Profile Intelligence | → AI / Semantic Talent Search |
| Talent Profile Intelligence | → Talent Recommendation Engine |
| Skill Intelligence | → Skill Graph |
| Skill Intelligence | → AI / Semantic Talent Search |
| Skill Intelligence | → Skill Gap Analysis |
| Skill Intelligence | → Technology Trend Analysis |
| Skill Intelligence | → Talent Recommendation Engine |
| Skill Graph | → AI / Semantic Talent Search |
| Skill Graph | → Talent Recommendation Engine |
| AI / Semantic Talent Search | → Talent Recommendation Engine |

**What this means in practice:** each row represents a downstream module relying on data or contracts exposed by the upstream module (for example, through its API), not a shared database table or direct SQL access across module boundaries. See [Module Boundary Rules](#9-module-boundary-rules) for how these dependencies should be implemented.

## 9. Module Boundary Rules

- Each module **owns its domain** — its data, its business logic, and its API contract.
- Other modules **consume contracts**, not internal implementation details.
- **No direct cross-module database access** unless explicitly approved.
- **No duplicate ownership** — a given piece of data or logic has exactly one owning module.
- **Cross-module changes require coordination** between the owning backend/frontend owners.
- **Business logic belongs to the correct module** — do not implement another module's business rules locally to work around a dependency.
- **Shared infrastructure should remain shared** — infrastructure (see [Section 12](#12-shared-infrastructure)) is used in common, not duplicated per module.

## 10. Frontend / Backend API Contract

TalentIQ follows a **contract-first** approach: the API contract for a module is defined before frontend/backend implementation begins, and the frontend consumes only that contract.

Contracts should define:

- **Endpoint naming** — consistent, resource-oriented naming per module.
- **HTTP methods** — standard REST verbs (`GET`, `POST`, `PUT`/`PATCH`, `DELETE`) used consistently with their conventional meaning.
- **Request** — a defined, versioned request shape.
- **Response** — a defined, versioned response shape.
- **Error format** — a consistent error response shape (see [Section 27](#27-error-handling)).
- **Validation** — input validation at the API boundary.
- **Pagination** — where a response returns a collection.
- **Authentication / authorization** — where required by the endpoint. *(Specific auth provider/mechanism: To Be Confirmed — see [Section 36](#36-open-decisions--to-be-confirmed).)*
- **Versioning** — see [Section 33](#33-api-versioning).
- **API ownership** — each API is owned by the backend owner of the module it belongs to.

**Illustrative example only — not a confirmed TalentIQ API:**

```
GET /api/v1/{module}/{resource}

Response 200:
{
  "data": [ { "id": "string", "...": "..." } ],
  "pagination": { "page": 1, "pageSize": 20, "total": 0 }
}

Response 4xx/5xx:
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

## 11. Technology Direction

The following technology **directions** have been identified for TalentIQ. Items marked *To Be Confirmed* have not had a specific product/provider selected; no cloud provider, deployment platform, or containerization stack is claimed here, since none is supported by the current project material.

**Frontend:**
- React

**Backend:**
- Backend services exposing APIs consumed by the frontend (see [Section 6](#6-backend-architecture))

**Data / infrastructure directions:**
- PostgreSQL
- Redis
- Kafka
- RabbitMQ
- Vector DB *(specific product — To Be Confirmed)*

**AI / search directions:**
- Embeddings *(specific model/provider — To Be Confirmed)*
- Semantic retrieval
- Vector search
- Ranking
- Indexing

**Observability directions:**
- Logging
- Monitoring
- Observability

> These are directional technology areas identified in the project material, not a confirmed, finalized technology stack. See [Section 36](#36-open-decisions--to-be-confirmed).

## 12. Shared Infrastructure

Shared infrastructure is used in common across modules rather than duplicated per module (see [Module Boundary Rules](#9-module-boundary-rules)).

- PostgreSQL
- Redis
- Kafka
- RabbitMQ
- Vector DB
- Observability
- Logging
- Monitoring
- Security
- Deployment

**Confirmed project direction:** the areas above are the identified shared-infrastructure directions.
**Requiring confirmation:** specific product choices (e.g., Vector DB product, exact deployment approach), operational ownership of shared infrastructure, and specific security/secret-management tooling. See [Section 36](#36-open-decisions--to-be-confirmed).

## 13. Standard Development Workflow

Every change to TalentIQ should follow this standard flow:

```
Requirement
   ↓
Identify module
   ↓
Confirm owner
   ↓
Understand dependencies
   ↓
Define API / contract
   ↓
Create feature branch
   ↓
Implement
   ↓
Run tests
   ↓
Review changes
   ↓
Commit
   ↓
Push
   ↓
Create Pull Request
   ↓
Code review
   ↓
Merge to develop
   ↓
Integration testing
   ↓
Approval
   ↓
main
```

1. **Requirement** — a requirement or task is identified.
2. **Identify module** — determine which of the seven modules owns this requirement (see [Section 3](#3-current-seven-module-architecture)).
3. **Confirm owner** — confirm the backend/frontend owner responsible for that module.
4. **Understand dependencies** — check [Cross-Module Dependencies](#8-cross-module-dependencies) for upstream/downstream impact.
5. **Define API / contract** — define or confirm the API contract before implementation (see [Section 10](#10-frontend--backend-api-contract)).
6. **Create feature branch** — branch from `develop` (see [Section 14](#14-git-branch-strategy)).
7. **Implement** — implement the change within the owning module's boundary.
8. **Run tests** — run relevant tests (see [Section 25](#25-testing-strategy)).
9. **Review changes** — self-review the diff before committing.
10. **Commit** — commit using [conventional commit messages](#23-commit-message-conventions).
11. **Push** — push the feature branch to the remote.
12. **Create Pull Request** — open a PR into `develop`.
13. **Code review** — the PR is reviewed by the relevant owner(s).
14. **Merge to develop** — once approved, merge into `develop`.
15. **Integration testing** — the change is validated as part of `develop`.
16. **Approval** — release approval is given.
17. **main** — the change is promoted to `main`, the stable/approved branch.

---

## 14. Git Branch Strategy

| Branch | Purpose |
|---|---|
| `main` | Stable / approved branch |
| `develop` | Integration / development branch |
| `feature/*` | Individual development work |

**Flow:** `feature/*` → `develop` → `main`

A feature branch is created from `develop`, developed and reviewed via Pull Request into `develop`, validated through integration/testing, and — once approved — `develop` is promoted to `main`.

> Feature branch names below (e.g., `feature/talent-profile-intelligence`) are **naming conventions/examples**, not a claim that these branches currently exist in the repository.

**Example feature branch naming conventions** (one per module):

```
feature/talent-profile-intelligence
feature/skill-intelligence
feature/skill-graph
feature/ai-semantic-talent-search
feature/skill-gap-analysis
feature/technology-trend-analysis
feature/talent-recommendation-engine
```

## 15. Git Setup / Clone Instructions

```bash
# Clone the repository
git clone https://github.com/MyHourly/TalentIQ.git

# Enter the project directory
cd TalentIQ

# View all branches (local and remote)
git branch -a

# Switch to the develop branch
git checkout develop

# Pull the latest changes
git pull origin develop
```

## 16. Fetching Branches

```bash
# Fetch updates from the remote
git fetch origin

# Fetch all remotes and prune deleted remote-tracking branches
git fetch --all --prune

# List all branches (local and remote) after fetching
git branch -a
```

`git fetch` updates your local knowledge of the remote (new branches, new commits) **without** modifying your working directory or current branch — it is the safe way to see what has changed before deciding to `pull` or `checkout`.

## 17. Switching / Creating a Feature Branch

```bash
# Switch to an existing local branch
git checkout <branch-name>

# Switch to a branch that exists on the remote but not yet locally
git checkout -b <branch-name> origin/<branch-name>

# Create a new feature branch from an up-to-date develop
git checkout develop
git pull origin develop
git checkout -b feature/<feature-name>
```

## 18. Pulling Latest Changes

```bash
# Get the latest develop
git checkout develop
git pull origin develop
```

**Updating a feature branch:** once `develop` is up to date, bring those changes into your feature branch carefully — merge or rebase only as appropriate for the team's convention, and **never rewrite shared history** (for example, force-push after a rebase) on a branch others are also working on, without explicit agreement.

## 19. Commit Workflow

```bash
# See what has changed
git status

# Review the actual changes
git diff

# Stage changes
git add .

# Commit with a conventional commit message
git commit -m "feat(module-name): short description"
```

Always inspect `git status` and `git diff` before committing — this catches unrelated or unintended changes (stray files, debug code, accidental secrets) before they enter history.

## 20. Push Workflow

```bash
# First push of a new feature branch (sets upstream)
git push -u origin feature/<feature-name>

# Subsequent pushes
git push
```

## 21. Pull Request Flow

Pull Requests move code through two review gates:

**1. Feature branch → develop**

```
feature/<feature-name> → Pull Request → Code review → Checks / tests → Approval → develop
```

**2. develop → main**

```
develop → Integration / system testing → Release approval → main
```

## 22. Pull Request Checklist

- [ ] Requirement implemented
- [ ] Correct module (per [Section 3](#3-current-seven-module-architecture))
- [ ] API contract defined
- [ ] Tests passed
- [ ] Error handling in place
- [ ] Security reviewed
- [ ] No secrets committed
- [ ] Documentation updated
- [ ] No unrelated changes
- [ ] Screenshots included for UI changes
- [ ] Cross-module dependencies reviewed (see [Section 8](#8-cross-module-dependencies))

## 23. Commit Message Conventions

TalentIQ uses **conventional commit** style:

| Prefix | Use for |
|---|---|
| `feat:` | A new feature |
| `fix:` | A bug fix |
| `refactor:` | Code change that neither fixes a bug nor adds a feature |
| `docs:` | Documentation-only changes |
| `test:` | Adding or updating tests |
| `chore:` | Maintenance work (tooling, config, dependencies) |

**Examples:**

```
feat(skill-graph): add relationship query endpoint
fix(ai-semantic-talent-search): correct ranking tie-break logic
refactor(talent-profile-intelligence): simplify certification mapper
docs(readme): update module ownership table
test(skill-gap-analysis): add gap calculation unit tests
chore(deps): bump backend dependency versions
```

**Do not use meaningless commit messages**, such as:

```
update
changes
final
test
new code
```

---

## 24. Team Working Rules

- One backend owner per module.
- Frontend consumes backend contracts — it does not reimplement backend business logic.
- No direct cross-module database access.
- Define contracts before integration begins.
- Reuse shared infrastructure rather than duplicating it.
- Respect module boundaries (see [Section 9](#9-module-boundary-rules)).
- Coordinate cross-module changes with the relevant owners.
- Do not commit secrets.
- Keep Pull Requests focused on a single change.
- Document important architectural changes.

## 25. Testing Strategy

- **Unit testing** — individual functions/components within a module.
- **API testing** — each module's API surface against its contract.
- **Integration testing** — behavior across module boundaries and dependencies (see [Section 8](#8-cross-module-dependencies)).
- **Frontend testing** — UI components and module pages.
- **Cross-module testing** — validating that dependent modules correctly consume upstream contracts.
- **Search / semantic testing** — relevance and ranking behavior, where applicable, for [AI / Semantic Talent Search](#44-ai--semantic-talent-search).

> Specific testing frameworks/tools are not specified here and should be confirmed per [Section 36](#36-open-decisions--to-be-confirmed).

## 26. Definition of Done

A feature is considered done only when:

- [ ] Requirement implemented
- [ ] Module ownership respected
- [ ] API contract defined
- [ ] Code reviewed
- [ ] Tests passed
- [ ] Error handling in place
- [ ] Documentation updated
- [ ] No secrets committed
- [ ] Integration verified
- [ ] PR approved
- [ ] Merged into `develop`

## 27. Error Handling

APIs should return errors in a **consistent** shape across modules.

**Generic example only — not a confirmed, final TalentIQ error schema:**

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "The requested resource could not be found.",
    "details": {}
  }
}
```

## 28. Logging & Observability

- **Logging** — consistent, structured logging across modules.
- **Monitoring** — visibility into system health per module.
- **Traceability** — ability to trace a request across module boundaries.
- **Error diagnosis** — logs should support diagnosing failures without requiring guesswork.
- **Avoid sensitive information in logs** — do not log credentials, tokens, or other sensitive personal/production data.

## 29. Security Rules

**Never commit:**
- Passwords
- Tokens
- API keys
- Private credentials
- Production secrets

Use environment variables / configuration / an approved secret-management approach instead of hardcoding secrets. *(A specific secret-management platform has not been specified and should be confirmed — see [Section 36](#36-open-decisions--to-be-confirmed).)*

## 30. AI / Semantic Search Engineering Guidelines

This section expands on the engineering approach for [AI / Semantic Talent Search](#44-ai--semantic-talent-search).

**Logical flow:**

```
User Query → Query Interpretation → Embedding → Semantic Retrieval → Candidate Retrieval → Ranking → Results
```

- **Search ownership** — owned entirely by the AI / Semantic Talent Search module; other modules consume its results, they do not reimplement search logic.
- **Semantic retrieval** — retrieval of candidates based on meaning/similarity, not just exact keyword match.
- **Vector DB** — the intended storage/retrieval mechanism for embeddings. *(Product — To Be Confirmed.)*
- **Embeddings** — used to represent talent/skill data for semantic comparison. *(Model/provider — To Be Confirmed.)*
- **Ranking** — ordering candidate results by relevance.
- **Indexing** — keeping the search index up to date as underlying profile/skill/graph data changes.
- **Data dependencies** — this module depends on data from Talent Profile Intelligence, Skill Intelligence, and Skill Graph (see [Section 8](#8-cross-module-dependencies)).
- **Search vs. recommendation boundary** — this module returns ranked search results for a query; the [Talent Recommendation Engine](#47-talent-recommendation-engine) is responsible for proactive, fit-scored recommendations that build on search/similarity output.

> No specific machine learning model or provider is named here, per the project material — see [Section 36](#36-open-decisions--to-be-confirmed).

## 31. Database Ownership

- Each module **owns its own domain** in the database.
- Other modules should **consume data through approved interfaces** (APIs/contracts), not direct database queries.
- **Avoid direct SQL access** into another module's domain tables.
- Shared database infrastructure (see [Section 12](#12-shared-infrastructure)) does **not** imply shared ownership of business tables — infrastructure can be shared while data ownership remains per-module.

## 32. Messaging / Events

Where Kafka and/or RabbitMQ are used (see [Section 11](#11-technology-direction)), event-driven communication between modules should follow these principles:

- **Event ownership** — the module that owns the underlying data owns the events describing changes to it.
- **Producer** — the owning module publishes events for changes within its domain.
- **Consumer** — dependent modules (see [Section 8](#8-cross-module-dependencies)) subscribe to the events relevant to them.
- **Payload** — event payloads should be well-defined and versioned, not ad hoc.
- **Retry** — consumers should handle transient failures with retry logic appropriate to the event.
- **Failure handling** — failed event processing should be observable and recoverable, not silently dropped.
- **Versioning** — event schemas should be versioned to allow safe evolution.

> Specific event names/topics are not defined in the current project material and are intentionally not invented here.

## 33. API Versioning

APIs should be versioned, for example:

```
/api/v1/talent-profile-intelligence/...
/api/v1/skill-intelligence/...
/api/v1/skill-graph/...
/api/v1/ai-semantic-talent-search/...
/api/v1/skill-gap-analysis/...
/api/v1/technology-trend-analysis/...
/api/v1/talent-recommendation-engine/...
```

> These paths are **illustrative examples** of a versioning convention, not confirmed, finalized TalentIQ endpoints.

---

## 34. Documentation Rules

- **`README.md`** — project overview, module architecture, ownership, workflow, and operational rules (this document).
- **`docs/`** — deeper documentation: architecture decisions, API specifications, and module-specific detail that goes beyond what belongs in the top-level README.
- Architecture, API, and module documentation should be **updated whenever** a module's contract, ownership, or boundary changes — documentation drift should be treated as a defect, not a formatting nitpick.

## 35. Governance

The following are governed at the project level:

- Module ownership
- API contracts
- Cross-module dependencies
- Database ownership
- Infrastructure
- Security
- Branch strategy
- Release process
- Architectural decisions
- Technical lead / project lead coordination

Architectural or ownership changes should be coordinated with the relevant technical/project lead(s) rather than made unilaterally within a single module.

## 36. Open Decisions / To Be Confirmed

The following items are **not yet finalized** and should not be treated as settled:

- Final team allocation confirmation (see [Section 3](#3-current-seven-module-architecture))
- Total people/assignment count clarification (a possible discrepancy — e.g., 10 vs. 11 — has been flagged and is unresolved)
- Certification/assessment boundary confirmation (currently documented as part of Talent Profile Intelligence — see [Section 4.1](#41-talent-profile-intelligence))
- Source systems feeding the platform
- Vector DB choice
- Embedding model / provider
- Shared infrastructure ownership (operational responsibility)
- API standards beyond the general contract-first approach described in [Section 10](#10-frontend--backend-api-contract)
- Authentication / authorization provider or mechanism
- Secret-management platform
- Testing frameworks/tools
- Any other architecture decisions not explicitly confirmed in the source project material

> These items are intentionally left open rather than resolved in this document. They should be tracked and closed through the project's normal decision-making process (see [Section 35](#35-governance)).

## 37. Current Repository Structure

```
TalentIQ/
├── backend/
│   ├── api/
│   ├── cmd/
│   ├── config/
│   ├── database/
│   └── internal/
│       ├── ai-semantic-talent-search/
│       ├── skill-gap-analysis/
│       ├── skill-graph/
│       ├── skill-intelligence/
│       ├── talent-profile-intelligence/
│       ├── talent-recommendation-engine/
│       └── technology-trend-analysis/
├── frontend/
│   └── src/
│       ├── components/
│       ├── hooks/
│       ├── modules/
│       │   ├── ai-semantic-talent-search/
│       │   ├── skill-gap-analysis/
│       │   ├── skill-graph/
│       │   ├── skill-intelligence/
│       │   ├── talent-profile-intelligence/
│       │   ├── talent-recommendation-engine/
│       │   └── technology-trend-analysis/
│       ├── services/
│       ├── types/
│       └── utils/
├── docs/
├── infrastructure/
└── tests/
```

This structure reflects the current **seven-module** architecture only; older module groupings are not represented here.

## 38. Quick Git Reference

```bash
git clone https://github.com/MyHourly/TalentIQ.git   # Clone
cd TalentIQ                                           # Enter directory
git checkout <branch>                                 # Switch branch
git fetch origin                                      # Fetch updates
git pull origin <branch>                              # Pull latest
git branch -a                                         # List branches
git status                                             # Check working tree status
git diff                                               # Review changes
git add .                                              # Stage changes
git commit -m "feat(module): message"                  # Commit
git push                                               # Push
```

## 39. Final Engineering Workflow

```
Requirement → Module → Owner → Contract → Branch → Implementation → Test → Review → PR → develop → Integration → Approval → main
```

## 40. Project Principles

- **Modularity** — seven clearly separated functional modules.
- **Clear ownership** — one backend owner and one frontend owner per module.
- **Contract-first development** — contracts defined before implementation.
- **Loose coupling** — modules depend on contracts, not internals.
- **Reusability** — shared components, hooks, services, types, and utilities on the frontend; shared infrastructure on the backend.
- **Testability** — unit, API, integration, frontend, and cross-module testing.
- **Security** — no secrets in source control; approved secret management.
- **Observability** — consistent logging, monitoring, and traceability.
- **Maintainability** — documentation kept current with architecture and contracts.
- **Controlled integration** — changes flow through feature branches, PR review, `develop`, integration testing, and approval before reaching `main`.

---

*This README reflects the current seven-module structure of TalentIQ. Items marked "To Be Confirmed" are open and should be resolved through the project's normal governance process (see [Section 35](#35-governance)) rather than assumed.*
