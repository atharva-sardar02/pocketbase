# Product Requirements Document
## PocketBase AI Query Assistant

**Version:** 2.0  
**Date:** December 2025  
**Status:** V1 Complete, V2 In Planning

---

## Version History

| Version | Date | Status | Description |
|---------|------|--------|-------------|
| 1.0 | Dec 2025 | ✅ Complete | Single-collection AI Query with filter generation |
| 2.0 | Dec 2025 | 🚧 Planning | Multi-table queries, SQL Terminal, dual output |

---

## 1. Executive Summary

This document outlines the requirements for building an **AI-powered natural language query assistant** as a feature extension to PocketBase. The feature enables users to query their database collections using plain English instead of learning PocketBase's filter syntax, making data exploration accessible to non-technical users while providing power users with a faster way to construct complex queries.

### Project Context

This is a **brownfield development project**. We will fork the PocketBase open-source repository (Go backend + Svelte Admin UI) and extend it with this new capability. The feature must integrate cleanly with PocketBase's existing architecture, hook system, and Admin dashboard.

**Repository:** https://github.com/pocketbase/pocketbase

---

## 2. User Stories

### 2.1 Primary User: Non-Technical Admin/Business User

**Profile:** Uses PocketBase Admin UI to manage application data. Not familiar with filter syntax or SQL. Needs to find and analyze data quickly.

| # | User Story |
|---|------------|
| 1 | As a business user, I want to type "show me all orders from last week that are still pending" and get the right records, so I don't have to learn the filter syntax. |
| 2 | As a business user, I want to see the generated filter expression alongside results, so I can learn the syntax over time if I choose. |
| 3 | As a business user, I want the AI to understand my collection's field names automatically, so I don't have to specify exact column names. |
| 4 | As a business user, I want helpful error messages when my query doesn't make sense, so I can rephrase and try again. |

### 2.2 Secondary User: Developer/Power User

**Profile:** Building applications on PocketBase. Knows filter syntax but wants faster query construction for complex filters.

| # | User Story |
|---|------------|
| 5 | As a developer, I want to describe a complex filter in plain English and get the equivalent PocketBase filter expression, so I can copy-paste it into my code. |
| 6 | As a developer, I want API access to the natural language query feature, so I can build AI-powered search into my application's frontend. |
| 7 | As a developer, I want the feature to respect existing collection API rules, so my security model isn't bypassed. |

### 2.3 Tertiary User: PocketBase Administrator

**Profile:** Manages the PocketBase instance. Responsible for configuration and security.

| # | User Story |
|---|------------|
| 8 | As an admin, I want to configure which LLM provider to use (OpenAI, local Ollama, etc.), so I can control costs and data privacy. |
| 9 | As an admin, I want to enable/disable the AI query feature globally or per-collection, so I can control where it's available. |
| 10 | As an admin, I want to see usage logs for AI queries, so I can monitor costs and detect abuse. |

### 2.4 V2 User Stories — Multi-Table & SQL Terminal

| # | User Story |
|---|------------|
| 11 | As a user, I want to query across multiple related tables (e.g., "orders with customer names"), so I can see combined data without manual joins. |
| 12 | As a user, I want to see both the PocketBase filter AND the SQL equivalent, so I can choose which to use. |
| 13 | As a user, I want to edit the generated query before executing it, so I can refine the results. |
| 14 | As a developer, I want a SQL terminal to run raw SQL commands, so I can manage my database efficiently. |
| 15 | As a developer, I want to create collections using SQL syntax (CREATE TABLE), so I can use familiar database commands. |
| 16 | As a developer, I want INSERT/UPDATE/DELETE operations in SQL Terminal to create real PocketBase records, so changes appear in Admin UI. |
| 17 | As a user, I want AI assistance in SQL Terminal, so I can describe what I want in plain English and get SQL generated. |

---

## 3. Key Features — Version 1.0 (MVP)

### 3.1 Natural Language Query Interface (Admin UI)

