# TalentIQ — Technology Trend Analysis (Module 4.6)

**Implementation Blueprint and Frontend/Backend Integration Contract**

*TalentIQ — Enterprise Talent Intelligence & Skill Graph Platform*

---

## Document Control

| Item | Detail |
|---|---|
| Module | 4.6 — Technology Trend Analysis |
| Document type | Implementation blueprint (single source of truth for Module 4.6) |
| Version | 1.0 (Draft for review) |
| Date | 1 October 2026 |
| Backend owner | Indra Kalyan Reddy ("Kalyan") |
| Frontend owner | Charan |
| Audience | Backend and frontend developers, reviewers, technical lead |
| Repository location | `docs/technology-trend-analysis/doc.md` |
| Backend code | `backend/internal/technology-trend-analysis/` |
| Frontend code | `frontend/src/modules/technology-trend-analysis/` |
| API base path | `/api/v1/technology-trend-analysis` |
| Branch flow | `feature/technology-trend-analysis` → `develop` → `main` |
| Contract authority | The frontend contract (`trend.types.ts`) is the source of truth; the backend conforms to it |

## Revision History

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026-10-01 | Module 4.6 team | Initial blueprint for review |

## Summary

Module 4.6 tracks technology and skill demand over time. It exposes seven read-only REST endpoints that return technology trends, per-technology history, emerging technologies, high-growth skills, capability gaps, an aggregated dashboard and a health check. The backend computes all derived values (growth, direction, classification, gap, severity) in one place, and the frontend renders them without recomputation.

**Conventions used in this document**

