# Module 4.4 — AI / Semantic Talent Search

## API Reference (Frontend/Backend Implementation Contract)

**Project:** TalentIQ — Enterprise Talent Intelligence & Skill Graph Platform
**Module:** 4.4 AI / Semantic Talent Search
**Backend owner:** Hari Babu · **Frontend owner:** Anaiza
**Backend path:** `backend/internal/ai-semantic-talent-search/`
**Frontend path:** `frontend/src/modules/ai-semantic-talent-search/`
**API base path:** `/api/v1/ai-semantic-talent-search`
**Backend feature branch:** `feature/ai-semantic-talent-search`

This document is the **implementation-ready API contract** for Module 4.4. Every endpoint, field, status code, and error code below is a concrete implementation decision for this module's frontend/backend integration — a frontend developer can build against this document directly, without waiting on further clarification. It does not restate cross-cutting, organization-level decisions that remain open at the project level (e.g., the specific Vector DB product or embedding model/provider) — those are internal backend implementation details that this API contract intentionally does not expose, and they do not affect anything a frontend developer needs to build against.

---

## Table of Contents

1. [Scope of This Document](#1-scope-of-this-document)
2. [Ownership & Module Paths](#2-ownership--module-paths)
3. [Cross-Module Ownership Rules](#3-cross-module-ownership-rules)
4. [Authentication & Common Headers](#4-authentication--common-headers)
5. [Common Error Envelope](#5-common-error-envelope)
6. [`POST /search`](#6-post-apiv1ai-semantic-talent-searchsearch)
7. [`GET /filters`](#7-get-apiv1ai-semantic-talent-searchfilters)
8. [`GET /health`](#8-get-apiv1ai-semantic-talent-searchhealth)
9. [TypeScript Types](#9-typescript-types)
10. [API Service Implementation](#10-api-service-implementation)
11. [Frontend Module Folder Structure](#11-frontend-module-folder-structure)
12. [Backend Flow](#12-backend-flow)
13. [Frontend Flow](#13-frontend-flow)
14. [Frontend State Behavior](#14-frontend-state-behavior)
15. [What the Frontend Must Not Do](#15-what-the-frontend-must-not-do)
16. [End-to-End Flow](#16-end-to-end-flow)

---

## 1. Scope of This Document

This document covers the **three frontend-facing HTTP endpoints** Module 4.4 exposes, and everything a frontend developer needs to integrate against them: request/response shapes, validation rules, status codes, error codes, pagination, empty-result behavior, TypeScript types, an example service implementation, folder structure, and the expected frontend state behavior.

It does **not** cover: internal backend package layout, the specific embedding model/provider, the specific Vector DB product, event/message schemas used for re-indexing, or other module-internal implementation details. Those remain backend-internal and are out of scope for a frontend-facing API reference — a frontend developer does not need any of them to implement against this contract.

## 2. Ownership & Module Paths

| Item | Value |
|---|---|
| Module | 4.4 AI / Semantic Talent Search |
| Backend owner | Hari Babu |
| Frontend owner | Anaiza |
| Backend path | `backend/internal/ai-semantic-talent-search/` |
| Frontend path | `frontend/src/modules/ai-semantic-talent-search/` |
| API base path | `/api/v1/ai-semantic-talent-search` |
| Feature branch | `feature/ai-semantic-talent-search` |
| Integration branch | `develop` |
| Stable branch | `main` |

## 3. Cross-Module Ownership Rules

Module 4.4 **consumes** data from upstream modules and **produces** ranked search output. It does not become the source of truth for any data it did not originate.

| Data / concern | Owned by | Module 4.4's relationship to it |
|---|---|---|
| Talent/profile records | 4.1 Talent Profile Intelligence | Consumed and indexed for search; 4.4 does not edit or store the authoritative record. |
| Skill definitions / taxonomy | 4.2 Skill Intelligence | Consumed as search signals and as the basis for the `skills` filter; 4.4 does not define or edit the skill taxonomy. |
| Skill/talent/technology graph relationships | 4.3 Skill Graph | Consumed as relevance-enriching context; 4.4 does not own graph relationships. |
| Final candidate recommendation scoring | 4.7 Talent Recommendation Engine | 4.4 produces `relevanceScore` / similarity output that 4.7 may consume as one input; 4.4 does not itself produce final recommendations. |
| Search ranking, query interpretation, embedding, semantic retrieval, indexing hooks | **4.4 (this module)** | Owned entirely by this module. No other module reimplements this logic. |

**Source-domain ownership stays with the owning module.** The `GET /filters` endpoint (Section 7) returns `skills`, `locations`, and `jobTitles` values — these are **derived from Module 4.4's own search index**, not republished as authoritative master data. The canonical skill taxonomy remains owned by 4.2; the canonical profile/location/title data remains owned by 4.1. If a consumer needs the authoritative master list (rather than "currently searchable values"), it must go to the owning module's own API, not to Module 4.4.

**Indexing and event processing are backend/internal.** Keeping 4.4's search index current as upstream data changes in 4.1/4.2/4.3 is a backend-internal concern (event consumption, re-embedding, re-indexing). It is **not** exposed through any endpoint in this document, has no frontend-facing surface, and is out of scope for this API reference.

---

## 4. Authentication & Common Headers

Authentication is platform-wide, not module-specific: Module 4.4 does not define its own authentication scheme, it enforces the same bearer-token check every other TalentIQ module's API enforces.

| Header | Required on | Value |
|---|---|---|
| `Authorization` | `POST /search`, `GET /filters` | `Bearer <token>` — a valid platform-issued access token. |
| `Content-Type` | `POST /search` | `application/json` |
| `Accept` | All endpoints | `application/json` |
| `X-Request-Id` | All endpoints (optional) | Client-supplied correlation ID (UUID recommended). If omitted, the server generates one. The value in effect for the request is always echoed back as `requestId` in an error body, and is recommended for correlating logs across the frontend and backend. |

**`GET /health` is unauthenticated.** It does not require `Authorization` and must not be gated behind the platform's auth check, since it exists for infrastructure health checks (container orchestration liveness/readiness probes) and operational monitoring, which run without a user session. See [Section 8](#8-get-apiv1ai-semantic-talent-searchhealth).

## 5. Common Error Envelope

`POST /search` and `GET /filters` use one consistent error shape for every non-2xx response:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {},
    "requestId": "string"
  }
}
```

| Field | Type | Description |
|---|---|---|
| `error.code` | string | One of the fixed codes below. Stable and safe to switch on in frontend code. |
| `error.message` | string | Human-readable description of the failure. Suitable for logging; not guaranteed to be end-user-friendly wording, so the frontend maps `code` to its own user-facing copy (see [Section 14](#14-frontend-state-behavior)). |
| `error.details` | object | Code-specific structured detail (e.g., field-level validation messages). May be an empty object `{}` when there is nothing further to add. |
| `error.requestId` | string | Correlation ID for this request, for cross-referencing frontend/backend logs. |

**Shared error codes** (an endpoint's own section states which of these it actually returns):

| `code` | HTTP status | Meaning |
|---|---|---|
| `VALIDATION_ERROR` | 400 | The request body failed validation. `details.fields` is an object mapping field paths to a human-readable validation message for that field. |
| `UNAUTHORIZED` | 401 | The request is missing a bearer token, or the token is invalid or expired. |
| `RATE_LIMIT_EXCEEDED` | 429 | The caller has exceeded the allowed request rate for this endpoint. `details.retryAfterSeconds` gives the minimum wait before retrying. |
| `DEPENDENCY_UNAVAILABLE` | 503 | A required upstream dependency (the module's semantic retrieval store, its embedding service, or an upstream module's data) is temporarily unavailable. Safe to retry with backoff. |
| `INTERNAL_SERVER_ERROR` | 500 | An unexpected failure occurred inside the module's own pipeline. |

`GET /health` does **not** use this envelope — see [Section 8](#8-get-apiv1ai-semantic-talent-searchhealth) for its response shape.

---

## 6. `POST /api/v1/ai-semantic-talent-search/search`

### Purpose

Accepts a natural-language query plus optional structured filters and pagination, and returns a ranked page of talent results produced by Module 4.4's semantic search pipeline. This is the module's primary endpoint and the only one that runs the full search pipeline (see [Section 12](#12-backend-flow)).

### Authentication & Headers

| Header | Required | Value |
|---|---|---|
| `Authorization` | Yes | `Bearer <token>` |
| `Content-Type` | Yes | `application/json` |
| `Accept` | Recommended | `application/json` |
| `X-Request-Id` | Optional | Client correlation ID |

### Path Parameters

None.

### Query Parameters

None. All inputs are supplied in the JSON request body.

### Request Body

```json
{
  "query": "distributed systems engineer with Kafka and Go experience",
  "filters": {
    "skills": ["Go", "Kafka"],
    "minExperienceYears": 3,
    "maxExperienceYears": 10,
    "locations": ["Bengaluru, India"],
    "jobTitles": ["Backend Engineer", "Senior Backend Engineer"]
  },
  "page": 1,
  "pageSize": 20
}
```

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `query` | string | **Yes** | — | The natural-language search query. |
| `filters` | object | No | `{}` | Structured filters, applied in addition to the semantic query. |
| `filters.skills` | string[] | No | none applied | Restrict results to talent whose indexed skills include at least one listed skill. |
| `filters.minExperienceYears` | number | No | none applied | Minimum total years of experience, inclusive. |
| `filters.maxExperienceYears` | number | No | none applied | Maximum total years of experience, inclusive. |
| `filters.locations` | string[] | No | none applied | Restrict results to talent located in one of the listed locations. |
| `filters.jobTitles` | string[] | No | none applied | Restrict results to talent whose current title is one of the listed titles. |
| `page` | integer | No | `1` | 1-indexed page number. |
| `pageSize` | integer | No | `20` | Number of results per page. |

### Required Fields

`query` is the only required field. `filters`, `page`, and `pageSize` are all optional.

### Optional Fields

`filters` (and every field inside it), `page`, `pageSize`.

### Validation Rules

| Field | Rule |
|---|---|
| `query` | Required. Must be a non-empty string after trimming whitespace. Maximum 500 characters. |
| `filters.skills` | If present, must be an array of non-empty strings. Maximum 20 entries. |
| `filters.minExperienceYears` | If present, must be a number ≥ 0. |
| `filters.maxExperienceYears` | If present, must be a number ≥ 0. If both `minExperienceYears` and `maxExperienceYears` are present, `maxExperienceYears` must be ≥ `minExperienceYears`. |
| `filters.locations` | If present, must be an array of non-empty strings. Maximum 20 entries. |
| `filters.jobTitles` | If present, must be an array of non-empty strings. Maximum 20 entries. |
| `page` | If present, must be an integer ≥ 1. |
| `pageSize` | If present, must be an integer between 1 and 100 inclusive. |
| Unrecognized fields | Ignored (not rejected), so additive, backward-compatible client changes don't break requests. |

All violations are returned together in a single `VALIDATION_ERROR` response (see below) — the frontend does not need to fix one field, resubmit, and discover the next.

### Success Response

**HTTP status:** `200 OK`

```json
{
  "results": [
    {
      "talentId": "TAL-10432",
      "name": "Ananya Rao",
      "currentTitle": "Senior Backend Engineer",
      "location": "Bengaluru, India",
      "experienceYears": 6,
      "skills": ["Go", "Kafka", "Kubernetes", "gRPC"],
      "relevanceScore": 0.91,
      "matchSummary": "Strong match on distributed systems, event-driven architecture with Kafka, and Go-based microservices.",
      "profileUrl": "/talent-profiles/TAL-10432"
    },
    {
      "talentId": "TAL-20981",
      "name": "Rohit Menon",
      "currentTitle": "Backend Engineer",
      "location": "Hyderabad, India",
      "experienceYears": 4,
      "skills": ["Go", "PostgreSQL", "Kafka"],
      "relevanceScore": 0.78,
      "matchSummary": "Relevant Go and Kafka experience; less direct evidence of distributed-systems architecture ownership.",
      "profileUrl": "/talent-profiles/TAL-20981"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 137,
    "totalPages": 7
  },
  "meta": {
    "query": "distributed systems engineer with Kafka and Go experience",
    "executionTimeMs": 184
  }
}
```

### Response Field Descriptions

| Field | Type | Description |
|---|---|---|
| `results` | array | Ranked talent results for the requested page, highest relevance first. |
| `results[].talentId` | string | Stable identifier for the talent record, owned by Module 4.1. |
| `results[].name` | string | Talent's display name. |
| `results[].currentTitle` | string | Talent's current job title. |
| `results[].location` | string | Talent's current location. |
| `results[].experienceYears` | number | Total years of professional experience. |
| `results[].skills` | string[] | Skills indexed for this talent record that are relevant to display alongside the result. |
| `results[].relevanceScore` | number | Relevance of this result to the query, in the range `0.0`–`1.0`, combining semantic similarity with applicable filters/business rules. Higher is more relevant. |
| `results[].matchSummary` | string | Short, human-readable explanation of why this result matched the query. |
| `results[].profileUrl` | string | Frontend route (client-side navigation path) into this talent's profile page, owned by Module 4.1. This is a frontend route, not a backend API path. |
| `pagination.page` | integer | The page number this response represents. |
| `pagination.pageSize` | integer | The page size used to produce this response. |
| `pagination.total` | integer | Total number of matching results across all pages. |
| `pagination.totalPages` | integer | Total number of pages, computed as `ceil(total / pageSize)`. |
| `meta.query` | string | The query string this response was generated for (echoed back for display/debugging). |
| `meta.executionTimeMs` | integer | Server-side time, in milliseconds, taken to process the search request end to end. |

### Error Responses

| HTTP status | `error.code` | When it happens |
|---|---|---|
| 400 | `VALIDATION_ERROR` | The request body fails one or more rules in [Validation Rules](#validation-rules) above. |
| 401 | `UNAUTHORIZED` | Missing or invalid `Authorization` header. |
| 429 | `RATE_LIMIT_EXCEEDED` | Too many search requests from the caller in the current window. |
| 503 | `DEPENDENCY_UNAVAILABLE` | The semantic retrieval store, the embedding service, or an upstream module's data required to serve the request is temporarily unreachable. |
| 500 | `INTERNAL_SERVER_ERROR` | An unexpected failure inside query interpretation, filtering, ranking, or response mapping. |

**Example error response** (`400 Bad Request`):

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request could not be validated.",
    "details": {
      "fields": {
        "query": "query is required and must be a non-empty string",
        "filters.maxExperienceYears": "must be greater than or equal to filters.minExperienceYears"
      }
    },
    "requestId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
}
```

### Pagination Behavior

- `page` and `pageSize` in the request control which slice of the full result set is returned.
- `pagination.total` is the count of all results matching the query and filters, across every page.
- `pagination.totalPages` is `ceil(total / pageSize)`; `0` when `total` is `0`.
- Requesting a `page` beyond `totalPages` is **not an error** — it returns `200 OK` with `results: []` and `pagination` reflecting the true `total`/`totalPages`, so the frontend can safely clamp its own page state against the response rather than needing to pre-validate the page number.

### Empty-Result Behavior

When no talent matches the query and filters, the response is still `200 OK`:

```json
{
  "results": [],
  "pagination": { "page": 1, "pageSize": 20, "total": 0, "totalPages": 0 },
  "meta": { "query": "quantum-computing hardware architect", "executionTimeMs": 41 }
}
```

An empty result set is a normal, successful outcome, never an error response.

### Frontend Usage

Called when the user submits a query from `SearchInput` (optionally combined with `SearchFilters` selections) and whenever `SearchPagination` requests a different page of the same query/filters. Consumed by `SearchResults` and rendered as a list of `TalentResultCard` components (see [Section 11](#11-frontend-module-folder-structure)).

### Backend Responsibility

Owns the entire pipeline for this request: validation, query interpretation, embedding, semantic retrieval, filtering, ranking, and response mapping (see [Section 12](#12-backend-flow)). Consumes profile/skill/graph data from Modules 4.1–4.3 without owning it (see [Section 3](#3-cross-module-ownership-rules)).

### Example Frontend Request

```typescript
import { searchTalent } from '../services/searchApi';

const response = await searchTalent({
  query: 'distributed systems engineer with Kafka and Go experience',
  filters: {
    skills: ['Go', 'Kafka'],
    minExperienceYears: 3,
    maxExperienceYears: 10,
    locations: ['Bengaluru, India'],
    jobTitles: ['Backend Engineer', 'Senior Backend Engineer'],
  },
  page: 1,
  pageSize: 20,
});
```

### Example Backend Response

See [Success Response](#success-response) above — that JSON is the complete, exact shape returned for this example request.

### End-to-End Flow

See [Section 16](#16-end-to-end-flow) for the full request path from `SearchInput` through the backend pipeline to `TalentResultCard`.

---

## 7. `GET /api/v1/ai-semantic-talent-search/filters`

### Purpose

Returns the set of filterable fields Module 4.4 supports, together with the values/bounds currently available in its own search index, so the frontend can render filter controls (skill multi-select, location multi-select, job-title multi-select, experience-years range) with real, current options instead of hard-coded lists.

### Authentication & Headers

| Header | Required | Value |
|---|---|---|
| `Authorization` | Yes | `Bearer <token>` |
| `Accept` | Recommended | `application/json` |
| `X-Request-Id` | Optional | Client correlation ID |

### Path Parameters

None.

### Query Parameters

None.

### Request Body

None — this is a `GET` request with no body.

### Required Fields / Optional Fields / Data Types (Request)

Not applicable — this endpoint takes no input beyond its required headers.

### Validation Rules

None beyond the standard authentication check in [Section 4](#4-authentication--common-headers).

### Success Response

**HTTP status:** `200 OK`

```json
{
  "skills": ["Go", "Java", "Kafka", "Kubernetes", "React", "PostgreSQL", "gRPC", "Python", "TypeScript"],
  "locations": ["Bengaluru, India", "Hyderabad, India", "Pune, India", "Remote"],
  "jobTitles": ["Backend Engineer", "Senior Backend Engineer", "Frontend Engineer", "Site Reliability Engineer", "Data Engineer"],
  "experienceYears": { "min": 0, "max": 30 }
}
```

### Response Field Descriptions

| Field | Type | Description |
|---|---|---|
| `skills` | string[] | Distinct skill values currently present in Module 4.4's search index, suitable for populating the `filters.skills` control. |
| `locations` | string[] | Distinct location values currently present in the index, suitable for `filters.locations`. |
| `jobTitles` | string[] | Distinct job-title values currently present in the index, suitable for `filters.jobTitles`. |
| `experienceYears.min` | number | Lowest `experienceYears` value currently present in the index. |
| `experienceYears.max` | number | Highest `experienceYears` value currently present in the index. |

If the number of distinct values for `skills`, `locations`, or `jobTitles` exceeds 200, this endpoint returns only the 200 most frequent values for that field, so the payload stays bounded regardless of index size.

### Error Responses

| HTTP status | `error.code` | When it happens |
|---|---|---|
| 401 | `UNAUTHORIZED` | Missing or invalid `Authorization` header. |
| 429 | `RATE_LIMIT_EXCEEDED` | Too many requests from the caller in the current window. |
| 503 | `DEPENDENCY_UNAVAILABLE` | The search index backing this endpoint is temporarily unreachable. |
| 500 | `INTERNAL_SERVER_ERROR` | An unexpected failure occurred while assembling the response. |

**Example error response** (`401 Unauthorized`):

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "A valid bearer token is required.",
    "details": {},
    "requestId": "f0e1d2c3-b4a5-6789-0123-456789abcdef"
  }
}
```

### Pagination Behavior

Not applicable. This endpoint always returns the complete (bounded, per the 200-value cap above) set of filter options in a single response; it is not paginated.

### Empty-Result Behavior

If the search index has no data yet, the response is still `200 OK`, with empty arrays and a zeroed range:

```json
{
  "skills": [],
  "locations": [],
  "jobTitles": [],
  "experienceYears": { "min": 0, "max": 0 }
}
```

### Frontend Usage

Called once when `SearchPage` (or its `SearchFilters` component) mounts, to populate the filter controls before the user runs their first search. Not called as part of every `POST /search` request.

### Backend Responsibility

Derive current filter facets from Module 4.4's own index — not from 4.1/4.2's master data directly — and serve this endpoint from a short-lived cache (Redis, per the project's shared-infrastructure direction) so it stays cheap to call on every page load.

### Example Frontend Request

```typescript
import { getSearchFilters } from '../services/searchApi';

const filterOptions = await getSearchFilters();
```

### Example Backend Response

See [Success Response](#success-response-1) above — that JSON is the complete, exact shape returned.

### End-to-End Flow

See [Section 16](#16-end-to-end-flow). This endpoint runs before the main search flow, to populate `SearchFilters`.

---

## 8. `GET /api/v1/ai-semantic-talent-search/health`

### Purpose

Reports whether Module 4.4's service and its critical dependencies are operational. Used by infrastructure health checks (container orchestration liveness/readiness probes) and by operational monitoring — not by the normal user-facing search flow.

### Authentication & Headers

| Header | Required | Value |
|---|---|---|
| `Authorization` | **No** | Not required — see [Section 4](#4-authentication--common-headers). |
| `Accept` | Recommended | `application/json` |

### Path Parameters

None.

### Query Parameters

None.

### Request Body

None — this is a `GET` request with no body.

### Required Fields / Optional Fields / Data Types (Request)

Not applicable — this endpoint takes no input.

### Validation Rules

None — there is no input to validate, and no authentication check to fail.

### Success Response

**HTTP status:** `200 OK` (returned when the module and its dependencies are healthy)

```json
{
  "status": "ok",
  "module": "ai-semantic-talent-search",
  "dependencies": {
    "semanticRetrievalStore": "ok",
    "embeddingService": "ok"
  },
  "timestamp": "2026-09-29T10:15:32Z"
}
```

### Response Field Descriptions

| Field | Type | Description |
|---|---|---|
| `status` | string | `"ok"` or `"unhealthy"` — overall status for this module. |
| `module` | string | Fixed value `"ai-semantic-talent-search"`, identifying which module this health report is for. |
| `dependencies` | object | Per-dependency status. Keys are dependency names; values are `"ok"` or `"unavailable"`. |
| `dependencies.semanticRetrievalStore` | string | Whether the module's semantic retrieval store responded to a lightweight reachability check. |
| `dependencies.embeddingService` | string | Whether the module's embedding service responded to a lightweight reachability check. |
| `timestamp` | string | ISO-8601 UTC timestamp for when this health check ran. |

### Error Responses / Unhealthy Response

This endpoint does **not** use the [common error envelope](#5-common-error-envelope) — a health report is a status, not a failed operation. When one or more dependencies are unavailable, it returns:

**HTTP status:** `503 Service Unavailable`

```json
{
  "status": "unhealthy",
  "module": "ai-semantic-talent-search",
  "dependencies": {
    "semanticRetrievalStore": "unavailable",
    "embeddingService": "ok"
  },
  "timestamp": "2026-09-29T10:16:04Z"
}
```

There is no `400`/`401`/`429` case for this endpoint: it takes no input to fail validation on, requires no authentication, and is exempt from rate limiting so that frequent, predictable infrastructure probes are never throttled.

### Pagination Behavior

Not applicable.

### Empty-Result Behavior

Not applicable — this endpoint always returns a status report, never an empty result set.

### Frontend Usage

**Not called by `SearchPage` or any of its components as part of the normal user-facing search flow.** This endpoint exists for infrastructure and operations tooling (deployment health probes, status dashboards), not for the search UI.

### Backend Responsibility

Run a lightweight, short-timeout reachability check against each critical dependency (the semantic retrieval store, the embedding service) and report aggregate plus per-dependency status, without executing a full search.

### Example Frontend Request

Not part of the search UI's request flow; included here for completeness, as this is still a request an ops tool or deployment probe would make against this module:

```typescript
const response = await fetch('/api/v1/ai-semantic-talent-search/health', {
  headers: { Accept: 'application/json' },
});
const health = await response.json();
```

### Example Backend Response

See [Success Response](#success-response-2) and the unhealthy example above — those cover the two possible response bodies.

### End-to-End Flow

This endpoint bypasses the search pipeline entirely; it is not part of the flow in [Section 16](#16-end-to-end-flow).

---

## 9. TypeScript Types

```typescript
// frontend/src/modules/ai-semantic-talent-search/types/index.ts

// ---- Request types ----

export interface SearchFilters {
  skills?: string[];
  minExperienceYears?: number;
  maxExperienceYears?: number;
  locations?: string[];
  jobTitles?: string[];
}

export interface SearchRequest {
  query: string;
  filters?: SearchFilters;
  page?: number;
  pageSize?: number;
}

// ---- Response types: POST /search ----

export interface TalentSearchResult {
  talentId: string;
  name: string;
  currentTitle: string;
  location: string;
  experienceYears: number;
  skills: string[];
  relevanceScore: number;
  matchSummary: string;
  profileUrl: string;
}

export interface PaginationMeta {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface SearchMeta {
  query: string;
  executionTimeMs: number;
}

export interface SearchResponse {
  results: TalentSearchResult[];
  pagination: PaginationMeta;
  meta: SearchMeta;
}

// ---- Response type: GET /filters ----

export interface ExperienceYearsRange {
  min: number;
  max: number;
}

export interface FilterOptionsResponse {
  skills: string[];
  locations: string[];
  jobTitles: string[];
  experienceYears: ExperienceYearsRange;
}

// ---- Response type: GET /health ----

export type DependencyStatus = 'ok' | 'unavailable';

export interface HealthResponse {
  status: 'ok' | 'unhealthy';
  module: string;
  dependencies: Record<string, DependencyStatus>;
  timestamp: string;
}

// ---- Shared error type ----

export type SearchErrorCode =
  | 'VALIDATION_ERROR'
  | 'UNAUTHORIZED'
  | 'RATE_LIMIT_EXCEEDED'
  | 'DEPENDENCY_UNAVAILABLE'
  | 'INTERNAL_SERVER_ERROR';

export interface ApiErrorBody {
  error: {
    code: SearchErrorCode;
    message: string;
    details: Record<string, unknown>;
    requestId: string;
  };
}
```

## 10. API Service Implementation

```typescript
// frontend/src/modules/ai-semantic-talent-search/services/searchApi.ts

import type {
  SearchRequest,
  SearchResponse,
  FilterOptionsResponse,
  HealthResponse,
  ApiErrorBody,
  SearchErrorCode,
} from '../types';

const BASE_PATH = '/api/v1/ai-semantic-talent-search';

/** Thrown for every non-2xx response from this module's API. */
export class SearchApiError extends Error {
  readonly status: number;
  readonly code: SearchErrorCode;
  readonly details: Record<string, unknown>;
  readonly requestId: string;

  constructor(status: number, body: ApiErrorBody) {
    super(body.error.message);
    this.name = 'SearchApiError';
    this.status = status;
    this.code = body.error.code;
    this.details = body.error.details;
    this.requestId = body.error.requestId;
  }
}

/**
 * This module's requests go through the platform's shared HTTP client in
 * `frontend/src/services/`, which already attaches `Authorization` and
 * `X-Request-Id`. `httpClient` below stands in for that shared client so
 * this file is self-contained; in the real module it is imported, not
 * redefined here.
 */
async function httpClient<T>(path: string, init: RequestInit): Promise<T> {
  const response = await fetch(`${BASE_PATH}${path}`, {
    ...init,
    headers: {
      Accept: 'application/json',
      ...(init.body ? { 'Content-Type': 'application/json' } : {}),
      ...(init.headers ?? {}),
    },
  });

  if (!response.ok) {
    const body = (await response.json()) as ApiErrorBody;
    throw new SearchApiError(response.status, body);
  }

  return response.json() as Promise<T>;
}

export function searchTalent(payload: SearchRequest): Promise<SearchResponse> {
  return httpClient<SearchResponse>('/search', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function getSearchFilters(): Promise<FilterOptionsResponse> {
  return httpClient<FilterOptionsResponse>('/filters', {
    method: 'GET',
  });
}

export function getSearchHealth(): Promise<HealthResponse> {
  return httpClient<HealthResponse>('/health', {
    method: 'GET',
  });
}
```

---

## 11. Frontend Module Folder Structure

```
frontend/src/modules/ai-semantic-talent-search/
├── pages/
│   └── SearchPage.tsx            → top-level page for this module's route
├── components/
│   ├── SearchInput.tsx           → natural-language query input
│   ├── SearchFilters.tsx         → skills / experience / locations / job titles controls
│   ├── SearchResults.tsx         → renders the ranked result list
│   ├── TalentResultCard.tsx      → a single result: name, title, skills, relevanceScore, matchSummary
│   ├── SearchPagination.tsx      → paging controls driven by `pagination`
│   ├── SearchLoadingState.tsx    → shown while a search request is in flight
│   ├── SearchEmptyState.tsx      → shown when `results` is empty
│   └── SearchErrorState.tsx      → shown for non-validation errors
├── hooks/
│   ├── useTalentSearch.ts        → owns query/filters/page state and calls searchApi
│   └── useSearchFilterOptions.ts → loads GET /filters once on mount
├── services/
│   └── searchApi.ts              → this module's typed API client (Section 10)
├── types/
│   └── index.ts                  → this module's TypeScript types (Section 9)
└── index.ts                      → module entry point / route registration
```

This module consumes, but does not redefine, the shared building blocks in `frontend/src/components/`, `frontend/src/hooks/`, `frontend/src/services/`, `frontend/src/types/`, and `frontend/src/utils/` (per the TalentIQ frontend architecture) — for example, the shared HTTP client that attaches `Authorization` and `X-Request-Id` lives once in `frontend/src/services/`, not duplicated inside `searchApi.ts`.

## 12. Backend Flow

```
Request → Validation → Query Interpretation → Embedding → Semantic Retrieval → Filtering → Ranking → Response Mapping → Response
```

| Stage | Responsibility |
|---|---|
| **Request** | `POST /search` is received; `Authorization` is checked (see [Section 4](#4-authentication--common-headers)). |
| **Validation** | The request body is checked against the rules in [Section 6](#6-post-apiv1ai-semantic-talent-searchsearch). Any failure short-circuits the pipeline and returns `VALIDATION_ERROR` immediately — no later stage runs. |
| **Query Interpretation** | `query` is parsed to extract search intent — the concepts, skills, and roles implied by the natural-language text — independent of `filters`, which are applied later, unmodified, as hard constraints. |
| **Embedding** | The interpreted query is converted into a vector representation for semantic comparison. |
| **Semantic Retrieval** | The query vector is compared against the module's semantic retrieval store to retrieve a candidate set of talent records ranked by conceptual similarity. |
| **Filtering** | `filters.skills`, `filters.minExperienceYears`, `filters.maxExperienceYears`, `filters.locations`, and `filters.jobTitles` are applied as hard constraints against the candidate set, narrowing it. |
| **Ranking** | The filtered candidates are ordered by combining semantic similarity with applicable business rules, producing each result's `relevanceScore`. |
| **Response Mapping** | Ranked candidates are shaped into `TalentSearchResult` objects (including `matchSummary` and `profileUrl`), and `pagination`/`meta` are assembled for the requested `page`/`pageSize`. |
| **Response** | The `SearchResponse` JSON (Section 6) is returned as `200 OK`. |

## 13. Frontend Flow

```
Search Page → Search Input → Filters → API Service → Search API → Results → Talent Result Card → Talent Profile
```

| Step | Component / concern | What happens |
|---|---|---|
| **Search Page** | `SearchPage` | Hosts the search experience for this module; owns overall page state via `useTalentSearch`. |
| **Search Input** | `SearchInput` | Captures the user's natural-language query text. |
| **Filters** | `SearchFilters` | Lets the user select `skills` / `minExperienceYears` / `maxExperienceYears` / `locations` / `jobTitles`, populated from `GET /filters` (via `useSearchFilterOptions`). |
| **API Service** | `searchApi.ts` | Builds and sends the `SearchRequest`; this is the only place in the module that calls `fetch`. |
| **Search API** | `POST /search` (backend) | Runs the pipeline in [Section 12](#12-backend-flow) and returns a `SearchResponse` or an error. |
| **Results** | `SearchResults` | Receives `results`/`pagination`/`meta` (or an error/empty state — see [Section 14](#14-frontend-state-behavior)) and renders the page of results plus `SearchPagination`. |
| **Talent Result Card** | `TalentResultCard` | Renders one `TalentSearchResult`: `name`, `currentTitle`, `location`, `experienceYears`, `skills`, `relevanceScore`, `matchSummary`. |
| **Talent Profile** | (Module 4.1) | Clicking a `TalentResultCard` navigates the app router to `results[].profileUrl`, handing off into Module 4.1's UI. Module 4.4 does not render the profile page itself. |

---

## 14. Frontend State Behavior

| State | When | What the frontend displays |
|---|---|---|
| **Initial** | `SearchPage` has mounted; the user has not yet submitted a query. | `SearchInput` empty; `SearchFilters` populated from `GET /filters` but with no selections made; no `SearchResults`, `SearchPagination`, loading, empty, or error UI rendered. A prompt such as "Enter a search query to find talent" may be shown in place of results. `POST /search` has not been called yet. |
| **Loading** | A `POST /search` request is in flight (initial submit, a filter change that re-runs the search, or a pagination change). | `SearchLoadingState` is shown in place of `SearchResults`. `SearchInput`/`SearchFilters` remain visible but should be treated as read-only/disabled for the duration of the request to prevent duplicate submissions. `SearchPagination` is hidden. |
| **Success** | `POST /search` resolves with `results.length > 0`. | `SearchResults` renders one `TalentResultCard` per entry in `results`, in the order returned. `SearchPagination` is shown, reflecting `pagination.page`/`pagination.totalPages`. `meta.executionTimeMs` may optionally be shown as a small "results in Xms" indicator. |
| **Empty** | `POST /search` resolves with `results.length === 0` (see [Empty-Result Behavior](#empty-result-behavior) in Section 6) — either no matches exist, or the requested `page` is beyond `pagination.totalPages`. | `SearchEmptyState` is shown instead of `SearchResults`, with a message indicating no matches and, where applicable, a suggestion to broaden the filters or return to page 1. `SearchPagination` is hidden or shown disabled. |
| **Error** | `searchApi` throws a `SearchApiError`. | Handling depends on `error.code`: `VALIDATION_ERROR` is shown **inline**, next to the offending field(s) in `SearchInput`/`SearchFilters` (using `details.fields`), not as a full error state, since it is the user's own input that needs correcting. `UNAUTHORIZED` redirects into the platform's re-authentication flow rather than showing a search-specific error. `RATE_LIMIT_EXCEEDED`, `DEPENDENCY_UNAVAILABLE`, and `INTERNAL_SERVER_ERROR` all render `SearchErrorState` with a code-specific, user-facing message and a retry action; `RATE_LIMIT_EXCEEDED`'s retry is disabled until `details.retryAfterSeconds` has elapsed. |

## 15. What the Frontend Must Not Do

Module 4.4's search logic is owned entirely by the backend. The frontend consumes the contract in this document — it must **not**:

- **Calculate semantic scores.** `relevanceScore` is computed server-side, during Ranking (Section 12); the frontend only renders the value it receives.
- **Perform vector retrieval.** Comparing embeddings against the semantic retrieval store happens only in the backend's Semantic Retrieval stage.
- **Generate embeddings.** Converting query or talent text into vector representations happens only in the backend's Embedding stage.
- **Implement backend ranking.** Ordering, scoring, or re-scoring results client-side is not permitted — the order returned in `results` is final.
- **Access databases directly.** The frontend calls only the three endpoints in this document; it never queries Module 4.4's index, PostgreSQL, or the semantic retrieval store directly.
- **Duplicate backend business logic.** Filtering, validation beyond basic input hygiene (e.g., trimming whitespace before display), and matching logic all live in the backend; the frontend reflects the backend's decisions, it does not re-derive them.

Violating any of the above breaks the module boundary described in [Section 3](#3-cross-module-ownership-rules) and the project's general module boundary rules.

## 16. End-to-End Flow

Combining the frontend flow (Section 13) and backend flow (Section 12) into one path, for the primary search endpoint:

```
User types a query into Search Input
        ↓
User optionally selects Filters (populated earlier from GET /filters)
        ↓
API Service (searchApi.ts) builds a SearchRequest and calls POST /search
        ↓
Backend: Request → Validation → Query Interpretation → Embedding →
         Semantic Retrieval → Filtering → Ranking → Response Mapping → Response
        ↓
API Service resolves with a SearchResponse, or throws a SearchApiError
        ↓
Search Page transitions to the Success, Empty, or Error state (Section 14)
        ↓
Results render as Talent Result Cards, paged via Search Pagination
        ↓
User clicks a Talent Result Card → app router navigates to profileUrl
        ↓
Talent Profile page (Module 4.1) — outside Module 4.4's scope
```

---

*This document is the implementation-ready API reference for Module 4.4 — AI / Semantic Talent Search. It supersedes the earlier API Contract, Error Handling, and Search Filters sections of the module documentation. Cross-cutting organization-level decisions such as the Vector DB product, embedding model/provider, authentication provider, and final team allocation remain governed by the project's shared documentation and are not exposed through this frontend API contract.*