| Feature | Description |
|---------|-------------|
| **Chat-style input** | Text input box in the Admin UI's collection records view for typing queries in plain English |
| **Schema-aware prompting** | System automatically injects current collection's schema (field names, types, relations) into LLM context |
| **Filter expression display** | Shows generated PocketBase filter syntax alongside results with a copy button |
| **Query refinement** | Users can modify natural language query and re-run, or edit the generated filter directly |

**UI Mockup Concept (Sidebar):**
```
┌──────────────────────────────────────────────────────────────────────────┐
│  PocketBase Admin                                                        │
├────────────────────┬─────────────────────────────────────────────────────┤
│                    │  Collection: orders                          [⚙️]  │
│  📁 Collections    ├─────────────────────────────────────────────────────┤
│    └─ orders  ←    │                                                     │
│    └─ users        │  Filter: [status = "pending" && total > 100    ]   │
│    └─ products     │                                                     │
│                    │  ┌─────────────────────────────────────────────┐   │
│  ─────────────     │  │ ID: abc123 | Total: $150 | Status: pending  │   │
│                    │  │ ID: def456 | Total: $200 | Status: pending  │   │
│  🤖 AI Query  ←    │  │ ID: ghi789 | Total: $175 | Status: pending  │   │
│  ─────────────     │  └─────────────────────────────────────────────┘   │
│                    │                                                     │
│  ⚙️ Settings       │                                                     │
│                    │                                                     │
└────────────────────┴─────────────────────────────────────────────────────┘

AI Query Sidebar Panel (when clicked):
┌────────────────────────────────────────┐
│  🤖 AI Query Assistant                 │
├────────────────────────────────────────┤
│  Collection: [orders          ▼]       │
│                                        │
│  Ask in plain English:                 │
│  ┌──────────────────────────────────┐  │
│  │ show me pending orders over $100 │  │
│  │ from this week                   │  │
│  └──────────────────────────────────┘  │
│                         [🔍 Search]    │
│                                        │
│  ─────────────────────────────────     │
│  Generated Filter:                     │
│  ┌──────────────────────────────────┐  │
│  │ status="pending" && total>100 && │  │
│  │ created>=@now-604800             │  │
│  └──────────────────────────────────┘  │
│            [📋 Copy] [▶ Apply Filter]  │
│                                        │
│  Results: 3 records found              │
│  → View in collection                  │
└────────────────────────────────────────┘
```

### 3.2 API Endpoint

| Endpoint | Description |
|----------|-------------|
| `POST /api/ai/query` | Accepts natural language query and collection name, returns filter expression and/or results |

**Request Schema:**
```json
{
  "collection": "orders",
  "query": "pending orders over $100 from last week",
  "execute": true,
  "page": 1,
  "perPage": 20
}
```

**Response Schema:**
```json
{
  "filter": "status='pending' && total>100 && created>=@now-604800",
  "results": [...],
  "totalItems": 42,
  "page": 1,
  "perPage": 20
}
```

**Security:**
- Authentication required (uses existing PocketBase auth)
- Respects collection API rules — only returns data the authenticated user can access
- Superuser-only option available via settings

### 3.3 LLM Configuration (Settings Panel)

New settings section in Admin UI: **Settings → AI Query**

| Setting | Options | Default |
|---------|---------|---------|
| **Enable AI Query** | On/Off | Off |
| **LLM Provider** | OpenAI, Ollama, Anthropic, Custom | OpenAI |
| **API Base URL** | Text input | `https://api.openai.com/v1` |
| **API Key** | Password input (encrypted) | Empty (required) |
| **Model** | Dropdown | `gpt-4o-mini` |
| **Temperature** | Slider 0-1 | 0.1 |
| **Test Connection** | Button | — |

**Recommended OpenAI Models:**
| Model | Speed | Quality | Cost | Best For |
|-------|-------|---------|------|----------|
| `gpt-4o-mini` | Fast | Good | ~$0.00015/query | **Default — best balance** |
| `gpt-4o` | Medium | Excellent | ~$0.005/query | Complex queries |
| `gpt-3.5-turbo` | Fastest | Decent | ~$0.0001/query | High volume, simple queries |

