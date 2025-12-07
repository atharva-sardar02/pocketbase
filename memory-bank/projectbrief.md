# Project Brief: PocketBase AI Query Assistant

## Project Overview

**Project Name:** PocketBase AI Query Assistant  
**Type:** Brownfield Development (Fork & Extend)  
**Status:** V1 ✅ Complete, V2 ✅ Complete, V3 🚧 In Progress  
**Current Version:** 3.0

## Core Mission

Build an **AI-powered natural language query assistant** and **developer tools** as feature extensions to PocketBase. These features enable users to query their database collections using plain English, execute SQL commands, monitor system health, and import data efficiently.

## Project Context

- **Base Repository:** Fork of https://github.com/pocketbase/pocketbase
- **Architecture:** Go backend + Svelte Admin UI
- **Integration Approach:** Must integrate cleanly with PocketBase's existing architecture, hook system, and Admin dashboard
- **Development Model:** 22 Pull Requests organized by feature area across 3 versions

## Primary Goals

### V1 Goals (✅ Complete)
1. **Accessibility:** Enable non-technical users to query databases without learning filter syntax
2. **Productivity:** Provide power users with faster query construction for complex filters
3. **Learning:** Show generated filter expressions so users can learn syntax over time
4. **Security:** Respect existing collection API rules and authentication
5. **Flexibility:** Support multiple LLM providers (OpenAI, Ollama, Anthropic, Custom)

### V2 Goals (✅ Complete)
1. **SQL Power:** Full SQL Terminal for developers who prefer SQL
2. **Multi-Table:** Query across related collections with JOINs
3. **Dual Output:** Show both PocketBase filter and SQL for flexibility
4. **Real Records:** SQL operations create actual PocketBase records

### V3 Goals (🚧 In Progress)
1. **Monitoring:** Real-time metrics dashboard for system health visibility
2. **Bulk Operations:** Data import wizard for CSV/JSON bulk loading
3. **Visual Insights:** Charts and graphs beyond text-based logs
4. **Productivity:** Reduce manual data entry with bulk import

## Success Criteria

### V1 Success Criteria (✅ Complete)
- ✅ User can type natural language query in Admin UI and see matching records
- ✅ Generated PocketBase filter expression is displayed with copy button
- ✅ LLM provider settings can be configured in Admin UI without code changes
- ✅ API endpoint `/api/ai/query` works for authenticated users
- ✅ Feature works with OpenAI API (`gpt-4o-mini` as default model)
- ✅ Ollama support for local/private deployments
- ✅ Invalid queries return helpful error messages
- ✅ Security: API respects collection rules

### V2 Success Criteria (✅ Complete)
- ✅ SQL Terminal page accessible from sidebar
- ✅ AI Mode generates SQL from natural language
- ✅ SQL Mode executes raw SQL commands
- ✅ CREATE TABLE creates real PocketBase collections
- ✅ INSERT/UPDATE/DELETE operations affect real records
- ✅ Multi-statement SQL execution supported
- ✅ Export results to CSV/JSON

### V3 Success Criteria (🚧 In Progress)
- [ ] Dashboard shows real-time requests, latency, error rate
- [ ] Charts display p50/p95/p99 latency percentiles
- [ ] Top endpoints visible in bar chart
- [ ] Collection record counts displayed
- [ ] CSV file import with field mapping
- [ ] JSON file import with field mapping
- [ ] Import progress and error logging

## Scope Boundaries

### In Scope (V1-V3)
| Version | Features |
|---------|----------|
| V1 | Natural language queries, filter display, LLM settings |
| V2 | SQL Terminal, multi-table queries, dual output, bulk SQL |
| V3 | Metrics dashboard, data import wizard |

### Out of Scope (Future)
- Conversation memory / chat history
- Embedding/vector search
- Usage billing/quotas
- Fine-tuned models
- Streaming responses
- Per-collection enable/disable

## Key Constraints

1. **Must maintain compatibility** with existing PocketBase architecture
2. **Must respect security model** - collection API rules enforced
3. **Must be configurable** - admin can enable/disable features
4. **Must be cost-conscious** - default to `gpt-4o-mini` for affordability
5. **Must handle errors gracefully** - no crashes, helpful error messages

## Timeline

| Version | Estimated | Actual | Status |
|---------|-----------|--------|--------|
| V1 | 35-45h | 38h | ✅ Complete |
| V2 | 41-50h | 45h | ✅ Complete |
| V3 | 17-23h | -- | 🚧 In Progress |
| **Total** | **93-118h** | **83h+** | |

## Key Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| LLM Provider | OpenAI with `gpt-4o-mini` | Best balance of speed, quality, and cost |
| UI Placement | Sidebar panels | Dedicated space, doesn't clutter main view |
| Settings Storage | Existing `_params` table | No schema changes needed |
| API Key | Encrypted at rest | Follow existing OAuth secret pattern |
| LLM Library | Raw `net/http` calls | Fewer dependencies, full control |
| SQL Execution | Via PocketBase APIs | Creates real records, respects rules |
| Dashboard Data | From `_logs` table | No new tables, uses existing data |
| Import Format | CSV + JSON | Most common data formats |

## Repository

**GitHub:** https://github.com/atharva-sardar02/pocketbase.git  
**Branch Strategy:**
- `master` - stable releases
- `feat/v3-dashboard-import` - V3 development