- **[DECISION]** marks a rule the source material does not define, with the proposal used here.
- **[ASSUMPTION]** marks an environment detail to verify against the codebase.
- Every **[DECISION]** and **[ASSUMPTION]** is consolidated in [Section 31](#31-open-items--decisions) for confirmation by Kalyan and Charan.

## Table of Contents

1. [Module Overview](#1-module-overview)
2. [Business Objective](#2-business-objective)
3. [Architecture](#3-architecture)
4. [Folder Structure](#4-folder-structure)
5. [Frontend Structure](#5-frontend-structure)
6. [Backend Structure](#6-backend-structure)
7. [Data Models](#7-data-models)
8. [DTOs](#8-dtos)
9. [API Contracts](#9-api-contracts)
10. [Request Examples](#10-request-examples)
11. [Response Examples](#11-response-examples)
12. [Error Responses](#12-error-responses)
13. [Validation Rules](#13-validation-rules)
14. [Query Parameters](#14-query-parameters)
15. [Business Logic](#15-business-logic)
16. [Growth Calculation Rules](#16-growth-calculation-rules)
17. [Capability Gap Calculation Rules](#17-capability-gap-calculation-rules)
18. [Classification Rules](#18-classification-rules)
19. [Frontend Integration Guide](#19-frontend-integration-guide)
20. [Component Mapping](#20-component-mapping)
21. [Backend Processing Flow](#21-backend-processing-flow)
22. [Frontend Processing Flow](#22-frontend-processing-flow)
23. [Future PostgreSQL Integration](#23-future-postgresql-integration)
24. [Future Redis Integration](#24-future-redis-integration)
25. [Future Kafka Integration](#25-future-kafka-integration)
26. [Testing Strategy](#26-testing-strategy)
27. [Mock Data Strategy](#27-mock-data-strategy)
28. [Deployment Considerations](#28-deployment-considerations)
29. [Ownership Matrix](#29-ownership-matrix)
30. [End-to-End Integration Flow](#30-end-to-end-integration-flow)
31. [Open Items & Decisions](#31-open-items--decisions)

---

## 1. Module Overview

**Technology Trend Analysis** tracks technology and skill demand signals over time. It answers four questions for the organization:

1. Which technologies are gaining or losing demand?
2. Which technologies are *emerging*?
3. Which skills are growing fastest (high-growth skills)?
4. Where does the organization's available capability fall short of what is required (capability gaps)?

### 1.1 What the module owns

Per the TalentIQ README (Section 4.6) and the module brief, this module owns:

- Technology demand analytics
- Historical trend analysis
- Growth analytics
- Emerging technology detection
- High-growth skill analytics
- Capability gap calculations (technology/skill capability view)
- Trend dashboards
- Trend insights

### 1.2 What the module does NOT own

| Concern | Owner |
|---|---|
| Skill definitions, taxonomy, skill IDs | 4.2 Skill Intelligence |
| Person/profile records | 4.1 Talent Profile Intelligence |
| Skill ↔ talent ↔ technology graph | 4.3 Skill Graph |
| Search and semantic ranking | 4.4 AI / Semantic Talent Search |
| Required-vs-inventory gap analysis for workforce planning | 4.5 Skill Gap Analysis |
| Recommendation scoring | 4.7 Talent Recommendation Engine |

> **Boundary note (4.5 vs 4.6).** Module 4.6 exposes a `capability-gaps` view derived from *trend/demand context* (required level vs. available level per skill). Module 4.5 owns the broader workforce skill-gap analysis. Both use the same gap math (Section 17) but own different data and endpoints. Do not read another module's tables directly; consume its API or events (README Section 9).

### 1.3 Upstream dependency

Skill Intelligence (4.2) → Technology Trend Analysis (4.6). Skill and technology identifiers (`skillId`, `technologyId`) originate from Skill Intelligence. This module stores them as references and does not redefine them.

---

## 2. Business Objective

From the TalentIQ project objective, this module delivers:

| Business benefit | How Module 4.6 delivers it |
|---|---|
| Identification of emerging and shortage skills | Emerging list (`growth > 40%`), capability gaps with severity |
| Better technology capability planning | Trend direction/classification per technology, history charts |
| Identification of training and upskilling needs | High/critical capability gaps surface where to train |
| Data-driven technology and workforce planning | Dashboard summarizing trend, growth and gap signals |
| Reduction of manual skill analysis | Automated growth, classification and gap calculations |

## 3. Architecture

### 3.1 Layered request flow (mandatory)

```
Request
  ↓
Router        (routes/routes.go)
  ↓
Handler       (handler/trend_handler.go)      — HTTP only: bind, validate, map errors
  ↓
Service       (service/trend_service.go)      — business orchestration, calls calculator
  ↓
Repository    (repository/trend_repository.go) — data access only
  ↓
Database      (in-memory mock now → PostgreSQL later)
  ↓
DTO           (dto/*.go)                       — map model → contract shape
  ↓
Response      (standard envelope)
```

### 3.2 Layer responsibilities

| Layer | Allowed | Forbidden |
|---|---|---|
| Router | Register routes, attach middleware | Business logic |
| Handler | Parse path/query, validate input, call service, write response | Calculations, DB access |
| Service | Orchestrate repository calls, apply `trend_calculator`, build DTOs | Gin types (`*gin.Context`), SQL |
| Repository | Fetch/persist models; filtering and sorting at data level | Business rules, HTTP concerns |
| Utils (`trend_calculator`) | Pure functions: growth, direction, classification, gap, severity | I/O, state |
| DTO | Contract shapes, JSON tags | Logic beyond mapping |

## 4. Folder Structure

### Repository layout (relevant portion)

```
TalentIQ/
├── backend/
│   ├── api/
│   ├── cmd/
│   ├── config/
│   ├── database/
│   └── internal/
│       └── technology-trend-analysis/        ← this module (backend)
├── frontend/
│   └── src/
│       ├── components/ hooks/ services/ types/ utils/   (shared)
│       └── modules/
│           └── technology-trend-analysis/    ← this module (frontend)
├── docs/
│   └── technology-trend-analysis/
│       └── doc.md                             ← this document
├── infrastructure/
└── tests/
```

## 5. Frontend Structure

```
frontend/src/modules/technology-trend-analysis/
├── components/
│   ├── CapabilityGapTable.tsx
│   ├── EmergingTechList.tsx
│   ├── HighGrowthSkillsList.tsx
│   ├── TrendCard.tsx
│   └── TrendChart.tsx
├── hooks/
│   └── useTrends.ts
├── pages/
│   ├── TrendDashboardPage.tsx
│   └── TrendDetailPage.tsx
├── services/
│   └── trendService.ts
└── types/
    └── trend.types.ts
```

| File | Responsibility |
|---|---|
| `types/trend.types.ts` | `TechnologyTrend`, `SkillTrend`, `CapabilityGap` interfaces (**contract source of truth**) |
| `services/trendService.ts` | Axios/fetch wrappers for the 5 contract calls; unwraps `data` |
| `hooks/useTrends.ts` | Loading/error/data state for each call |
| `pages/TrendDashboardPage.tsx` | Composes cards, lists, table |
| `pages/TrendDetailPage.tsx` | Technology detail + history chart (route param `technologyId`) |
| `components/*` | Presentational components (Section 20) |

---

## 6. Backend Structure

```
backend/internal/technology-trend-analysis/
├── constants/
│   └── constants.go
├── dto/
│   ├── dashboard_response.go
│   ├── technology_trend_response.go
│   ├── skill_trend_response.go
│   └── capability_gap_response.go
├── handler/
│   └── trend_handler.go
├── model/
│   ├── technology_trend.go
│   ├── skill_trend.go
│   ├── capability_gap.go
│   └── trend_history.go
├── repository/
│   └── trend_repository.go
├── routes/
│   └── routes.go
├── service/
│   └── trend_service.go
└── utils/
    └── trend_calculator.go
```

> **Go naming note.** The directory is `technology-trend-analysis` (hyphens, per the repo convention). A Go *package name* cannot contain hyphens, so use short package names: `constants`, `dto`, `handler`, `model`, `repository`, `routes`, `service`, `utils`. Import via full path, e.g. `github.com/MyHourly/TalentIQ/backend/internal/technology-trend-analysis/service`. **[ASSUMPTION]** Confirm the Go module path against `go.mod`.

---|---|
| `constants/constants.go` | Thresholds, enums, defaults, error codes |
| `model/*.go` | Domain/storage structs |
| `dto/*.go` | JSON contract structs + mappers |
| `repository/trend_repository.go` | `TrendRepository` interface + in-memory implementation |
| `service/trend_service.go` | Orchestration |
| `handler/trend_handler.go` | Gin handlers |
| `routes/routes.go` | Route registration |
| `utils/trend_calculator.go` | Pure calculation functions |

---

## 7. Data Models

Models are backend-internal. They may differ from DTOs, and DTOs must match the frontend contract exactly.

### 7.1 `constants/constants.go`

```go
package constants

// Direction values (frontend: 'up' | 'down' | 'stable')
const (
	DirectionUp     = "up"
	DirectionDown   = "down"
	DirectionStable = "stable"
)

// Classification values (frontend: 'emerging' | 'growing' | 'declining' | 'stable')
const (
	ClassEmerging  = "emerging"
	ClassGrowing   = "growing"
	ClassDeclining = "declining"
	ClassStable    = "stable"
)

// Severity values (frontend: 'low' | 'medium' | 'high' | 'critical')
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Direction thresholds (growth percent)
const (
	DirectionUpThreshold   = 10.0  // growth >  10 -> up
	DirectionDownThreshold = -10.0 // growth < -10 -> down
)

// Classification thresholds (growth percent)
const (
	EmergingThreshold  = 40.0  // growth >  40 -> emerging
	GrowingThreshold   = 15.0  // growth >  15 -> growing
	DecliningThreshold = -15.0 // growth < -15 -> declining
)

// Severity upper bounds on gapPercent (inclusive)
const (
	SeverityLowMax    = 20.0
	SeverityMediumMax = 40.0
	SeverityHighMax   = 60.0 // > 60 -> critical
)

// [DECISION] High-growth skill threshold (not defined in the brief).
const HighGrowthThreshold = 25.0 // growthPercent > 25 -> isHighGrowth

// [DECISION] Growth value when previousCount == 0 and demandCount > 0.
const NewDemandGrowthPercent = 100.0

// Pagination
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
	DefaultLimit    = 10
	MaxLimit        = 50
)

// Periods
const (
	Period30D  = "30d"
	Period90D  = "90d"
	Period180D = "180d"
	Period365D = "365d"
	DefaultPeriod = Period90D
)

// Error codes
const (
	ErrCodeInvalidParam     = "INVALID_PARAMETER"
	ErrCodeInvalidID        = "INVALID_TECHNOLOGY_ID"
	ErrCodeNotFound         = "RESOURCE_NOT_FOUND"
	ErrCodeInternal         = "INTERNAL_SERVER_ERROR"
	ErrCodeUnavailable      = "SERVICE_UNAVAILABLE"
)

const ModuleName = "technology-trend-analysis"
```

### 7.2 `model/technology_trend.go`

```go
package model

import "time"

// TechnologyTrend is a computed trend row for one technology and one period.
type TechnologyTrend struct {
	ID             string
	TechnologyID   string
	TechnologyName string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	DemandCount    int
	PreviousCount  int
	GrowthPercent  float64 // computed by utils.CalculateGrowth
	Direction      string  // computed by utils.DetermineDirection
	Classification string  // computed by utils.Classify
}
```

### 7.3 `model/skill_trend.go`

```go
package model

import "time"

type SkillTrend struct {
	ID            string
	SkillID       string
	SkillName     string
	PeriodStart   time.Time
	PeriodEnd     time.Time
	DemandCount   int
	PreviousCount int     // needed to compute growth; NOT exposed in the contract
	GrowthPercent float64
	IsHighGrowth  bool
}
```

> `previousCount` is stored on the model because growth is derived from it, but the frontend `SkillTrend` contract has no `previousCount` field, so the DTO omits it.

### 7.4 `model/capability_gap.go`

```go
package model

type CapabilityGap struct {
	ID             string
	SkillID        string
	SkillName      string
	RequiredLevel  float64
	AvailableLevel float64
	GapValue       float64
	GapPercent     float64
	Severity       string
}
```

### 7.5 `model/trend_history.go`

```go
package model

import "time"

// TrendHistoryPoint is one historical period for a technology.
// History is the stored raw series; TechnologyTrend rows are derived from it.
type TrendHistoryPoint struct {
	ID             string
	TechnologyID   string
	TechnologyName string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	DemandCount    int
}
```

### 7.6 Frontend contract (authoritative)

`frontend/src/modules/technology-trend-analysis/types/trend.types.ts`

```typescript
export interface TechnologyTrend {
  id: string;
  technologyId: string;
  technologyName: string;
  periodStart: string;
  periodEnd: string;
  demandCount: number;
  previousCount: number;
  growthPercent: number;
  direction: 'up' | 'down' | 'stable';
  classification: 'emerging' | 'growing' | 'declining' | 'stable';
}

export interface SkillTrend {
  id: string;
  skillId: string;
  skillName: string;
  periodStart: string;
  periodEnd: string;
  demandCount: number;
  growthPercent: number;
  isHighGrowth: boolean;
}

export interface CapabilityGap {
  id: string;
  skillId: string;
  skillName: string;
  requiredLevel: number;
  availableLevel: number;
  gapValue: number;
  gapPercent: number;
  severity: 'low' | 'medium' | 'high' | 'critical';
}
```

## 8. DTOs

DTOs mirror the frontend contract exactly. Mapper functions convert models to DTOs. Do not add extra fields to the three contract types (extra fields are tolerated by JS, but keep the contract pure).

### 8.1 `dto/technology_trend_response.go`

```go
package dto

import (
	"time"

	"github.com/MyHourly/TalentIQ/backend/internal/technology-trend-analysis/model"
)

type TechnologyTrendDTO struct {
	ID             string  `json:"id"`
	TechnologyID   string  `json:"technologyId"`
	TechnologyName string  `json:"technologyName"`
	PeriodStart    string  `json:"periodStart"`
	PeriodEnd      string  `json:"periodEnd"`
	DemandCount    int     `json:"demandCount"`
	PreviousCount  int     `json:"previousCount"`
	GrowthPercent  float64 `json:"growthPercent"`
	Direction      string  `json:"direction"`
	Classification string  `json:"classification"`
}

func ToTechnologyTrendDTO(m model.TechnologyTrend) TechnologyTrendDTO {
	return TechnologyTrendDTO{
		ID: m.ID, TechnologyID: m.TechnologyID, TechnologyName: m.TechnologyName,
		PeriodStart: m.PeriodStart.UTC().Format(time.RFC3339),
		PeriodEnd:   m.PeriodEnd.UTC().Format(time.RFC3339),
		DemandCount: m.DemandCount, PreviousCount: m.PreviousCount,
		GrowthPercent: m.GrowthPercent, Direction: m.Direction, Classification: m.Classification,
	}
}

func ToTechnologyTrendDTOs(ms []model.TechnologyTrend) []TechnologyTrendDTO {
	out := make([]TechnologyTrendDTO, 0, len(ms)) // never nil -> JSON [] not null
	for _, m := range ms {
		out = append(out, ToTechnologyTrendDTO(m))
	}
	return out
}
```

> **Critical:** always return an empty slice (`[]`), never `nil`. A `nil` slice serializes to `null` and breaks `.map()` in React.

### 8.2 `dto/skill_trend_response.go`

```go
package dto

type SkillTrendDTO struct {
	ID            string  `json:"id"`
	SkillID       string  `json:"skillId"`
	SkillName     string  `json:"skillName"`
	PeriodStart   string  `json:"periodStart"`
	PeriodEnd     string  `json:"periodEnd"`
	DemandCount   int     `json:"demandCount"`
	GrowthPercent float64 `json:"growthPercent"`
	IsHighGrowth  bool    `json:"isHighGrowth"`
}
// ToSkillTrendDTO / ToSkillTrendDTOs follow the same pattern as 8.1.
```

### 8.3 `dto/capability_gap_response.go`

```go
package dto

type CapabilityGapDTO struct {
	ID             string  `json:"id"`
	SkillID        string  `json:"skillId"`
	SkillName      string  `json:"skillName"`
	RequiredLevel  float64 `json:"requiredLevel"`
	AvailableLevel float64 `json:"availableLevel"`
	GapValue       float64 `json:"gapValue"`
	GapPercent     float64 `json:"gapPercent"`
	Severity       string  `json:"severity"`
}
// ToCapabilityGapDTO / ToCapabilityGapDTOs follow the same pattern as 8.1.
```

### 8.4 `dto/dashboard_response.go`

The dashboard endpoint is not in the frontend's five-call service list (see Section 31, item 1). Its shape is therefore proposed:

```go
package dto

type DashboardSummary struct {
	TotalTechnologies   int     `json:"totalTechnologies"`
	EmergingCount       int     `json:"emergingCount"`
	GrowingCount        int     `json:"growingCount"`
	DecliningCount      int     `json:"decliningCount"`
	StableCount         int     `json:"stableCount"`
	HighGrowthSkillCount int    `json:"highGrowthSkillCount"`
	CriticalGapCount    int     `json:"criticalGapCount"`
	AverageGrowthPercent float64 `json:"averageGrowthPercent"`
}

type DashboardResponse struct {
	Summary             DashboardSummary     `json:"summary"`
	TopTrends           []TechnologyTrendDTO `json:"topTrends"`
	EmergingTechnologies []TechnologyTrendDTO `json:"emergingTechnologies"`
	HighGrowthSkills    []SkillTrendDTO      `json:"highGrowthSkills"`
	CapabilityGaps      []CapabilityGapDTO   `json:"capabilityGaps"`
	Insights            []string             `json:"insights"`
	GeneratedAt         string               `json:"generatedAt"`
}
```

Each embedded list is already in contract shape, so a dashboard page can pass them straight into `EmergingTechList`, `HighGrowthSkillsList` and `CapabilityGapTable`.

### 8.5 Standard envelopes

```go
package dto

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type SuccessResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorBody struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Module    string `json:"module"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}
```

Health endpoint example: `GET /health` → `{ "data": { "status": "ok", "module": "technology-trend-analysis", "version": "v1", "timestamp": "2026-10-01T09:30:00Z" } }`

---

## 9. API Contracts

- **Base URL:** `/api/v1/technology-trend-analysis`
- **Method:** all endpoints are `GET` (read-only module).
- **Content type:** `application/json; charset=utf-8`
- **Auth:** To Be Confirmed project-wide (README Section 36). The module is written so auth middleware can be attached in `routes.go` with no handler changes.

| # | Endpoint | Frontend call? | `data` type | Purpose |
|---|---|---|---|---|
| 1 | `GET /dashboard` | Optional (not in service list) | `DashboardResponse` (object) | Aggregated dashboard |
| 2 | `GET /trends` | **Yes** | `TechnologyTrend[]` | List technology trends |
| 3 | `GET /trends/:technologyId` | **Yes** | `TechnologyTrend[]` (history) | One technology's history |
| 4 | `GET /emerging` | **Yes** | `TechnologyTrend[]` | Emerging technologies |
| 5 | `GET /high-growth` | **Yes** | `SkillTrend[]` | High-growth skills |
| 6 | `GET /capability-gaps` | **Yes** | `CapabilityGap[]` | Capability gaps |
| 7 | `GET /health` | No | `HealthResponse` (object) | Liveness/readiness |

> **Frontend path note.** The frontend lists paths like `GET /technology-trend-analysis/trends`. The `/api/v1` prefix is supplied by the shared API client's `baseURL`. Final URL = `baseURL (/api/v1)` + `/technology-trend-analysis/trends`. Confirm that `baseURL` in `frontend/src/services/` includes `/api/v1`.

### 9.1 Response envelope

Success (collection):
```json
{ "data": [ ... ], "pagination": { "page": 1, "pageSize": 20, "total": 7 } }
```
Success (object):
```json
{ "data": { ... } }
```
Error:
```json
{ "error": { "code": "RESOURCE_NOT_FOUND", "message": "...", "details": {} } }
```

The frontend reads `response.data.data`. The `pagination` key is optional and ignored by the current frontend. It is included for README conformity and future use. It is only emitted on `/trends`.

### 9.2 Endpoint 3 — detail semantics **[DECISION]**

The brief gives no frontend type for the detail response and states that every call returns `{ "data": [...] }`. So `GET /trends/:technologyId` returns **an array of `TechnologyTrend` ordered by `periodStart` ascending**, one entry per historical period of that technology. This:

- reuses the existing `TechnologyTrend` type (no new type is needed),
- feeds `TrendChart.tsx` directly (x = `periodStart`, y = `demandCount`),
- gives `TrendDetailPage.tsx` its "current" record as the **last** element.

If Charan prefers a single-object detail response, change `data` to the latest entry and add `history: TechnologyTrend[]`. That is a contract change and must be agreed first.

## 10. Request Examples

All examples assume `http://localhost:8080` as the host.

```bash
# Dashboard
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/dashboard"

# All trends, default period (90d), default pagination
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/trends"

# Declining technologies only, sorted by growth ascending
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/trends?classification=declining&sortBy=growthPercent&order=asc"

# Search by name, 180-day window
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/trends?search=go&period=180d"

# Detail / history for one technology
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/trends/tech-go"

# Emerging technologies, top 5
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/emerging?limit=5"

# High-growth skills
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/high-growth?limit=10"

# Critical + high capability gaps
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/capability-gaps?severity=critical"

# Health
curl -s "http://localhost:8080/api/v1/technology-trend-analysis/health"
```

Frontend (`trendService.ts`) equivalents are in Section 19.

---

## 11. Response Examples

Reference data used across all examples (current period = Q3 2026: `2026-07-01` → `2026-09-30`; previous period = Q2 2026). All numbers below were computed with the rules in Sections 16–18 and are valid test fixtures.

### 11.1 `GET /trends` → 200

```json
{
  "data": [
    {
      "id": "trend-go-2026q3", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 1240,
      "previousCount": 780, "growthPercent": 58.97,
      "direction": "up", "classification": "emerging"
    },
    {
      "id": "trend-vector-db-2026q3", "technologyId": "tech-vector-db",
      "technologyName": "Vector Database", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 380,
      "previousCount": 190, "growthPercent": 100.0,
      "direction": "up", "classification": "emerging"
    },
    {
      "id": "trend-kubernetes-2026q3", "technologyId": "tech-kubernetes",
      "technologyName": "Kubernetes", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 2100,
      "previousCount": 1750, "growthPercent": 20.0,
      "direction": "up", "classification": "growing"
    },
    {
      "id": "trend-java-2026q3", "technologyId": "tech-java",
      "technologyName": "Java", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 900,
      "previousCount": 1000, "growthPercent": -10.0,
      "direction": "stable", "classification": "stable"
    },
    {
      "id": "trend-jquery-2026q3", "technologyId": "tech-jquery",
      "technologyName": "jQuery", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 300,
      "previousCount": 450, "growthPercent": -33.33,
      "direction": "down", "classification": "declining"
    }
  ],
  "pagination": {"page": 1, "pageSize": 20, "total": 5}
}
```

(Abbreviated to 5 of the 7 fixture technologies; The full fixture also contains the full fixture also contains Rust `410/260 → 57.69, up, emerging` and SQL `1500/1450 → 3.45, stable, stable`.)

> Edge example: Java is exactly `-10.00`. Direction uses a strict `< -10`, so it is `stable`, not `down`.

### 11.2 `GET /trends/tech-go` → 200 (history, ascending)

```json
{
  "data": [
    {
      "id": "trend-go-2025q4", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2025-10-01T00:00:00Z",
      "periodEnd": "2025-12-31T23:59:59Z", "demandCount": 640,
      "previousCount": 600, "growthPercent": 6.67,
      "direction": "stable", "classification": "stable"
    },
    {
      "id": "trend-go-2026q1", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2026-01-01T00:00:00Z",
      "periodEnd": "2026-03-31T23:59:59Z", "demandCount": 700,
      "previousCount": 640, "growthPercent": 9.38,
      "direction": "stable", "classification": "stable"
    },
    {
      "id": "trend-go-2026q2", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2026-04-01T00:00:00Z",
      "periodEnd": "2026-06-30T23:59:59Z", "demandCount": 780,
      "previousCount": 700, "growthPercent": 11.43,
      "direction": "up", "classification": "stable"
    },
    {
      "id": "trend-go-2026q3", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 1240,
      "previousCount": 780, "growthPercent": 58.97,
      "direction": "up", "classification": "emerging"
    }
  ]
}
```

> Q2 shows `direction: up` with `classification: stable` (growth 11.43 is above 10 but not above 15). This is correct per the rules and is not a bug. The UI must display both independently.

### 11.3 `GET /emerging` → 200

Technologies whose current-period `classification == "emerging"`, sorted by `growthPercent` desc.

```json
{
  "data": [
    {
      "id": "trend-vector-db-2026q3", "technologyId": "tech-vector-db",
      "technologyName": "Vector Database", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 380,
      "previousCount": 190, "growthPercent": 100.0,
      "direction": "up", "classification": "emerging"
    },
    {
      "id": "trend-go-2026q3", "technologyId": "tech-go",
      "technologyName": "Go", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 1240,
      "previousCount": 780, "growthPercent": 58.97,
      "direction": "up", "classification": "emerging"
    },
    {
      "id": "trend-rust-2026q3", "technologyId": "tech-rust",
      "technologyName": "Rust", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 410,
      "previousCount": 260, "growthPercent": 57.69,
      "direction": "up", "classification": "emerging"
    }
  ]
}
```

### 11.4 `GET /high-growth` → 200

Skills with `isHighGrowth == true`, sorted by `growthPercent` desc.

```json
{
  "data": [
    {
      "id": "skilltrend-vector-db-2026q3", "skillId": "skill-vector-db",
      "skillName": "Vector Database", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 380,
      "growthPercent": 100.0, "isHighGrowth": true
    },
    {
      "id": "skilltrend-cloud-security-2026q3", "skillId": "skill-cloud-security",
      "skillName": "Cloud Security", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 900,
      "growthPercent": 50.0, "isHighGrowth": true
    },
    {
      "id": "skilltrend-kafka-2026q3", "skillId": "skill-kafka",
      "skillName": "Kafka", "periodStart": "2026-07-01T00:00:00Z",
      "periodEnd": "2026-09-30T23:59:59Z", "demandCount": 1320,
      "growthPercent": 32.0, "isHighGrowth": true
    }
  ]
}
```

### 11.5 `GET /capability-gaps` → 200

Sorted by `gapPercent` desc.

```json
{
  "data": [
    {
      "id": "gap-cloud-security", "skillId": "skill-cloud-security",
      "skillName": "Cloud Security", "requiredLevel": 5,
      "availableLevel": 1.5, "gapValue": 3.5,
      "gapPercent": 70.0, "severity": "critical"
    },
    {
      "id": "gap-vector-db", "skillId": "skill-vector-db",
      "skillName": "Vector Database", "requiredLevel": 4,
      "availableLevel": 1.6, "gapValue": 2.4,
      "gapPercent": 60.0, "severity": "high"
    },
    {
      "id": "gap-kubernetes", "skillId": "skill-kubernetes",
      "skillName": "Kubernetes", "requiredLevel": 4,
      "availableLevel": 3, "gapValue": 1,
      "gapPercent": 25.0, "severity": "medium"
    },
    {
      "id": "gap-kafka", "skillId": "skill-kafka",
      "skillName": "Kafka", "requiredLevel": 4,
      "availableLevel": 3.2, "gapValue": 0.8,
      "gapPercent": 20.0, "severity": "low"
    }
  ]
}
```

> Boundary examples: Vector Database is exactly `60.00` → `high` (the critical band starts above 60). Kafka is exactly `20.00` → `low`.

### 11.6 `GET /dashboard` → 200

```jsonc
{
  "data": {
    "summary": {
      "totalTechnologies": 7,
      "emergingCount": 3,
      "growingCount": 1,
      "decliningCount": 1,
      "stableCount": 2,
      "highGrowthSkillCount": 3,
      "criticalGapCount": 1,
      "averageGrowthPercent": 28.11
    },
    "topTrends": [ /* 5 TechnologyTrend, top by demandCount desc: Kubernetes, SQL, Go, Java, Rust */ ],
    "emergingTechnologies": [ /* same items as 11.3 */ ],
    "highGrowthSkills": [ /* same items as 11.4 */ ],
    "capabilityGaps": [ /* same items as 11.5 */ ],
    "insights": [
      "Vector Database demand doubled (+100.0%) versus the previous period.",
      "3 technologies are classified as emerging: Vector Database, Go, Rust.",
      "Cloud Security has a critical capability gap (70.0% below the required level).",
      "jQuery demand fell 33.33%, so it is classified as declining."
    ],
    "generatedAt": "2026-10-01T09:30:00Z"
  }
}
```

## 12. Error Responses

All errors use the README standard shape. `details` is always an object (use `{}` when empty).

```json
{ "error": { "code": "<CODE>", "message": "<human readable>", "details": {} } }
```

| HTTP | Code | When | Example `details` |
|---|---|---|---|
| 400 | `INVALID_PARAMETER` | Bad query parameter value | `{"field":"period","allowed":["30d","90d","180d","365d"]}` |
| 400 | `INVALID_TECHNOLOGY_ID` | `:technologyId` empty/malformed | `{"field":"technologyId"}` |
| 404 | `RESOURCE_NOT_FOUND` | Unknown `technologyId` on detail | `{"technologyId":"tech-xyz"}` |
| 405 | `METHOD_NOT_ALLOWED` | Non-GET method | `{}` |
| 500 | `INTERNAL_SERVER_ERROR` | Unexpected failure | `{}` (no internals leaked) |
| 503 | `SERVICE_UNAVAILABLE` | Data source down (future DB) | `{}` |

### 12.1 Examples

`GET /trends?period=7d` → 400
```json
{ "error": { "code": "INVALID_PARAMETER",
  "message": "Query parameter 'period' must be one of: 30d, 90d, 180d, 365d.",
  "details": { "field": "period", "allowed": ["30d","90d","180d","365d"] } } }
```

`GET /trends/tech-unknown` → 404
```json
{ "error": { "code": "RESOURCE_NOT_FOUND",
  "message": "No trend data found for technology 'tech-unknown'.",
  "details": { "technologyId": "tech-unknown" } } }
```

### 12.2 Frontend error handling contract

- Read `error.response.data.error.message` for display and `error.response.data.error.code` for logic.
- 404 on the detail page renders a "technology not found" state.
- 5xx renders a retry state.
- An empty list (`data: []`) is not an error, so render an empty state.

---

## 13. Validation Rules

Validation happens in the **handler** (HTTP-level) before the service is called. Business-rule validation (e.g., negative levels in stored data) is handled in the service/repository and treated as data errors.

| Parameter | Rule | On failure |
|---|---|---|
| `technologyId` (path) | Required, 1–64 chars, regex `^[A-Za-z0-9][A-Za-z0-9_-]*$` | 400 `INVALID_TECHNOLOGY_ID` |
| `period` | One of `30d`, `90d`, `180d`, `365d`; default `90d` | 400 |
| `page` | Integer ≥ 1; default 1 | 400 |
| `pageSize` | Integer 1–100; default 20 | 400 |
| `limit` | Integer 1–50; default 10 (emerging, high-growth, capability-gaps) | 400 |
| `classification` | One of `emerging`, `growing`, `declining`, `stable` | 400 |
| `direction` | One of `up`, `down`, `stable` | 400 |
| `severity` | One of `low`, `medium`, `high`, `critical` | 400 |
| `sortBy` | Allow-list per endpoint (Section 14) | 400 |
| `order` | `asc` or `desc`; default `desc` | 400 |
| `search` | 0–64 chars; trimmed; case-insensitive contains match | 400 if > 64 |
| unknown query params | Ignored (do not fail) | — |

### Data integrity invariants (service/repository level)

- `demandCount >= 0`, `previousCount >= 0`.
- `periodEnd > periodStart`.
- `requiredLevel >= 0`, `availableLevel >= 0`.
- A row that violates an invariant is **skipped and logged**; it is never returned with fabricated values.

## 14. Query Parameters

| Endpoint | Parameters |
|---|---|
| `/dashboard` | `period` |
| `/trends` | `period`, `page`, `pageSize`, `search`, `classification`, `direction`, `sortBy`, `order` |
| `/trends/:technologyId` | `period` (history span: how far back to return), `limit` (max periods; default 12, max 24) |
| `/emerging` | `period`, `limit`, `sortBy`, `order` |
| `/high-growth` | `period`, `limit`, `sortBy`, `order` |
| `/capability-gaps` | `severity`, `limit`, `sortBy`, `order` |
| `/health` | none |

### 14.1 `sortBy` allow-lists

| Endpoint | Allowed `sortBy` | Default |
|---|---|---|
| `/trends`, `/emerging` | `growthPercent`, `demandCount`, `technologyName` | `growthPercent desc` |
| `/high-growth` | `growthPercent`, `demandCount`, `skillName` | `growthPercent desc` |
| `/capability-gaps` | `gapPercent`, `gapValue`, `skillName` | `gapPercent desc` |

Ties are broken by the name field ascending, so ordering is deterministic and testable.

### 14.2 `period` semantics **[DECISION]**

`period` sets the window length. The *current* window ends "now" (or at the latest data date), and the *previous* window is the immediately preceding window of equal length.

| `period` | Current window | Previous window |
|---|---|---|
| `30d` | last 30 days | the 30 days before that |
| `90d` (default) | last 90 days | the 90 days before that |
| `180d` | last 180 days | the 180 days before that |
| `365d` | last 365 days | the 365 days before that |

> The fixtures in Section 11 use calendar quarters for readability. With mock data, the repository returns the pre-aggregated fixtures regardless of `period`. Real period handling arrives with PostgreSQL (Section 23).

---

## 15. Business Logic

### 15.1 Service interface

```go
type TrendService interface {
	GetDashboard(ctx context.Context, f DashboardFilter) (*dto.DashboardResponse, error)
	ListTrends(ctx context.Context, f TrendFilter) ([]dto.TechnologyTrendDTO, int, error) // items, total
	GetTrendHistory(ctx context.Context, technologyID string, limit int) ([]dto.TechnologyTrendDTO, error)
	ListEmerging(ctx context.Context, f ListFilter) ([]dto.TechnologyTrendDTO, error)
	ListHighGrowthSkills(ctx context.Context, f ListFilter) ([]dto.SkillTrendDTO, error)
	ListCapabilityGaps(ctx context.Context, f GapFilter) ([]dto.CapabilityGapDTO, error)
	Health(ctx context.Context) dto.HealthResponse
}
```

### 15.2 Processing pipeline per endpoint

| Endpoint | Steps |
|---|---|
| `/trends` | repo `ListTechnologyTrends` → for each row compute growth/direction/classification → filter → sort → paginate → DTO |
| `/trends/:id` | repo `GetHistory(technologyID)` → `ErrNotFound` if empty → compute each period's derived fields from its own `demandCount` and the previous point's `demandCount` → ascending → DTO |
| `/emerging` | same as `/trends` then keep `classification == emerging` → sort → limit |
| `/high-growth` | repo `ListSkillTrends` → compute growth → `isHighGrowth = growth > 25` → keep true → sort → limit |
| `/capability-gaps` | repo `ListSkillLevels` → compute gap/percent/severity → filter by severity → sort → limit |
| `/dashboard` | run the above (period-scoped) and aggregate counts, average growth, top trends and insights |

### 15.3 Derived values are always computed, never trusted from storage

Even if a future DB stores `growth_percent`, the service **recomputes** from `demandCount` and `previousCount`. This guarantees the rules in this document are the only source of classification.

## 16. Growth Calculation Rules

### 16.1 Formula

```
growthPercent = ((demandCount - previousCount) / previousCount) × 100
```

Rounded to **2 decimal places** (half away from zero) after calculation. All comparisons in Sections 16–18 use the **rounded** value, so what the user sees always matches the classification they see.

### 16.2 Edge cases **[DECISION]**

The brief's formula divides by `previousCount`. These cases are needed to avoid division by zero.

| Case | `growthPercent` | Rationale |
|---|---|---|
| `previousCount > 0` | formula | normal |
| `previousCount == 0` and `demandCount > 0` | `100.0` (`NewDemandGrowthPercent`) | New demand appearing from nothing is flagged as strong growth (→ `up`, `emerging`), without returning `Infinity` (invalid JSON) |
| `previousCount == 0` and `demandCount == 0` | `0.0` | No signal → `stable` |
| `demandCount < previousCount` | negative | decline |
| negative inputs | rejected as invalid data | Section 13.1 |

### 16.3 Direction

| Condition | Direction |
|---|---|
| `growthPercent > 10` | `up` |
| `growthPercent < -10` | `down` |
| otherwise (`-10 ≤ g ≤ 10`) | `stable` |

### 16.4 Reference implementation — `utils/trend_calculator.go`

```go
package utils

import (
	"math"

	"github.com/MyHourly/TalentIQ/backend/internal/technology-trend-analysis/constants"
)

// Round2 rounds half away from zero to 2 decimals.
func Round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// CalculateGrowth returns growth percent rounded to 2 dp.
func CalculateGrowth(current, previous int) float64 {
	switch {
	case previous > 0:
		return Round2((float64(current-previous) / float64(previous)) * 100)
	case current > 0:
		return constants.NewDemandGrowthPercent
	default:
		return 0
	}
}

func DetermineDirection(growth float64) string {
	switch {
	case growth > constants.DirectionUpThreshold:
		return constants.DirectionUp
	case growth < constants.DirectionDownThreshold:
		return constants.DirectionDown
	default:
		return constants.DirectionStable
	}
}

func Classify(growth float64) string {
	switch {
	case growth > constants.EmergingThreshold:
		return constants.ClassEmerging
	case growth > constants.GrowingThreshold:
		return constants.ClassGrowing
	case growth < constants.DecliningThreshold:
		return constants.ClassDeclining
	default:
		return constants.ClassStable
	}
}

func IsHighGrowth(growth float64) bool {
	return growth > constants.HighGrowthThreshold
}

// CalculateGap returns gapValue, gapPercent, severity.
func CalculateGap(required, available float64) (gapValue, gapPercent float64, severity string) {
	gapValue = required - available
	if gapValue < 0 {
		gapValue = 0 // surplus is not a gap
	}
	gapValue = Round2(gapValue)
	if required > 0 {
		gapPercent = Round2((gapValue / required) * 100)
	}
	return gapValue, gapPercent, DetermineSeverity(gapPercent)
}

func DetermineSeverity(gapPercent float64) string {
	switch {
	case gapPercent <= constants.SeverityLowMax:
		return constants.SeverityLow
	case gapPercent <= constants.SeverityMediumMax:
		return constants.SeverityMedium
	case gapPercent <= constants.SeverityHighMax:
		return constants.SeverityHigh
	default:
		return constants.SeverityCritical
	}
}
```

## 17. Capability Gap Calculation Rules

### 17.1 Formulas

```
gapValue   = requiredLevel - availableLevel          (clamped at 0 if negative)
gapPercent = (gapValue / requiredLevel) × 100        (0 if requiredLevel == 0)
```

Both values are rounded to 2 decimals.

### 17.2 Severity (from `gapPercent`)

The brief states bands as `0-20`, `21-40`, `41-60`, `61+`. Because `gapPercent` is a decimal, the bands are defined without gaps:

| `gapPercent` | Severity |
|---|---|
| `≤ 20` | `low` |
| `> 20` and `≤ 40` | `medium` |
| `> 40` and `≤ 60` | `high` |
| `> 60` | `critical` |

### 17.3 Edge cases **[DECISION]**

| Case | Result |
|---|---|
| `availableLevel > requiredLevel` | `gapValue = 0`, `gapPercent = 0`, `low` (a surplus is not a gap) |
| `requiredLevel == 0` | `gapValue = 0`, `gapPercent = 0`, `low` (avoids divide-by-zero) |
| negative levels | invalid data, skipped and logged |

### 17.4 Level scale **[ASSUMPTION]**

The brief does not specify the scale of `requiredLevel`/`availableLevel`. This document uses a **0–5 proficiency scale** in examples. Because severity depends on a *percentage*, the logic is scale-independent, so a different scale (e.g. 0–10) requires no code change. Confirm the scale with Skill Intelligence (4.2), which owns proficiency definitions. The `numeric` columns in Section 23 do not constrain the scale.

## 18. Classification Rules

### 18.1 Technology classification (evaluated top-down, first match wins)

| Order | Condition | Classification |
|---|---|---|
| 1 | `growth > 40` | `emerging` |
| 2 | `growth > 15` | `growing` |
| 3 | `growth < -15` | `declining` |
| 4 | otherwise | `stable` |

### 18.2 Direction vs. classification matrix

The two fields are independent. Both are always returned.

| Growth | Direction | Classification |
|---|---|---|
| `> 40` | up | emerging |
| `> 15 – 40` | up | growing |
| `> 10 – 15` | up | **stable** |
| `-10 – 10` | stable | stable |
| `-15 – < -10` | down | **stable** |
| `< -15` | down | declining |

The bold rows are intentional (e.g. growth 12 → up/stable; growth −12 → down/stable). They must be covered by unit tests (Section 26) and displayed correctly by the UI.

### 18.3 High-growth skill **[DECISION]**

`isHighGrowth = growthPercent > 25` (`HighGrowthThreshold`). The brief names "High-growth skills" but gives no threshold. 25 sits between "growing" (>15) and "emerging" (>40). It is a single constant. Change it in `constants.go` only, and never in the frontend.

### 18.4 Boundary summary (test oracle)

| Growth | Direction | Classification |
|---|---|---|
| 40.00 | up | growing |
| 40.01 | up | emerging |
| 15.00 | up | stable |
| 15.01 | up | growing |
| 10.00 | stable | stable |
| 10.01 | up | stable |
| −10.00 | stable | stable |
| −10.01 | down | stable |
| −15.00 | down | stable |
| −15.01 | down | declining |

---

## 19. Frontend Integration Guide

### 19.1 Rules for Charan

1. Use the interfaces in `trend.types.ts` **unchanged**. If a field needs to change, update this document and notify Kalyan first.
2. Do **not** recompute `growthPercent`, `direction`, `classification`, `gapValue`, `gapPercent`, `severity` or `isHighGrowth`. Render them.
3. Always unwrap one level: `response.data.data`.
4. Treat `periodStart`/`periodEnd` as ISO-8601 UTC strings and parse them only at display time.
5. Never assume a list is non-empty. Handle `[]`.
6. Use the shared API client in `frontend/src/services/` (no hardcoded host).

### 19.2 `services/trendService.ts`

```typescript
import api from '../../../services/apiClient'; // shared axios instance; baseURL ends with /api/v1
import { TechnologyTrend, SkillTrend, CapabilityGap } from '../types/trend.types';

const BASE = '/technology-trend-analysis';

interface ApiResponse<T> { data: T; pagination?: { page: number; pageSize: number; total: number } }

export interface TrendQuery {
  period?: '30d' | '90d' | '180d' | '365d';
  page?: number;
  pageSize?: number;
  search?: string;
  classification?: TechnologyTrend['classification'];
  direction?: TechnologyTrend['direction'];
  sortBy?: 'growthPercent' | 'demandCount' | 'technologyName';
  order?: 'asc' | 'desc';
}

export const trendService = {
  getTrends: async (q?: TrendQuery): Promise<TechnologyTrend[]> =>
    (await api.get<ApiResponse<TechnologyTrend[]>>(`${BASE}/trends`, { params: q })).data.data,

  getTrendByTechnology: async (technologyId: string): Promise<TechnologyTrend[]> =>
    (await api.get<ApiResponse<TechnologyTrend[]>>(`${BASE}/trends/${encodeURIComponent(technologyId)}`)).data.data,

  getEmerging: async (limit?: number): Promise<TechnologyTrend[]> =>
    (await api.get<ApiResponse<TechnologyTrend[]>>(`${BASE}/emerging`, { params: { limit } })).data.data,

  getHighGrowth: async (limit?: number): Promise<SkillTrend[]> =>
    (await api.get<ApiResponse<SkillTrend[]>>(`${BASE}/high-growth`, { params: { limit } })).data.data,

  getCapabilityGaps: async (severity?: CapabilityGap['severity']): Promise<CapabilityGap[]> =>
    (await api.get<ApiResponse<CapabilityGap[]>>(`${BASE}/capability-gaps`, { params: { severity } })).data.data,
};
```

> `import api from ...` assumes the shared client's file name. **[ASSUMPTION]** Adjust to the real shared client.

### 19.3 Dashboard endpoint usage

The five contract calls are enough to build the dashboard page. `GET /dashboard` is an optimization that returns the same lists plus summary counts in one round trip. The frontend may adopt it later without changing component props, because its embedded arrays use the same types.

## 20. Component Mapping

| Component | Data source (type) | API | Key fields rendered |
|---|---|---|---|
| `TrendCard.tsx` | one `TechnologyTrend` | `/trends` (item) | `technologyName`, `demandCount`, `growthPercent`, `direction`, `classification` |
| `TrendChart.tsx` | `TechnologyTrend[]` (history) | `/trends/:technologyId` | x: `periodStart`; y: `demandCount` (optionally `previousCount`) |
| `EmergingTechList.tsx` | `TechnologyTrend[]` | `/emerging` | `technologyName`, `growthPercent`, `demandCount` |
| `HighGrowthSkillsList.tsx` | `SkillTrend[]` | `/high-growth` | `skillName`, `growthPercent`, `demandCount`, `isHighGrowth` |
| `CapabilityGapTable.tsx` | `CapabilityGap[]` | `/capability-gaps` | `skillName`, `requiredLevel`, `availableLevel`, `gapValue`, `gapPercent`, `severity` |
| `TrendDashboardPage.tsx` | all of the above | `/trends`, `/emerging`, `/high-growth`, `/capability-gaps` (or `/dashboard`) | composes cards, lists, table |
| `TrendDetailPage.tsx` | `TechnologyTrend[]` | `/trends/:technologyId` | header from last item; `TrendChart` from full array |

## 21. Backend Processing Flow

### 21.1 Example: `GET /api/v1/technology-trend-analysis/trends?classification=declining`

```
1. Gin router matches GET /api/v1/technology-trend-analysis/trends
2. Middleware: request ID → logging → recovery (→ auth, future)
3. Handler.ListTrends
     a. parse + validate query (period, page, pageSize, classification, ...)
     b. invalid → 400 ErrorResponse (stop)
4. Service.ListTrends(ctx, filter)
     a. repo.ListTechnologyTrends(ctx, period)            → []model.TechnologyTrend (raw counts)
     b. for each row: growth = CalculateGrowth(demand, previous)
                       direction = DetermineDirection(growth)
                       classification = Classify(growth)
     c. filter by classification/direction/search
     d. sort (sortBy, order; tie-break by name)
     e. total = len(filtered); paginate
5. dto.ToTechnologyTrendDTOs(...)   (non-nil slice)
6. Handler writes 200 {"data":[...],"pagination":{...}}
```

### 21.2 Repository interface

```go
type TrendRepository interface {
	ListTechnologyTrends(ctx context.Context, period string) ([]model.TechnologyTrend, error)
	GetTrendHistory(ctx context.Context, technologyID string, limit int) ([]model.TrendHistoryPoint, error)
	ListSkillTrends(ctx context.Context, period string) ([]model.SkillTrend, error)
	ListCapabilityGaps(ctx context.Context) ([]model.CapabilityGap, error) // required/available levels
	Ping(ctx context.Context) error
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnavailable  = errors.New("unavailable")
)
```

### 21.3 Routes — `routes/routes.go`

```go
func RegisterRoutes(rg *gin.RouterGroup, h *handler.TrendHandler) {
	g := rg.Group("/technology-trend-analysis") // rg is already /api/v1
	{
		g.GET("/health", h.Health)
		g.GET("/dashboard", h.Dashboard)
		g.GET("/trends", h.ListTrends)
		g.GET("/trends/:technologyId", h.GetTrendByTechnology)
		g.GET("/emerging", h.ListEmerging)
		g.GET("/high-growth", h.ListHighGrowth)
		g.GET("/capability-gaps", h.ListCapabilityGaps)
	}
}
```

Wiring (in `cmd/`): `repo := repository.NewInMemoryTrendRepository()` → `svc := service.NewTrendService(repo)` → `h := handler.NewTrendHandler(svc)` → `routes.RegisterRoutes(v1, h)`.

## 22. Frontend Processing Flow

### 22.1 Dashboard load

```
TrendDashboardPage mounts
  ├─ useTrends()            → GET /trends            → TrendCard[]
  ├─ useEmerging(5)         → GET /emerging?limit=5  → EmergingTechList
  ├─ useHighGrowth(5)       → GET /high-growth       → HighGrowthSkillsList
  └─ useCapabilityGaps()    → GET /capability-gaps   → CapabilityGapTable
Each hook independently manages loading / error / data.
Sections render skeletons while loading, an error panel with retry on failure,
and an empty state on [].
```

### 22.2 Detail load

```
User clicks TrendCard → navigate /technology-trends/:technologyId
TrendDetailPage → useTrendDetail(technologyId) → GET /trends/:technologyId
  ├─ 200 → header uses last element; TrendChart plots all elements
  ├─ 404 → "Technology not found" state
  └─ 5xx → retry state
```

## 23. Future PostgreSQL Integration

Phase 1 uses an in-memory repository. Phase 2 swaps in a PostgreSQL repository behind the same `TrendRepository` interface, with **no handler, service, DTO or frontend change**. Table names are prefixed `tta_` to mark module ownership (README Section 31: no cross-module table access).

### 23.1 Proposed schema

```sql
-- Raw demand signal per technology per period (source for TechnologyTrend)
CREATE TABLE tta_technology_demand (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    technology_id    VARCHAR(64)  NOT NULL,     -- reference to Skill Intelligence (4.2); no FK across modules
    technology_name  VARCHAR(128) NOT NULL,     -- denormalized snapshot for display
    period_start     TIMESTAMPTZ  NOT NULL,
    period_end       TIMESTAMPTZ  NOT NULL,
    demand_count     INTEGER      NOT NULL CHECK (demand_count >= 0),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CHECK (period_end > period_start),
    UNIQUE (technology_id, period_start, period_end)
);
CREATE INDEX idx_tta_tech_demand_tech_period ON tta_technology_demand (technology_id, period_start DESC);

CREATE TABLE tta_skill_demand (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id      VARCHAR(64)  NOT NULL,
    skill_name    VARCHAR(128) NOT NULL,
    period_start  TIMESTAMPTZ  NOT NULL,
    period_end    TIMESTAMPTZ  NOT NULL,
    demand_count  INTEGER      NOT NULL CHECK (demand_count >= 0),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CHECK (period_end > period_start),
    UNIQUE (skill_id, period_start, period_end)
);
CREATE INDEX idx_tta_skill_demand_period ON tta_skill_demand (period_start DESC);

CREATE TABLE tta_capability_level (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id        VARCHAR(64)  NOT NULL UNIQUE,
    skill_name      VARCHAR(128) NOT NULL,
    required_level  NUMERIC(5,2) NOT NULL CHECK (required_level  >= 0),
    available_level NUMERIC(5,2) NOT NULL CHECK (available_level >= 0),
    as_of           TIMESTAMPTZ  NOT NULL DEFAULT now()
);
```

### 23.2 Design notes

- **Derived values are not stored.** `growthPercent`, `direction`, `classification`, `gapValue`, `gapPercent`, `severity` and `isHighGrowth` are computed in Go (Section 15.3). Optionally cache them later (Redis, or a materialized view).
- `previousCount` is read from the immediately preceding window using a window function:

```sql
SELECT t.technology_id, t.technology_name, t.period_start, t.period_end,
       t.demand_count,
       COALESCE(LAG(t.demand_count) OVER (PARTITION BY t.technology_id ORDER BY t.period_start), 0) AS previous_count
FROM tta_technology_demand t
WHERE t.period_start >= $1;
```
- The IDs `technology_id` / `skill_id` are **soft references** to Skill Intelligence IDs. There is no cross-module foreign key; resolve names via the 4.2 API or events (Section 25).
- Use `pgx` or `database/sql` with parameterized queries only; no string-built SQL.
- Migrations are managed under `backend/database/` using the repo's migration tool (**To Be Confirmed**).
- Where does the demand signal come from? The source system is **To Be Confirmed** (README Section 36), so the ingestion path (events or batch) is an open item.

## 24. Future Redis Integration

Redis is already in the TalentIQ stack for API caching. Candidate uses for this module:

| Use | Key pattern | TTL | Notes |
|---|---|---|---|
| Cached endpoint payloads | `tta:v1:trends:{period}:{hash(query)}` | 5–15 min | Cache the serialized `data` after calculation |
| Dashboard aggregate | `tta:v1:dashboard:{period}` | 5 min | Most expensive aggregate |
| Emerging / high-growth lists | `tta:v1:emerging:{period}`, `tta:v1:highgrowth:{period}` | 10 min | |
| Capability gaps | `tta:v1:gaps` | 10 min | |
| Rate limiting | `tta:rl:{clientId}` | per window | Shared platform rate limiter |

### 24.1 Rules

- Cache lives in the **service layer**, as a decorator around the repository or service. Handlers stay unchanged.
- Cache **only successful** results. Never cache errors.
- Keys carry a version (`v1`) so the contract can evolve safely.
- Cache failures must **fail open**: if Redis is down, serve from the DB and log a warning.
- Invalidate on relevant Kafka events (Section 25) or rely on TTL initially.
- No PII is stored here, since this module has no personal data.

---

## 25. Future Kafka Integration

The project architecture routes talent events through Kafka to consumers including Trend Analysis. Event topic names and schemas are **not yet defined** in the project material (README Section 32), so this section lists *candidate* inputs only and does not invent final names.

### 25.1 Candidate consumed events (from the project's Kafka design)

| Event (project example list) | Possible use in Module 4.6 |
|---|---|
| `skill.updated` | Refresh skill names/metadata; adjust skill demand records |
| `technology.updated` | Refresh technology names/metadata for `tta_technology_demand` |
| `project.completed` | Count technology/skill usage as a demand signal |
| `candidate.profile.updated` | Update available-level inputs for capability gaps |
| `assessment.completed` | Update available-level inputs for capability gaps |
| `certification.added` | Update available-level inputs for capability gaps |
| `experience.updated` | Update technology/skill usage counts |

### 25.2 Consumer design (when approved)

```
Kafka topic → consumer group "technology-trend-analysis"
   → validate envelope + schema version
   → idempotent upsert into tta_* tables (keyed by event ID)
   → invalidate matching Redis keys
   → on failure: retry with backoff → dead-letter (observable, never silently dropped)
```

## 26. Testing Strategy

Frameworks are not yet fixed project-wide (README Section 36). The proposals below use Go's standard `testing` package (plus `net/http/httptest`) and Jest/React Testing Library for the frontend. **[ASSUMPTION]** Confirm the frontend runner.

### 26.1 Backend

| Level | Target | What to verify |
|---|---|---|
| Unit | `utils/trend_calculator.go` | All boundary values in 18.4 and 17.3. Table-driven. Target **100% coverage**. |
| Unit | `service/trend_service.go` | Filtering, sorting, pagination, tie-breaks, empty results, insights, invalid-data skipping (with a fake repository) |
| Unit | `dto` mappers | Exact JSON field names; empty slices serialize to `[]`; time format RFC3339 UTC |
| Handler / API | `handler` + `routes` via `httptest` | Status codes, envelope, validation errors, every query param, 404 on unknown ID |
| Contract | JSON golden files | Marshal DTOs and compare with the frontend interfaces (Section 26.3) |
| Repository | `TrendRepository` contract suite | Same tests run against the memory and (later) Postgres implementations |
| Integration | Router → service → in-memory repo | Full request on all 7 endpoints |
| Race | `go test -race ./...` | In-memory repo concurrency |

### 26.2 Frontend

| Level | Target | What to verify |
|---|---|---|
| Unit | `trendService.ts` | Unwraps `data.data`; correct URLs and params (mock axios) |
| Hook | `useTrends` family | loading → success, loading → error, refetch, abort on unmount |
| Component | 5 components | Renders each field; empty state; badge mapping per enum value |
| Page | Dashboard, Detail | Loading/error/empty/404 states; navigation by `technologyId` |

### 26.3 Contract test (prevents mismatch)

1. Backend golden-file tests marshal each DTO and assert exact key sets: `TechnologyTrend` has exactly 10 keys, `SkillTrend` 8, and `CapabilityGap` 8.
2. Frontend tests feed the **same golden JSON** into the components and the service, and type-check it against `trend.types.ts`.
3. Store the golden JSON in `tests/contracts/technology-trend-analysis/` (shared) so both sides test the same bytes.

## 27. Mock Data Strategy

Because the backend must be usable by the frontend immediately, Phase 1 ships an **in-memory repository** seeded with deterministic fixtures. The fixture values match Section 11 exactly, so the doc, the tests and the UI all agree.

### 27.1 Rules

- Mock data lives in `repository` (e.g. `trend_repository.go`, with the seed in an unexported function) and is replaced wholesale by the Postgres implementation later.
- The fixtures hold **raw inputs only** (counts, levels). Growth, classification and gap values are *computed* by the service, so the mock exercises the real logic.
- Fixtures are deterministic (no `rand`, no `time.Now()` in the data), so tests are stable. Only `generatedAt` and `health.timestamp` use the clock.

### 27.2 Fixture set

**Technologies** (current = Q3 2026, previous = Q2 2026):

| technologyId | name | demand | previous |
|---|---|---|---|
| `tech-go` | Go | 1240 | 780 |
| `tech-rust` | Rust | 410 | 260 |
| `tech-vector-db` | Vector Database | 380 | 190 |
| `tech-kubernetes` | Kubernetes | 2100 | 1750 |
| `tech-sql` | SQL | 1500 | 1450 |
| `tech-java` | Java | 900 | 1000 |
| `tech-jquery` | jQuery | 300 | 450 |

**History for `tech-go`** (Q4 2025 → Q3 2026 demand): `600 (prev window), 640, 700, 780, 1240` (as in 11.2).
Other technologies may return a 2-point history (previous + current) in v1.

**Skills:** `skill-vector-db` (380 / 190), `skill-cloud-security` (900 / 600), `skill-kafka` (1320 / 1000), `skill-kubernetes` (2100 / 1750), `skill-sql` (1500 / 1450).

**Capability levels:**

| skillId | required | available |
|---|---|---|
| `skill-cloud-security` | 5 | 1.5 |
| `skill-vector-db` | 4 | 1.6 |
| `skill-kubernetes` | 4 | 3 |
| `skill-kafka` | 4 | 3.2 |

### 27.3 Frontend mocking

Charan may develop against the real backend mock. For offline work, use MSW (or equivalent) with the same golden JSON from `tests/contracts/technology-trend-analysis/`, so frontend mocks cannot drift from the backend.

### 27.4 Coverage of the fixtures

Together the fixtures exercise every enum value: directions (`up`, `down`, `stable`), classifications (`emerging`, `growing`, `declining`, `stable`), severities (`low`, `medium`, `high`, `critical`), the boundary cases (Java at −10, Vector DB gap at 60, Kafka gap at 20), and the Q2 `up/stable` mismatch row.

---

## 28. Deployment Considerations

Cloud provider and deployment platform are **To Be Confirmed** (README Section 36). The project stack lists Docker and Kubernetes, so the module is written to be deployment-neutral:

| Concern | Guidance |
|---|---|
| Packaging | Part of the shared Go backend binary (`backend/cmd`), registered via `routes.RegisterRoutes`. No separate process in v1. |
| Config | Environment variables / config files only, with no secrets in source. Examples: `TTA_REPOSITORY`, `TTA_DEFAULT_PERIOD`, `TTA_HIGH_GROWTH_THRESHOLD` (optional override of constants). |
| Health | `GET /api/v1/technology-trend-analysis/health` for liveness and readiness. Readiness may call `repo.Ping` once a DB exists. |
| Statelessness | Safe to run multiple replicas. All state is in the DB (and Redis cache later). |
| Observability | Structured logs with request ID; Prometheus metrics (request count/latency per route); OpenTelemetry tracing across layers. |
| API docs | Add OpenAPI/Swagger annotations to handlers. The OpenAPI schema must match Section 8 exactly. |
| CORS | Configured at the shared API layer for the frontend origin. |
| Versioning | `/api/v1`. Breaking contract changes require `/api/v2`. Adding optional fields is non-breaking. |
| Security | Input validated against allow-lists. Parameterized SQL only. No sensitive data is logged. Auth middleware attached at the router group once the project's auth mechanism is confirmed. |
| Performance | Targets for the mock/DB-backed read endpoints: p95 < 300 ms uncached. Add Redis caching (Section 24) if exceeded. |
| CI/CD | Pipeline runs `go vet`, `go test -race ./...`, lint, frontend tests, and the contract tests (Section 26.3). |

## 29. Ownership Matrix

| Area | Owner | Reviewer / Consulted |
|---|---|---|
| Module 4.6 backend (all of `backend/internal/technology-trend-analysis/`) | **Indra Kalyan Reddy (Kalyan)** | Project/tech lead |
| Module 4.6 frontend (all of `frontend/src/modules/technology-trend-analysis/`) | **Charan** | Project/tech lead |
| `trend.types.ts` (contract source of truth) | Charan | Kalyan must approve changes |
| DTOs (`dto/*.go`) matching the contract | Kalyan | Charan verifies |
| Calculation rules (`trend_calculator.go`, `constants.go`) | Kalyan | Charan consulted on thresholds |
| Contract golden files (`tests/contracts/...`) | Kalyan + Charan (joint) | — |
| This `doc.md` | Kalyan (backend sections), Charan (frontend sections) | Project/tech lead |
| Skill/technology IDs and proficiency scale | Skill Intelligence (4.2): Vikas (BE), Naveendra (FE) | Kalyan consults |
| Shared API client / `baseURL` | Shared frontend | Charan consults |
| Shared infrastructure (PostgreSQL, Redis, Kafka) | Project-level (ownership **TBC**) | Kalyan consults |
| Event schemas (`skill.updated`, `technology.updated`, ...) | Owning module of each event | Kalyan consumes |

## 30. End-to-End Integration Flow

### 30.1 Sequence — dashboard

```
User -> TrendDashboardPage (Charan)
        useTrends / useEmerging / useHighGrowth / useCapabilityGaps
        |
        |  GET /api/v1/technology-trend-analysis/trends
        v
     Router -> Handler (Kalyan)
        validate query parameters
        Service: repository -> calculator -> filter/sort/paginate -> DTO
        |
        |  200 { data: TechnologyTrend[], pagination }
        v
     Frontend renders TrendCard

  Parallel calls, same pattern:
     GET /emerging       -> EmergingTechList
     GET /high-growth    -> HighGrowthSkillsList
     GET /capability-gaps -> CapabilityGapTable
```

### 30.2 Sequence — detail

```
User clicks TrendCard(technologyId = "tech-go")
 → /technology-trends/tech-go
 → GET /technology-trend-analysis/trends/tech-go
 → 200 { data: [Q4'25, Q1'26, Q2'26, Q3'26] }
 → TrendDetailPage header = last item; TrendChart = all items
```

### 30.3 Contract guarantees at the boundary

| Guarantee | Enforced by |
|---|---|
| JSON keys exactly match `trend.types.ts` | DTO structs + golden-file contract tests |
| Enum strings exactly as in the TS unions | Constants + tests |
| Lists never `null` | DTO mappers return non-nil slices |
| Derived values computed only in backend | Architecture rule (Section 15.3) |
| Standard envelope and error shape | `SuccessResponse` / `ErrorResponse` |

## 31. Open Items & Decisions

Items the source material left undefined, with the proposal used in this document. Confirm or change each item during the first review.

| # | Item | Proposal in this doc | Who confirms |
|---|---|---|---|
| 1 | `/dashboard` and `/health` are required by the brief but are not in the frontend's five-call list; the dashboard response shape is undefined | Dashboard shape in 8.4 / 11.6; health shape in 8.5. Frontend can ignore both for now. | Charan |
| 2 | Shape of `GET /trends/:technologyId` (no TS type given; brief says all responses are `data: [...]`) | Array of `TechnologyTrend`, ascending history (9.2) | Charan |
| 3 | `previousCount == 0` growth handling | `100.0` if `demandCount > 0`, else `0` (16.2) | Kalyan + Charan |
| 4 | High-growth threshold (`isHighGrowth`) | `growthPercent > 25` (18.3) | Kalyan + Charan |
| 5 | Severity band boundaries on decimals | `≤20, ≤40, ≤60, >60` (17.2) | Kalyan |
| 6 | Capability level scale | 0–5 assumed; logic is scale-independent (17.4) | Skill Intelligence owners |
| 7 | Surplus (`available > required`) | `gapValue = 0`, `low` (17.3) | Kalyan |
| 8 | `period` semantics (current vs previous window) | Equal-length consecutive windows (14.2) | Kalyan |
| 9 | Pagination on `/trends` | Optional `pagination` sibling of `data`; ignored by the current frontend | Charan |
| 10 | Frontend `baseURL` includes `/api/v1`; shared client file name and route paths | Assumed (9, 19.2) | Charan |
| 11 | Go module path and package naming | `github.com/MyHourly/TalentIQ/backend/...`; short package names (6) | Kalyan |
| 12 | Source system for the demand signal | To Be Confirmed project-wide (README Sec. 36); mock data until then | Project lead |
| 13 | Boundary with 4.5 Skill Gap Analysis for gap data | 4.6 exposes a trend-context view, 4.5 owns workforce gap analysis (1.2) | Kalyan + Manjunath |
| 14 | Auth mechanism | To Be Confirmed project-wide; middleware attach point reserved | Project lead |
| 15 | Name: brief says "Kalyan", README says "Indra Kalyan Reddy" | Treated as the same person | — |

---

*End of document.*