### 3.4 Query Translation Engine

**Core Components:**

1. **System Prompt Template**
   - Teaches LLM the PocketBase filter syntax
   - Includes all operators, datetime macros, relation syntax
   - Provides few-shot examples

2. **Schema Injection**
   - Dynamically builds context from collection schema
   - Includes: field names, field types, relation targets
   - Example: `Fields: id (text), title (text), status (select: draft|published|archived), author (relation→users), created (datetime)`

3. **Validation Layer**
   - Parses generated filter before execution
   - Validates field names exist in collection
   - Validates operators are appropriate for field types
   - Rejects obviously malformed syntax

4. **Error Handling**
   - LLM timeout → "Query is taking too long, please try again"
   - Invalid filter → "I couldn't understand that query. Try rephrasing with specific field names."
   - No results → "No records match your query" (not an error)

---

## 4. Tech Stack

| Layer | Technology | Notes |
|-------|------------|-------|
| **Backend** | Go 1.21+ | Existing PocketBase language. New API routes in `/apis` |
| **Frontend** | Svelte 4 | Existing Admin UI framework. Components in `/ui/src` |
| **Database** | SQLite | Existing. AI settings stored in `_params` table |
| **LLM Communication** | HTTP/REST | OpenAI-compatible API format (works with Ollama, OpenAI, etc.) |
| **Build - Backend** | `go build` | CGO_ENABLED=0 for static binary |
| **Build - Frontend** | `npm run build` | Vite-based, outputs to `/ui/dist` |

### Key Files to Modify/Create

**V1 Files (✅ Complete):**
```
pocketbase/
├── apis/
│   └── ai_query.go              # API endpoint handler
├── core/
│   └── ai_settings.go           # AI settings struct
├── services/ai/
│   ├── openai_client.go         # LLM API client
│   ├── schema_extractor.go      # Collection schema extraction
│   ├── prompt_builder.go        # Prompt construction
│   └── filter_validator.go      # Filter validation
├── ui/src/
│   ├── components/ai/
│   │   ├── AIQueryPanel.svelte  # Main query panel
│   │   ├── AIQueryInput.svelte  # Query input
│   │   ├── AIFilterDisplay.svelte # Filter display
│   │   ├── AIQueryResults.svelte  # Results display
│   │   └── AISettingsForm.svelte  # Settings form
│   └── pages/settings/
│       └── AI.svelte            # Settings page
```

**V2 Files (🚧 Planned):**
```
pocketbase/
├── apis/
│   └── sql_terminal.go          # NEW: SQL Terminal API endpoints
├── services/sql/
│   ├── parser.go                # NEW: SQL statement parser
│   ├── executor.go              # NEW: SQL execution via PocketBase APIs
│   └── mapper.go                # NEW: SQL type → PocketBase field mapper
├── services/ai/
│   └── schema_extractor.go      # MODIFY: Extract ALL collections
├── ui/src/
│   ├── components/ai/
│   │   ├── QueryTabs.svelte     # NEW: Filter/SQL tab switcher
│   │   └── EditableCodeBlock.svelte # NEW: Editable query component
│   ├── components/sql/
│   │   ├── SQLEditor.svelte     # NEW: Code editor with syntax highlight
│   │   ├── SchemaExplorer.svelte # NEW: Collections sidebar browser
│   │   ├── ResultsTable.svelte  # NEW: Dynamic results display
│   │   └── QueryHistory.svelte  # NEW: Command history
│   ├── pages/
│   │   └── SQLTerminal.svelte   # NEW: Main SQL Terminal page
│   └── stores/
│       └── sql.js               # NEW: SQL terminal state
```

---

## 5. Out of Scope (V1) — Now in V2

The following features were excluded from V1 MVP but are now being added in V2:

