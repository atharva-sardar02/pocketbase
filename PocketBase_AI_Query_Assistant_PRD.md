# Product Requirements Document
## PocketBase AI Query Assistant

**Version:** 1.0  
**Date:** December 2025  
**Status:** Draft for Review

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

```
pocketbase/
├── apis/
│   └── ai_query.go          # NEW: API endpoint handler
├── core/
│   └── settings_ai.go       # NEW: AI settings struct
├── ui/src/
│   ├── components/
│   │   └── ai/
│   │       ├── AIQueryInput.svelte    # NEW
│   │       ├── AIQueryResults.svelte  # NEW
│   │       └── AISettings.svelte      # NEW
│   └── pages/
│       └── collections/
│           └── Records.svelte         # MODIFY: Add AI query UI
└── examples/base/
    └── main.go              # Entry point for testing
```

---

## 5. Out of Scope (V1)

The following features are **explicitly excluded** from MVP to maintain scope:

| Feature | Reason | Future Version? |
|---------|--------|-----------------|
| Multi-collection joins | Complex `@collection` syntax requires advanced prompting | V2 |
| Query history/favorites | Requires new database tables, UI work | V2 |
| Natural language CREATE/UPDATE/DELETE | Security risk, scope creep | Maybe V3 |
| Conversation memory | "Show me more like last query" context | V2 |
| Embedding/vector search | Requires embedding model, index storage | V3 |
| Usage billing/quotas | Rate limiting, cost tracking | V2 |
| Fine-tuned models | Custom training on PocketBase syntax | Never (use prompt engineering) |
| Streaming responses | Real-time token display | V2 |
| Per-collection enable/disable | Granular control | V2 |

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

The feature is **complete** when:

- [ ] User can type natural language query in Admin UI and see matching records
- [ ] Generated PocketBase filter expression is displayed with copy button
- [ ] LLM provider settings can be configured in Admin UI without code changes
- [ ] API endpoint `/api/ai/query` works for authenticated users
- [ ] Feature works with OpenAI API (`gpt-4o-mini` as default model)
- [ ] Optional: Ollama support for local/private deployments
- [ ] Invalid queries return helpful error messages, not crashes
- [ ] Security: API respects collection rules, doesn't expose unauthorized data
- [ ] Documentation complete: README, architecture overview, setup instructions
- [ ] Demo video recorded showing feature end-to-end

---

## 8. Proposed Timeline

| Day | Focus | Deliverables |
|-----|-------|--------------|
| **Day 1** | Setup & Learning | Fork repo, dev environment working, Go/Svelte basics, codebase map |
| **Day 2** | Architecture | Detailed technical design, identify all files to modify, POC of LLM call |
| **Day 3** | Backend Core | API endpoint, LLM integration, prompt engineering, basic filter generation |
| **Day 4** | Backend Polish | Schema injection, validation layer, error handling, settings storage |
| **Day 5** | Frontend Core | AI query input component, results display, integration with records page |
| **Day 6** | Frontend Polish | Settings panel, loading states, copy button, error messages |
| **Day 7** | Finalize | Testing, bug fixes, documentation, demo video |

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
