---
name: Dashboard and Import Features
overview: "Add two new features to PocketBase: (1) A Real-time Metrics Dashboard with visual charts for performance monitoring, and (2) A Data Import Wizard for bulk importing CSV/JSON data into collections."
todos:
  - id: metrics-api
    content: Create apis/metrics.go with overview, requests, latency, errors, endpoints, collections endpoints
    status: pending
  - id: metrics-core
    content: Add metrics query functions in core (LogsMetricsStats, DatabaseSize)
    status: pending
  - id: dashboard-page
    content: Create Dashboard.svelte main page with layout and auto-refresh
    status: pending
  - id: dashboard-components
    content: "Build dashboard components: MetricCard, RequestsChart, LatencyChart, EndpointsChart, CollectionsTable"
    status: pending
  - id: dashboard-store
    content: Create dashboard.js store for state management
    status: pending
  - id: import-api
    content: Create apis/import.go with preview, validate, execute endpoints
    status: pending
  - id: import-page
    content: Create ImportWizard.svelte with 4-step wizard flow
    status: pending
  - id: import-components
    content: "Build import components: FileUpload, DataPreview, FieldMapper, ImportProgress"
    status: pending
  - id: import-store
    content: Create import.js store for wizard state
    status: pending
  - id: routes-sidebar
    content: Add routes and sidebar icons for Dashboard and Import
    status: pending
  - id: build-test
    content: Build UI, rebuild Go binary, test both features end-to-end
    status: pending
---

# PocketBase V3: Dashboard + Data Import

## Feature 1: Real-time Metrics Dashboard

### Backend API (`apis/metrics.go`)

New endpoints under `/api/metrics` (superuser only):

| Endpoint | Returns |

|----------|---------|

| `GET /api/metrics/overview` | Request count, avg latency, error rate, DB size |

| `GET /api/metrics/requests` | Time-series of requests per minute (last 24h) |

| `GET /api/metrics/latency` | Avg/p50/p95/p99 latency over time |

| `GET /api/metrics/errors` | Error count by status code over time |

| `GET /api/metrics/endpoints` | Top 10 endpoints by request count |

| `GET /api/metrics/collections` | Record counts per collection |

**Implementation approach:**

- Query existing `_logs` table with aggregations on `data.execTime`, `data.status`, `data.url`
- Use SQL percentile calculations: `GROUP BY strftime('%Y-%m-%d %H:%M', created)`
- DB size via SQLite pragma: `PRAGMA page_count` * `PRAGMA page_size`
```go
// apis/metrics.go - key structures
type MetricsOverview struct {
    TotalRequests   int64   `json:"totalRequests"`
    AvgLatency      float64 `json:"avgLatency"`
    ErrorRate       float64 `json:"errorRate"`
    DatabaseSize    int64   `json:"databaseSize"`
    TotalCollections int    `json:"totalCollections"`
    TotalRecords    int64   `json:"totalRecords"`
}
```


### Frontend UI (`ui/src/pages/Dashboard.svelte`)

**Layout:**

```
+--------------------------------------------------+
| [Overview Cards - 4 metrics]                      |
| Requests/24h | Avg Latency | Error Rate | DB Size |
+--------------------------------------------------+
| [Requests Over Time]     | [Latency Distribution] |
| Line chart (24h)         | Line chart (p50/p95)   |
+--------------------------------------------------+
| [Top Endpoints]          | [Collections]          |
| Bar chart                | Table with counts      |
+--------------------------------------------------+
```

**Components:**

- `ui/src/pages/Dashboard.svelte` - Main page
- `ui/src/components/dashboard/MetricCard.svelte` - Stat card with icon
- `ui/src/components/dashboard/RequestsChart.svelte` - Time-series line chart
- `ui/src/components/dashboard/LatencyChart.svelte` - Multi-line percentile chart
- `ui/src/components/dashboard/EndpointsChart.svelte` - Horizontal bar chart
- `ui/src/components/dashboard/CollectionsTable.svelte` - Table with record counts
- `ui/src/stores/dashboard.js` - State management

**Auto-refresh:** Poll every 30 seconds (configurable)

---

## Feature 2: Data Import Wizard

### Backend API (`apis/import.go`)

New endpoints under `/api/import`:

| Endpoint | Purpose |

|----------|---------|

| `POST /api/import/preview` | Parse file, return headers + sample rows |

| `POST /api/import/validate` | Validate mapping against collection schema |

| `POST /api/import/execute` | Perform bulk import with progress |