| Feature | V1 Status | V2 Status |
|---------|-----------|-----------|
| Multi-collection joins | ❌ Excluded | ✅ **Now in V2** |
| Natural language CREATE/UPDATE/DELETE | ❌ Excluded | ✅ **Now in V2 (SQL Terminal)** |
| Query history/favorites | ❌ Excluded | 🔄 Future |
| Conversation memory | ❌ Excluded | 🔄 Future |
| Embedding/vector search | ❌ Excluded | 🔄 Future |
| Usage billing/quotas | ❌ Excluded | 🔄 Future |
| Streaming responses | ❌ Excluded | 🔄 Future |
| Per-collection enable/disable | ❌ Excluded | 🔄 Future |

---

## 6. Version 2.0 Features

### 6.1 Enhanced AI Query — Multi-Table Support

**What's New:**
- Query across **multiple collections** with JOINs
- Support for **aggregates** (COUNT, SUM, AVG, etc.)
- **Dual output** — Both Filter and SQL shown in tabs
- **Editable queries** — Modify generated query before executing
- **Complex queries** — GROUP BY, HAVING, subqueries

**Dual Output UI:**
```
┌─────────────────────────────────────────────────┐
│  AI Query Results                               │
├─────────────────────────────────────────────────┤
│  ┌────────┐ ┌─────┐                             │
│  │ Filter │ │ SQL │  ← Click to switch tabs     │
│  └────────┘ └─────┘                             │
├─────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────┐    │
│  │ status = "active" && total > 100        │    │  ← Editable!
│  └─────────────────────────────────────────┘    │
│  [Copy] [Execute] [Apply to Collection]         │
├─────────────────────────────────────────────────┤
│  Results: 42 records                            │
└─────────────────────────────────────────────────┘
```

**Example Multi-Table Queries:**

| Natural Language | Generated SQL |
|------------------|---------------|
| "orders with customer names" | `SELECT o.*, c.name FROM orders o JOIN customers c ON o.customer = c.id` |
| "total revenue by category" | `SELECT category, SUM(total) FROM orders GROUP BY category` |
| "customers who spent over $1000" | `SELECT c.*, SUM(o.total) as spent FROM customers c JOIN orders o ON c.id = o.customer GROUP BY c.id HAVING spent > 1000` |

### 6.2 SQL Terminal — Full Database Console

A complete database management tool built into PocketBase Admin UI.

**Key Features:**
- **Two Modes:**
  - **AI Mode** — Type natural language, AI generates and executes SQL
  - **SQL Mode** — Type raw SQL directly
- **Full Database Access:**
  - SELECT, INSERT, UPDATE, DELETE
  - CREATE TABLE, ALTER TABLE, DROP TABLE
  - All operations create **real PocketBase collections** (not raw SQLite tables)
- **Developer Tools:**
  - Syntax highlighting
  - Auto-complete for table/column names
  - Command history (up/down arrows)
  - Export results to CSV/JSON
  - Schema browser sidebar

