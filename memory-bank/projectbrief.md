# Project Brief: PocketBase AI Query Assistant

## Project Overview

**Project Name:** PocketBase AI Query Assistant  
**Type:** Brownfield Development (Fork & Extend)  
**Status:** Planning → Implementation  
**Version:** 1.0 MVP

## Core Mission

Build an **AI-powered natural language query assistant** as a feature extension to PocketBase. This feature enables users to query their database collections using plain English instead of learning PocketBase's filter syntax, making data exploration accessible to non-technical users while providing power users with a faster way to construct complex queries.

## Project Context

- **Base Repository:** Fork of https://github.com/pocketbase/pocketbase
- **Architecture:** Go backend + Svelte Admin UI
- **Integration Approach:** Must integrate cleanly with PocketBase's existing architecture, hook system, and Admin dashboard
- **Development Model:** 9 Pull Requests organized by feature area

## Primary Goals

1. **Accessibility:** Enable non-technical users to query databases without learning filter syntax
2. **Productivity:** Provide power users with faster query construction for complex filters
3. **Learning:** Show generated filter expressions so users can learn syntax over time
4. **Security:** Respect existing collection API rules and authentication
5. **Flexibility:** Support multiple LLM providers (OpenAI, Ollama, Anthropic, Custom)

## Success Criteria

The feature is complete when:
- ✅ User can type natural language query in Admin UI and see matching records
- ✅ Generated PocketBase filter expression is displayed with copy button
- ✅ LLM provider settings can be configured in Admin UI without code changes
- ✅ API endpoint `/api/ai/query` works for authenticated users
- ✅ Feature works with OpenAI API (`gpt-4o-mini` as default model)
- ✅ Optional: Ollama support for local/private deployments
- ✅ Invalid queries return helpful error messages, not crashes
- ✅ Security: API respects collection rules, doesn't expose unauthorized data
- ✅ Documentation complete: README, architecture overview, setup instructions
- ✅ Demo video recorded showing feature end-to-end

## Scope Boundaries

### In Scope (V1 MVP)
- Natural language query interface in Admin UI
- Schema-aware prompting (automatic field name understanding)
- Filter expression display with copy functionality
- API endpoint for programmatic access
- LLM configuration settings panel
- OpenAI integration (default)
- Basic validation and error handling

### Out of Scope (V1)
- Multi-collection joins
- Query history/favorites
- Natural language CREATE/UPDATE/DELETE
- Conversation memory
- Embedding/vector search
- Usage billing/quotas
- Fine-tuned models
- Streaming responses
- Per-collection enable/disable

## Key Constraints

1. **Must maintain compatibility** with existing PocketBase architecture
2. **Must respect security model** - collection API rules enforced
3. **Must be configurable** - admin can enable/disable and choose LLM provider
4. **Must be cost-conscious** - default to `gpt-4o-mini` for affordability
5. **Must handle errors gracefully** - no crashes, helpful error messages

## Timeline

**Estimated Total Time:** 35-45 hours  
**Approach:** 9 sequential PRs over 7-9 days  
**Current Phase:** Planning → Setup (PR #1)

## Key Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| LLM Provider | OpenAI with `gpt-4o-mini` | Best balance of speed, quality, and cost |
| UI Placement | Sidebar panel | Dedicated space, doesn't clutter main view |
| Settings Storage | Existing `_params` table | No schema changes needed |
| API Key | Encrypted at rest | Follow existing OAuth secret pattern |
| LLM Library | Raw `net/http` calls | Fewer dependencies, full control |
| Filter Execution | Both (filter + results) | Use `execute` parameter for flexibility |

## Open Questions

1. **Rate limiting?** Should we add basic rate limiting (e.g., 10 queries/minute)?
2. **Telemetry?** Log AI queries to PocketBase logs for debugging? Privacy implications?
3. **Fallback behavior?** If OpenAI API is unavailable, show error or hide AI feature entirely?
4. **Cost warning?** Show estimated cost per query in UI? (e.g., "~$0.0002 per query")