```go
// apis/import.go - key structures
type ImportPreviewRequest struct {
    File       *multipart.FileHeader `form:"file"`
    Collection string                `form:"collection"`
}

type ImportPreviewResponse struct {
    Headers    []string         `json:"headers"`
    SampleRows [][]string       `json:"sampleRows"`
    TotalRows  int              `json:"totalRows"`
    Fields     []FieldInfo      `json:"fields"` // collection fields
}

type ImportExecuteRequest struct {
    Collection string            `json:"collection"`
    FileData   string            `json:"fileData"` // base64 or temp file ref
    Mapping    map[string]string `json:"mapping"`  // csv_header -> field_name
    SkipHeader bool              `json:"skipHeader"`
}
```

### Frontend UI (`ui/src/pages/ImportWizard.svelte`)

**4-Step Wizard:**

```
Step 1: Select Collection & Upload File
+----------------------------------+
| [Collection Dropdown]            |
| [File Drop Zone]                 |
| Supports: CSV, JSON              |
+----------------------------------+

Step 2: Preview Data
+----------------------------------+
| File: data.csv (1,234 rows)      |
| [Preview Table - first 5 rows]   |
+----------------------------------+

Step 3: Map Fields
+----------------------------------+
| CSV Column    →  Collection Field|
| --------------|------------------|
| [name]        →  [name ▼]        |
| [email]       →  [email ▼]       |
| [age]         →  [-- skip -- ▼]  |
| Auto-detect button               |
+----------------------------------+

Step 4: Import
+----------------------------------+
| [Progress Bar: 45%]              |
| Imported: 556 / 1,234            |
| Errors: 3                        |
| [View Error Log]                 |
+----------------------------------+
```

**Components:**

- `ui/src/pages/ImportWizard.svelte` - Main wizard page
- `ui/src/components/import/FileUpload.svelte` - Drag-drop file zone
- `ui/src/components/import/DataPreview.svelte` - Preview table
- `ui/src/components/import/FieldMapper.svelte` - Mapping UI
- `ui/src/components/import/ImportProgress.svelte` - Progress + errors
- `ui/src/stores/import.js` - State management

---

## File Structure

```
apis/
├── metrics.go          # Dashboard API endpoints
├── import.go           # Import wizard API endpoints

ui/src/
├── pages/
│   ├── Dashboard.svelte
│   └── ImportWizard.svelte
├── components/
│   ├── dashboard/
│   │   ├── MetricCard.svelte
│   │   ├── RequestsChart.svelte
│   │   ├── LatencyChart.svelte
│   │   ├── EndpointsChart.svelte
│   │   └── CollectionsTable.svelte
│   └── import/
│       ├── FileUpload.svelte
│       ├── DataPreview.svelte
│       ├── FieldMapper.svelte
│       └── ImportProgress.svelte
├── stores/
│   ├── dashboard.js
│   └── import.js
└── scss/
    ├── _dashboard.scss
    └── _import.scss
```

---

## Integration Points

- **Sidebar:** Add Dashboard icon (chart) and Import icon (upload)
- **Routes:** `/dashboard` and `/import`
- **Base API:** Register in `apis/base.go`

---

## Estimated Time

| Task | Hours |

|------|-------|

| Documentation Updates (PRD, TaskList, Memory Bank) | 1-2h |

| Dashboard Backend | 3-4h |

| Dashboard Frontend | 4-5h |

| Import Backend | 3-4h |

| Import Frontend | 4-5h |

| Testing + Polish | 2-3h |

| **Total** | **17-23h** |

---

## Phase 0: Documentation Updates (Do First)

### 1. Update PRD (`PocketBase_AI_Query_Assistant_PRD.md`)

Add V3 section with:

- Version 3.0 in version history table
- V3 User Stories (Dashboard monitoring, bulk data import)
- V3 Features: Real-time Metrics Dashboard, Data Import Wizard
- V3 Success Criteria
- V3 Timeline (PRs #18-21)

### 2. Update Task List (`PocketBase_AI_Query_TaskList.md`)

Add V3 PRs:

- PR #18: Metrics Backend API
- PR #19: Dashboard UI
- PR #20: Import Backend API  
- PR #21: Import Wizard UI

### 3. Update Memory Bank

Files to update:

- `memory-bank/activeContext.md` - Current focus: V3 Dashboard + Import
- `memory-bank/progress.md` - Add V3 status section
- `memory-bank/projectbrief.md` - Add V3 goals
- `memory-bank/productContext.md` - Add Dashboard/Import user stories