**UI Mockup:**
```
┌─────────────────────────────────────────────────────────────────┐
│  🖥️ SQL Terminal                           [AI Mode ●] [SQL Mode]│
├─────────────────────────────────────────────────────────────────┤
│  Schema Browser:              │  Query Editor:                  │
│  ┌─────────────────────┐      │  ┌─────────────────────────────┐│
│  │ 📁 Collections      │      │  │ CREATE TABLE products (     ││
│  │   ├─ orders        │      │  │   name TEXT NOT NULL,       ││
│  │   ├─ customers     │      │  │   price NUMBER,             ││
│  │   └─ products ←NEW │      │  │   category TEXT             ││
│  └─────────────────────┘      │  │ );                          ││
│                               │  └─────────────────────────────┘│
│                               │  [▶ Run] [Clear] [History ▼]    │
├───────────────────────────────┴─────────────────────────────────┤
│  Output:                                        [Export CSV]    │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ ✅ Collection 'products' created successfully               ││
│  │    Fields: name (text), price (number), category (text)     ││
│  │    → View in Collections                                    ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

**SQL → PocketBase Collection Mapping:**

| SQL Command | PocketBase Action |
|-------------|-------------------|
| `CREATE TABLE products (name TEXT, price REAL)` | Creates PocketBase collection with fields |
| `ALTER TABLE products ADD COLUMN stock INT` | Adds field to collection schema |
| `DROP TABLE products` | Deletes collection |
| `INSERT INTO products VALUES (...)` | Creates record via PocketBase API |
| `UPDATE products SET price = 19.99 WHERE id = 'x'` | Updates record via PocketBase API |
| `DELETE FROM products WHERE stock = 0` | Deletes records via PocketBase API |

**SQL Type → PocketBase Field Type Mapping:**

| SQL Type | PocketBase Field |
|----------|------------------|
| `TEXT`, `VARCHAR` | text |
| `INTEGER`, `INT` | number |
| `REAL`, `FLOAT` | number |
| `BOOLEAN` | bool |
| `DATE`, `DATETIME` | date |
| `TEXT REFERENCES table(id)` | relation |
| `TEXT CHECK(value IN (...))` | select |

### 6.3 New API Endpoints (V2)

**Enhanced AI Query:**

`POST /api/ai/query` (updated response)
```json
{
  "query": "orders with customer names where total > 100",
  "execute": true
}

// Response now includes both filter AND SQL:
{
  "filter": "total > 100",
  "filterCollection": "orders",
  "sql": "SELECT o.*, c.name FROM orders o JOIN customers c ON o.customer = c.id WHERE o.total > 100",
  "canUseFilter": true,
  "results": [...],
  "columns": ["id", "customer", "total", "customer_name"]
}
```

**SQL Terminal Execute:**

`POST /api/sql/execute`
```json
{
  "sql": "CREATE TABLE products (name TEXT NOT NULL, price REAL)"
}

// Response:
{
  "type": "create",
  "collection": "products",
  "fields": [
    {"name": "name", "type": "text", "required": true},
    {"name": "price", "type": "number"}
  ],
  "message": "Collection 'products' created successfully"
}
```

**SQL Terminal AI Mode:**

`POST /api/sql/ai`
```json
{
  "query": "Create a products table with name, price, and category"
}

// Response:
{
  "sql": "CREATE TABLE products (name TEXT NOT NULL, price REAL, category TEXT)",
  "explanation": "Creates a 'products' collection with three fields"
}
```

### 6.4 Access Control (V2)

| Feature | Access Level |
|---------|--------------|
| AI Query (read) | Any authenticated user |
| SQL Terminal (read) | Any authenticated user |
| SQL Terminal (write) | Any authenticated user |
| SQL Terminal Settings | Superuser only |

### 6.5 Security Considerations (V2)

| Risk | Mitigation |
|------|------------|
| SQL Injection | Validate table names exist, parameterize where possible |
| Accidental Data Loss | Confirmation dialog for DROP/DELETE operations |
| Resource Exhaustion | Query timeout (30s), result limit (10,000 rows) |
| Unauthorized Schema Changes | Optional superuser-only mode in settings |

---

## 6. Technical Pitfalls & Considerations

### 6.1 Go Backend Challenges

| Risk | Mitigation |
|------|------------|
| **Learning curve** | Go's strict typing and error handling require adjustment. Budget 1-2 days for familiarization. Use AI assistance heavily. |
| **PocketBase hook system** | Must understand `OnServe()` and router binding. Study `/apis/record.go` as reference. |
| **HTTP client patterns** | Go requires explicit timeouts. Use `context.WithTimeout()` for LLM calls. Default 30s timeout. |
| **Dependency management** | `go mod tidy` after adding imports. Avoid external LLM libraries — use raw `net/http`. |

### 6.2 Svelte Frontend Challenges

| Risk | Mitigation |
|------|------------|
| **Component architecture** | PocketBase uses specific patterns for forms, buttons, modals. Copy existing component structures. |
| **Build process** | Must run `npm run build` in `/ui` after every change, then rebuild Go binary. Create script to automate. |
| **State management** | Uses Svelte stores. Study `/ui/src/stores/` before creating new state. |
| **API client** | Use existing `ApiClient` in `/ui/src/utils/ApiClient.js` for consistency. |

### 6.3 LLM Integration Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| **Prompt injection** | HIGH | User could craft query to manipulate LLM. Sanitize inputs, validate outputs against schema. Never execute raw LLM output as code. |
| **Hallucination** | MEDIUM | LLM invents field names or syntax. Validate all field names exist in collection before executing filter. |
| **Latency** | LOW | OpenAI `gpt-4o-mini` typically responds in 0.5-2 seconds. Show loading spinner. |
| **Cost** | MEDIUM | Each query costs ~$0.00015 with `gpt-4o-mini`. Add API key requirement, consider rate limiting. |
| **API availability** | LOW | OpenAI has 99.9% uptime. Show clear error if API fails, allow retry. |

### 6.4 Architecture Decisions

> **These require your input before implementation:**

| Decision | Option A | Option B | **Decision** |
|----------|----------|----------|----------------|
| **Settings storage** | Existing `_params` table | New `ai_settings` collection | **Option A** — existing `_params` table ✓ |
| **API key encryption** | Plain text in DB | Encrypted at rest | **Option B** — follow existing OAuth secret pattern |
| **LLM library** | Raw `net/http` calls | Use `langchaingo` package | **Option A** — fewer dependencies, full control |
| **Filter execution** | Return filter only | Return filter + results | **Both** — use `execute` parameter |
| **UI placement** | New sidebar tab | Inline in records view | **Sidebar** — dedicated AI Query panel ✓ |
| **Default LLM** | Ollama (local) | OpenAI (cloud) | **OpenAI** — using `gpt-4o-mini` ✓ |

---

## 7. Success Criteria

### V1 Success Criteria (✅ Complete)

- [x] User can type natural language query in Admin UI and see matching records
- [x] Generated PocketBase filter expression is displayed with copy button
- [x] LLM provider settings can be configured in Admin UI without code changes
- [x] API endpoint `/api/ai/query` works for authenticated users
- [x] Feature works with OpenAI API (`gpt-4o-mini` as default model)
- [x] Optional: Ollama support for local/private deployments
- [x] Invalid queries return helpful error messages, not crashes
- [x] Security: API respects collection rules, doesn't expose unauthorized data
- [x] Documentation complete: README, architecture overview, setup instructions

### V2 Success Criteria (🚧 In Progress)

- [ ] Multi-table queries work with JOINs across related collections
- [ ] Dual output shows both Filter and SQL in switchable tabs
- [ ] Users can edit generated queries before executing
- [ ] SQL Terminal page accessible from sidebar
- [ ] SQL Terminal supports AI Mode (natural language → SQL)
- [ ] SQL Terminal supports SQL Mode (direct SQL execution)
- [ ] CREATE TABLE creates real PocketBase collections
- [ ] INSERT/UPDATE/DELETE operations create/modify real records
- [ ] Schema browser shows all collections in SQL Terminal
- [ ] Query history saved and accessible
- [ ] Export results to CSV/JSON
- [ ] Changes immediately visible in Admin UI

---

## 8. Proposed Timeline

### V1 Timeline (✅ Complete — 38 hours)

| Day | Focus | Status |
|-----|-------|--------|
| Day 1-2 | Setup & Settings | ✅ Complete |
| Day 3-4 | Backend (API, LLM, Validation) | ✅ Complete |
| Day 5-6 | Frontend (AI Query, Settings) | ✅ Complete |
| Day 7 | Testing & Documentation | ✅ Complete |

### V2 Timeline (Estimated 41-50 hours)

| Phase | Focus | Est. Hours |
|-------|-------|------------|
| **Phase 1** | Multi-Collection Schema (PR #10) | 4-5 |
| **Phase 2** | Dual Output Backend (PR #11) | 5-6 |
| **Phase 3** | Editable Query UI (PR #12) | 4-5 |
| **Phase 4** | SQL Parser & Mapper (PR #13) | 6-7 |
| **Phase 5** | SQL Executor (PR #14) | 6-7 |
| **Phase 6** | SQL Terminal API (PR #15) | 5-6 |
| **Phase 7** | SQL Terminal UI (PR #16) | 8-10 |
| **Phase 8** | Documentation (PR #17) | 3-4 |
| **Total** | | **41-50 hours** |

---

## 9. Appendix

### A. PocketBase Filter Syntax Reference

The AI must generate valid syntax using these patterns:

**Operators:**
| Operator | Meaning | Example |
|----------|---------|---------|
| `=` | Equals | `status = "active"` |
| `!=` | Not equals | `status != "deleted"` |
| `>` `<` `>=` `<=` | Comparison | `total > 100` |
| `~` | Contains (LIKE) | `title ~ "hello"` |
| `!~` | Not contains | `title !~ "spam"` |
| `?=` | Any equals (arrays) | `tags ?= "urgent"` |
| `?~` | Any contains (arrays) | `tags ?~ "imp"` |

**Logical Operators:**
- `&&` — AND
- `||` — OR
- `()` — Grouping

**Datetime Macros:**
| Macro | Meaning |
|-------|---------|
| `@now` | Current datetime |
| `@second` | Current second (0-59) |
| `@minute` | Current minute (0-59) |
| `@hour` | Current hour (0-23) |
| `@weekday` | Day of week (0-6, Sunday=0) |
| `@day` | Day of month (1-31) |
| `@month` | Month (1-12) |
| `@year` | Year (e.g., 2025) |

**Example Translations:**

| Natural Language | PocketBase Filter |
|-----------------|-------------------|
| "active users" | `status = "active"` |
| "orders over $100" | `total > 100` |
| "posts from last week" | `created >= @now - 604800` |
| "titles containing 'hello'" | `title ~ "hello"` |
| "my records" | `user = @request.auth.id` |
| "pending OR processing" | `status = "pending" \|\| status = "processing"` |
| "high priority items from today" | `priority = "high" && created >= @now - 86400` |

### B. Sample System Prompt (Draft)

```
You are a PocketBase filter query generator. Convert natural language queries into valid PocketBase filter syntax.

COLLECTION SCHEMA:
{schema_will_be_injected_here}

FILTER SYNTAX RULES:
- Use = for exact match: field = "value"
- Use ~ for contains: field ~ "partial"
- Use > < >= <= for numbers/dates
- Use && for AND, || for OR
- Use @now for current time, subtract seconds for past dates
- Wrap string values in double quotes
- Field names are case-sensitive

EXAMPLES:
User: "active users"
Filter: status = "active"

User: "orders over 100 dollars from this week"
Filter: total > 100 && created >= @now - 604800

User: "posts containing javascript in the title"
Filter: title ~ "javascript"

USER QUERY: {user_query}

Respond with ONLY the filter expression, no explanation.
```

---

## 10. Decisions Made

| Question | Decision |
|----------|----------|
| **LLM Provider** | OpenAI with `gpt-4o-mini` as default |
| **UI Placement** | Sidebar panel (dedicated AI Query section) |
| **Settings Storage** | Existing `_params` table |
| **API Key** | Required, encrypted at rest |

## 11. Remaining Open Questions

1. **Rate limiting?** Should we add basic rate limiting (e.g., 10 queries/minute) to prevent abuse/cost overrun?

2. **Telemetry?** Log AI queries to PocketBase logs for debugging? Privacy implications?

3. **Fallback behavior?** If OpenAI API is unavailable, show error or hide AI feature entirely?

4. **Cost warning?** Show estimated cost per query in UI? (e.g., "~$0.0002 per query")

---

**Document Status:** Ready for review  
**Next Steps:** Review and approve PRD → Begin Day 1 implementation
