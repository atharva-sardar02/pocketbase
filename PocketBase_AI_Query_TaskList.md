# PocketBase AI Query Assistant — Task List

## Project Overview

**Repository:** Fork of https://github.com/pocketbase/pocketbase  
**Feature:** AI-powered natural language query assistant  
**Total PRs:** 8  

---

## File Structure Overview

```
pocketbase/                          # Forked repository root
├── apis/
│   ├── ai_query.go                  # NEW — API endpoint handler
│   └── ai_query_test.go             # NEW — API endpoint tests
├── core/
│   ├── ai_settings.go               # NEW — AI settings struct & validation
│   ├── ai_settings_test.go          # NEW — Settings tests
│   └── settings.go                  # MODIFY — Add AI settings to main settings
├── services/
│   └── ai/
│       ├── openai_client.go         # NEW — OpenAI API client
│       ├── openai_client_test.go    # NEW — Client tests (mocked)
│       ├── prompt_builder.go        # NEW — System prompt construction
│       ├── prompt_builder_test.go   # NEW — Prompt tests
│       ├── schema_extractor.go      # NEW — Collection schema extraction
│       ├── schema_extractor_test.go # NEW — Schema extraction tests
│       ├── filter_validator.go      # NEW — Filter syntax validation
│       └── filter_validator_test.go # NEW — Validation tests
├── ui/
│   ├── src/
│   │   ├── components/
│   │   │   └── ai/
│   │   │       ├── AIQueryPanel.svelte      # NEW — Main sidebar panel
│   │   │       ├── AIQueryInput.svelte      # NEW — Query input component
│   │   │       ├── AIQueryResults.svelte    # NEW — Results display
│   │   │       ├── AIFilterDisplay.svelte   # NEW — Filter with copy button
│   │   │       └── AISettingsForm.svelte    # NEW — Settings form component
│   │   ├── pages/
│   │   │   └── settings/
│   │   │       └── AI.svelte                # NEW — AI settings page
│   │   ├── stores/
│   │   │   └── ai.js                        # NEW — AI-related state store
│   │   └── App.svelte                       # MODIFY — Add sidebar entry
│   └── package.json                         # MODIFY — Add any new dependencies
├── examples/base/
│   └── main.go                              # Entry point (no changes needed)
├── tests/
│   └── integration/
│       └── ai_query_integration_test.go     # NEW — E2E integration tests
└── docs/
    └── AI_QUERY_FEATURE.md                  # NEW — Feature documentation
```

---

## PR #1: Project Setup & Repository Configuration

**Branch:** `feat/ai-query-setup`  
**Estimated Time:** 2-3 hours  
**Dependencies:** None  

### Description
Fork the repository, set up development environment, and create the foundational file structure for the AI Query feature.

### Tasks

- [ ] **1.1** Fork PocketBase repository to personal GitHub
- [ ] **1.2** Clone forked repository locally
- [ ] **1.3** Verify Go environment (`go version` ≥ 1.21)
- [ ] **1.4** Verify Node.js environment (`node -v` ≥ 18)
- [ ] **1.5** Run initial build to confirm setup works
  ```powershell
  cd ui && npm install && npm run build
  cd ../examples/base
  $env:CGO_ENABLED="0"; go build
  .\base.exe serve
  ```
- [ ] **1.6** Create feature branch structure
- [ ] **1.7** Create empty directory structure for new files:
  - [ ] `services/ai/` directory
  - [ ] `ui/src/components/ai/` directory
  - [ ] `tests/integration/` directory
  - [ ] `docs/` directory
- [ ] **1.8** Add `.gitkeep` files to empty directories
- [ ] **1.9** Update main `README.md` with AI Query feature mention (placeholder)
- [ ] **1.10** Create `docs/AI_QUERY_FEATURE.md` with initial structure

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/ai/.gitkeep` | CREATE | Placeholder for AI services |
| `ui/src/components/ai/.gitkeep` | CREATE | Placeholder for UI components |
| `tests/integration/.gitkeep` | CREATE | Placeholder for integration tests |
| `docs/AI_QUERY_FEATURE.md` | CREATE | Feature documentation skeleton |
| `README.md` | MODIFY | Add feature mention |

### Tests
> ❌ **No tests required** — This PR is setup only, no functional code.

### Verification
- [ ] `go build` succeeds in `examples/base`
- [ ] `npm run build` succeeds in `ui`
- [ ] PocketBase starts and Admin UI loads at `http://127.0.0.1:8090/_/`

---

## PR #2: AI Settings Data Structure & Storage

**Branch:** `feat/ai-query-settings`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #1  

### Description
Create the Go data structures for AI settings and integrate them into PocketBase's existing settings system (`_params` table).

### Tasks

- [ ] **2.1** Create `core/ai_settings.go` with settings struct:
  ```go
  type AISettings struct {
      Enabled     bool   `json:"enabled"`
      Provider    string `json:"provider"`    // "openai", "ollama", "anthropic"
      BaseURL     string `json:"baseUrl"`
      APIKey      string `json:"apiKey"`      // encrypted
      Model       string `json:"model"`
      Temperature float64 `json:"temperature"`
  }
  ```
- [ ] **2.2** Add validation methods to `AISettings`:
  - [ ] `Validate()` — check required fields when enabled
  - [ ] `ValidateProvider()` — ensure provider is valid enum
  - [ ] `ValidateTemperature()` — ensure 0.0-1.0 range
- [ ] **2.3** Modify `core/settings.go` to include `AISettings` field
- [ ] **2.4** Add default values for AI settings
- [ ] **2.5** Implement API key encryption using existing PocketBase encryption helpers
- [ ] **2.6** Create migration (if needed) for settings schema update
- [ ] **2.7** Write unit tests for settings validation

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `core/ai_settings.go` | CREATE | AI settings struct and validation |
| `core/ai_settings_test.go` | CREATE | Unit tests for settings |
| `core/settings.go` | MODIFY | Add `AI AISettings` field |

### Tests
> ✅ **Unit Tests Required** — `core/ai_settings_test.go`

```go
// Test cases to implement:
func TestAISettings_Validate(t *testing.T)
func TestAISettings_ValidateProvider(t *testing.T)
func TestAISettings_ValidateTemperature(t *testing.T)
func TestAISettings_Defaults(t *testing.T)
func TestAISettings_APIKeyEncryption(t *testing.T)
```

**Test Scenarios:**
| Test | Input | Expected |
|------|-------|----------|
| Valid settings | enabled=true, provider="openai", apiKey="sk-xxx" | No error |
| Missing API key when enabled | enabled=true, apiKey="" | Error: "API key required" |
| Invalid provider | provider="invalid" | Error: "Invalid provider" |
| Temperature out of range | temperature=1.5 | Error: "Temperature must be 0-1" |
| Disabled settings skip validation | enabled=false, apiKey="" | No error |

### Verification
- [ ] `go test ./core/... -v` passes
- [ ] Settings can be saved/loaded from `_params` table

---

## PR #3: OpenAI Client & LLM Communication

**Branch:** `feat/ai-query-openai-client`  
**Estimated Time:** 5-6 hours  
**Dependencies:** PR #2  

### Description
Implement the HTTP client for communicating with OpenAI API (and compatible endpoints like Ollama).

### Tasks

- [ ] **3.1** Create `services/ai/openai_client.go`:
  - [ ] `NewOpenAIClient(settings AISettings)` constructor
  - [ ] `SendCompletion(ctx, systemPrompt, userMessage)` method
  - [ ] HTTP request building with proper headers
  - [ ] Response parsing (extract content from choices)
  - [ ] Error handling (API errors, timeouts, rate limits)
- [ ] **3.2** Implement timeout handling with `context.WithTimeout()` (30s default)
- [ ] **3.3** Add retry logic for transient failures (max 2 retries)
- [ ] **3.4** Create custom error types:
  - [ ] `AIClientError` — base error
  - [ ] `AIRateLimitError` — 429 responses
  - [ ] `AIAuthError` — 401 responses
  - [ ] `AITimeoutError` — context deadline exceeded
- [ ] **3.5** Write unit tests with mocked HTTP responses

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/ai/openai_client.go` | CREATE | OpenAI API client |
| `services/ai/openai_client_test.go` | CREATE | Client tests with mocks |
| `services/ai/errors.go` | CREATE | Custom error types |

### Tests
> ✅ **Unit Tests Required** — `services/ai/openai_client_test.go`

```go
// Test cases to implement:
func TestOpenAIClient_SendCompletion_Success(t *testing.T)
func TestOpenAIClient_SendCompletion_Timeout(t *testing.T)
func TestOpenAIClient_SendCompletion_RateLimit(t *testing.T)
func TestOpenAIClient_SendCompletion_AuthError(t *testing.T)
func TestOpenAIClient_SendCompletion_InvalidResponse(t *testing.T)
func TestOpenAIClient_Retry(t *testing.T)
```

**Mock Server Pattern:**
```go
func TestOpenAIClient_SendCompletion_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request format
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
        
        // Return mock response
        response := `{"choices":[{"message":{"content":"status = \"active\""}}]}`
        w.WriteHeader(200)
        w.Write([]byte(response))
    }))
    defer server.Close()
    
    client := NewOpenAIClient(AISettings{BaseURL: server.URL, APIKey: "test"})
    result, err := client.SendCompletion(context.Background(), "system", "user query")
    
    assert.NoError(t, err)
    assert.Equal(t, `status = "active"`, result)
}
```

### Verification
- [ ] `go test ./services/ai/... -v` passes
- [ ] Manual test with real OpenAI API key works

---

## PR #4: Schema Extraction & Prompt Building

**Branch:** `feat/ai-query-prompt-builder`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #3  

### Description
Build the system that extracts collection schemas and constructs optimized prompts for the LLM.

### Tasks

- [ ] **4.1** Create `services/ai/schema_extractor.go`:
  - [ ] `ExtractSchema(collection *core.Collection) string` — converts collection to prompt-friendly format
  - [ ] Handle all field types: text, number, bool, email, url, date, select, relation, file, json
  - [ ] Include relation target collection names
  - [ ] Include select field options
- [ ] **4.2** Create `services/ai/prompt_builder.go`:
  - [ ] `BuildSystemPrompt(schema string) string` — constructs full system prompt
  - [ ] Include PocketBase filter syntax rules
  - [ ] Include datetime macros documentation
  - [ ] Include few-shot examples
  - [ ] `BuildUserPrompt(query string) string` — wraps user query
- [ ] **4.3** Create prompt template as embedded string or file
- [ ] **4.4** Write comprehensive tests for schema extraction

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/ai/schema_extractor.go` | CREATE | Collection schema extraction |
| `services/ai/schema_extractor_test.go` | CREATE | Schema extraction tests |
| `services/ai/prompt_builder.go` | CREATE | Prompt construction |
| `services/ai/prompt_builder_test.go` | CREATE | Prompt building tests |
| `services/ai/prompt_template.go` | CREATE | System prompt template |

### Tests
> ✅ **Unit Tests Required** — Multiple test files

**`services/ai/schema_extractor_test.go`:**
```go
func TestExtractSchema_TextFields(t *testing.T)
func TestExtractSchema_NumberFields(t *testing.T)
func TestExtractSchema_SelectFields(t *testing.T)
func TestExtractSchema_RelationFields(t *testing.T)
func TestExtractSchema_AllFieldTypes(t *testing.T)
func TestExtractSchema_EmptyCollection(t *testing.T)
```

**`services/ai/prompt_builder_test.go`:**
```go
func TestBuildSystemPrompt_IncludesSchema(t *testing.T)
func TestBuildSystemPrompt_IncludesSyntaxRules(t *testing.T)
func TestBuildSystemPrompt_IncludesExamples(t *testing.T)
func TestBuildUserPrompt_WrapsQuery(t *testing.T)
```

**Test Scenarios:**
| Test | Input Collection | Expected Schema Output |
|------|-----------------|----------------------|
| Text field | `{name: "title", type: "text"}` | `title (text)` |
| Select field | `{name: "status", type: "select", options: ["active","inactive"]}` | `status (select: active\|inactive)` |
| Relation field | `{name: "author", type: "relation", collectionId: "users"}` | `author (relation → users)` |

### Verification
- [ ] `go test ./services/ai/... -v` passes
- [ ] Schema output is human-readable and LLM-friendly

---

## PR #5: Filter Validation & Query Execution

**Branch:** `feat/ai-query-validation`  
**Estimated Time:** 5-6 hours  
**Dependencies:** PR #4  

### Description
Implement validation layer to verify LLM-generated filters before execution, preventing hallucinated field names and invalid syntax.

### Tasks

- [ ] **5.1** Create `services/ai/filter_validator.go`:
  - [ ] `ValidateFilter(filter string, collection *core.Collection) error`
  - [ ] Extract field names from filter expression
  - [ ] Verify each field exists in collection schema
  - [ ] Check operator compatibility with field types
  - [ ] Validate datetime macro usage
- [ ] **5.2** Implement filter tokenizer/parser (basic):
  - [ ] Split on operators (`=`, `!=`, `>`, `<`, `~`, `&&`, `||`)
  - [ ] Extract field references (left side of operators)
  - [ ] Handle parentheses grouping
- [ ] **5.3** Create validation error messages:
  - [ ] `"Unknown field: {fieldName}. Available fields: {list}"`
  - [ ] `"Invalid operator '{op}' for field type '{type}'"`
  - [ ] `"Malformed filter syntax near: {context}"`
- [ ] **5.4** Write thorough validation tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/ai/filter_validator.go` | CREATE | Filter validation logic |
| `services/ai/filter_validator_test.go` | CREATE | Validation tests |
| `services/ai/filter_tokenizer.go` | CREATE | Basic filter parsing |

### Tests
> ✅ **Unit Tests Required** — `services/ai/filter_validator_test.go`

```go
func TestValidateFilter_ValidSimple(t *testing.T)
func TestValidateFilter_ValidComplex(t *testing.T)
func TestValidateFilter_UnknownField(t *testing.T)
func TestValidateFilter_InvalidOperator(t *testing.T)
func TestValidateFilter_MalformedSyntax(t *testing.T)
func TestValidateFilter_DatetimeMacros(t *testing.T)
func TestValidateFilter_RelationFields(t *testing.T)
```

**Test Scenarios:**
| Test | Filter | Collection Fields | Expected |
|------|--------|-------------------|----------|
| Valid simple | `status = "active"` | `[status]` | ✅ Pass |
| Unknown field | `invalid_field = "x"` | `[status, name]` | ❌ Error: Unknown field |
| Valid complex | `status = "active" && total > 100` | `[status, total]` | ✅ Pass |
| Invalid op for type | `name > 100` | `[name (text)]` | ❌ Error: Invalid operator |
| Datetime macro | `created >= @now - 86400` | `[created (date)]` | ✅ Pass |

### Verification
- [ ] `go test ./services/ai/... -v` passes
- [ ] Invalid filters are rejected with helpful messages

---

## PR #6: API Endpoint Implementation

**Branch:** `feat/ai-query-api`  
**Estimated Time:** 6-7 hours  
**Dependencies:** PR #5  

### Description
Create the `/api/ai/query` endpoint that ties together all backend components and exposes AI query functionality via REST API.

### Tasks

- [ ] **6.1** Create `apis/ai_query.go`:
  - [ ] Register route: `POST /api/ai/query`
  - [ ] Request validation (collection, query required)
  - [ ] Authentication check (require logged-in user or superuser)
  - [ ] Load AI settings, check if enabled
  - [ ] Load collection schema
- [ ] **6.2** Implement query flow:
  1. Extract schema from collection
  2. Build system + user prompts
  3. Call OpenAI client
  4. Validate generated filter
  5. Optionally execute filter and return results
  6. Return response with filter + results
- [ ] **6.3** Implement request/response structs:
  ```go
  type AIQueryRequest struct {
      Collection string `json:"collection"`
      Query      string `json:"query"`
      Execute    bool   `json:"execute"`
      Page       int    `json:"page"`
      PerPage    int    `json:"perPage"`
  }
  
  type AIQueryResponse struct {
      Filter     string        `json:"filter"`
      Results    []interface{} `json:"results,omitempty"`
      TotalItems int           `json:"totalItems,omitempty"`
      Page       int           `json:"page,omitempty"`
      PerPage    int           `json:"perPage,omitempty"`
      Error      string        `json:"error,omitempty"`
  }
  ```
- [ ] **6.4** Add collection API rule enforcement (respect listRule)
- [ ] **6.5** Implement error responses for all failure modes
- [ ] **6.6** Register API route in PocketBase app initialization
- [ ] **6.7** Write integration tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `apis/ai_query.go` | CREATE | API endpoint handler |
| `apis/ai_query_test.go` | CREATE | API endpoint tests |
| `apis/base.go` | MODIFY | Register AI query route |

### Tests
> ✅ **Integration Tests Required** — `apis/ai_query_test.go`

```go
func TestAIQueryAPI_Success(t *testing.T)
func TestAIQueryAPI_Unauthorized(t *testing.T)
func TestAIQueryAPI_AIDisabled(t *testing.T)
func TestAIQueryAPI_InvalidCollection(t *testing.T)
func TestAIQueryAPI_EmptyQuery(t *testing.T)
func TestAIQueryAPI_ExecuteResults(t *testing.T)
func TestAIQueryAPI_RespectsAPIRules(t *testing.T)
func TestAIQueryAPI_LLMError(t *testing.T)
func TestAIQueryAPI_ValidationError(t *testing.T)
```

**Integration Test Setup:**
```go
func setupTestApp(t *testing.T) *tests.TestApp {
    app, err := tests.NewTestApp()
    require.NoError(t, err)
    
    // Enable AI settings
    settings := app.Settings()
    settings.AI.Enabled = true
    settings.AI.Provider = "openai"
    settings.AI.APIKey = "test-key"
    settings.AI.Model = "gpt-4o-mini"
    app.Save(settings)
    
    return app
}
```

### Verification
- [ ] `go test ./apis/... -v` passes
- [ ] Manual API test with curl/Postman works:
  ```bash
  curl -X POST http://127.0.0.1:8090/api/ai/query \
    -H "Authorization: Bearer {token}" \
    -H "Content-Type: application/json" \
    -d '{"collection":"posts","query":"recent posts","execute":true}'
  ```

---

## PR #7: Admin UI — AI Query Sidebar Panel

**Branch:** `feat/ai-query-ui-panel`  
**Estimated Time:** 6-8 hours  
**Dependencies:** PR #6  

### Description
Build the Svelte components for the AI Query sidebar panel in the Admin UI.

### Tasks

- [ ] **7.1** Create `ui/src/stores/ai.js`:
  - [ ] `aiQuery` store (current query text)
  - [ ] `aiFilter` store (generated filter)
  - [ ] `aiResults` store (query results)
  - [ ] `aiLoading` store (loading state)
  - [ ] `aiError` store (error message)
- [ ] **7.2** Create `ui/src/components/ai/AIQueryInput.svelte`:
  - [ ] Textarea for natural language query
  - [ ] Collection dropdown selector
  - [ ] Search button with loading state
  - [ ] Keyboard shortcut (Ctrl+Enter to search)
- [ ] **7.3** Create `ui/src/components/ai/AIFilterDisplay.svelte`:
  - [ ] Display generated filter in code block
  - [ ] Copy to clipboard button
  - [ ] "Apply Filter" button (navigates to collection with filter)
- [ ] **7.4** Create `ui/src/components/ai/AIQueryResults.svelte`:
  - [ ] Results count display
  - [ ] Basic record list preview (id, first few fields)
  - [ ] "View in Collection" link
- [ ] **7.5** Create `ui/src/components/ai/AIQueryPanel.svelte`:
  - [ ] Combines Input, Filter, Results components
  - [ ] Handles API calls to `/api/ai/query`
  - [ ] Error display
- [ ] **7.6** Modify `ui/src/App.svelte`:
  - [ ] Add "AI Query" entry to sidebar navigation
  - [ ] Add route for AI Query panel
- [ ] **7.7** Style components to match PocketBase Admin UI design

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/stores/ai.js` | CREATE | AI state management |
| `ui/src/components/ai/AIQueryInput.svelte` | CREATE | Query input component |
| `ui/src/components/ai/AIFilterDisplay.svelte` | CREATE | Filter display component |
| `ui/src/components/ai/AIQueryResults.svelte` | CREATE | Results component |
| `ui/src/components/ai/AIQueryPanel.svelte` | CREATE | Main panel component |
| `ui/src/App.svelte` | MODIFY | Add sidebar entry |
| `ui/src/scss/_ai.scss` | CREATE | AI component styles |

### Tests
> ❌ **No automated tests** — UI components tested manually.

**Manual Test Checklist:**
- [ ] AI Query appears in sidebar when AI is enabled
- [ ] Collection dropdown populates with all collections
- [ ] Query input accepts text and submits on button click
- [ ] Loading spinner shows during API call
- [ ] Generated filter displays correctly
- [ ] Copy button copies filter to clipboard
- [ ] Apply Filter navigates to collection with filter applied
- [ ] Errors display clearly in UI
- [ ] UI matches PocketBase design language

### Verification
- [ ] `npm run build` succeeds in `/ui`
- [ ] Rebuilt Go binary includes new UI
- [ ] AI Query panel renders and functions in browser

---

## PR #8: Admin UI — AI Settings Page

**Branch:** `feat/ai-query-ui-settings`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #7  

### Description
Build the Settings page for configuring AI Query feature (provider, API key, model, etc.).

### Tasks

- [ ] **8.1** Create `ui/src/pages/settings/AI.svelte`:
  - [ ] Enable/Disable toggle
  - [ ] Provider dropdown (OpenAI, Ollama, Anthropic, Custom)
  - [ ] API Base URL input (auto-fills based on provider)
  - [ ] API Key input (password field)
  - [ ] Model dropdown/input
  - [ ] Temperature slider (0.0 - 1.0)
  - [ ] Save button
- [ ] **8.2** Create `ui/src/components/ai/AISettingsForm.svelte`:
  - [ ] Reusable form component
  - [ ] Field validation
  - [ ] "Test Connection" button
- [ ] **8.3** Implement Test Connection functionality:
  - [ ] Call backend endpoint to verify LLM connectivity
  - [ ] Show success/failure toast message
- [ ] **8.4** Add AI Settings to settings navigation
- [ ] **8.5** Implement settings save/load via PocketBase API
- [ ] **8.6** Add conditional UI (hide API key field for Ollama)

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/pages/settings/AI.svelte` | CREATE | Settings page |
| `ui/src/components/ai/AISettingsForm.svelte` | CREATE | Settings form component |
| `ui/src/pages/settings/Index.svelte` | MODIFY | Add AI link to settings nav |
| `apis/ai_settings.go` | CREATE | Settings API endpoints |

### Tests
> ❌ **No automated tests** — Settings UI tested manually.

**Manual Test Checklist:**
- [ ] AI Settings page accessible from Settings menu
- [ ] Toggle enables/disables AI feature
- [ ] Provider selection updates default Base URL
- [ ] API Key field is password-masked
- [ ] Temperature slider works with 0.1 increments
- [ ] Test Connection shows success for valid config
- [ ] Test Connection shows error for invalid API key
- [ ] Settings persist after save and page reload
- [ ] AI Query panel hidden when AI disabled

### Verification
- [ ] Settings can be saved and retrieved
- [ ] Test Connection works with real OpenAI API key
- [ ] UI enables/disables based on settings

---

## PR #9: Documentation & Final Polish

**Branch:** `feat/ai-query-docs`  
**Estimated Time:** 3-4 hours  
**Dependencies:** PR #8  

### Description
Complete documentation, final bug fixes, and prepare for demo.

### Tasks

- [ ] **9.1** Complete `docs/AI_QUERY_FEATURE.md`:
  - [ ] Feature overview
  - [ ] Setup instructions
  - [ ] Configuration guide
  - [ ] API reference
  - [ ] Troubleshooting guide
- [ ] **9.2** Update main `README.md`:
  - [ ] Add AI Query to features list
  - [ ] Add quick start guide
  - [ ] Add screenshots
- [ ] **9.3** Create `CHANGELOG.md` entry for AI Query feature
- [ ] **9.4** Final code review and cleanup:
  - [ ] Remove debug logging
  - [ ] Fix any TODO comments
  - [ ] Ensure consistent error messages
- [ ] **9.5** Run full test suite:
  ```powershell
  go test ./... -v
  ```
- [ ] **9.6** Build final release binary
- [ ] **9.7** Record demo video (5 minutes)

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `docs/AI_QUERY_FEATURE.md` | MODIFY | Complete documentation |
| `README.md` | MODIFY | Add feature to main readme |
| `CHANGELOG.md` | MODIFY | Add release notes |

### Tests
> ✅ **Full Test Suite Run Required**

```powershell
# Run all tests
go test ./... -v -cover

# Expected output: All tests pass, >80% coverage on new code
```

### Verification
- [ ] All tests pass
- [ ] Documentation is complete and accurate
- [ ] Demo video recorded successfully
- [ ] Feature works end-to-end

---

## Summary: Test Coverage Matrix

| PR | Unit Tests | Integration Tests | Manual Tests |
|----|------------|-------------------|--------------|
| PR #1: Setup | ❌ | ❌ | ✅ Build verification |
| PR #2: Settings | ✅ `ai_settings_test.go` | ❌ | ❌ |
| PR #3: OpenAI Client | ✅ `openai_client_test.go` | ❌ | ✅ Real API test |
| PR #4: Prompt Builder | ✅ `schema_extractor_test.go`, `prompt_builder_test.go` | ❌ | ❌ |
| PR #5: Validation | ✅ `filter_validator_test.go` | ❌ | ❌ |
| PR #6: API Endpoint | ✅ `ai_query_test.go` | ✅ Full API tests | ✅ curl/Postman |
| PR #7: UI Panel | ❌ | ❌ | ✅ Full UI testing |
| PR #8: UI Settings | ❌ | ❌ | ✅ Full UI testing |
| PR #9: Docs | ❌ | ✅ Full suite run | ✅ Demo recording |

---

## Quick Reference: All New Files

```
NEW FILES (18):
├── apis/ai_query.go
├── apis/ai_query_test.go
├── apis/ai_settings.go
├── core/ai_settings.go
├── core/ai_settings_test.go
├── services/ai/openai_client.go
├── services/ai/openai_client_test.go
├── services/ai/prompt_builder.go
├── services/ai/prompt_builder_test.go
├── services/ai/prompt_template.go
├── services/ai/schema_extractor.go
├── services/ai/schema_extractor_test.go
├── services/ai/filter_validator.go
├── services/ai/filter_validator_test.go
├── services/ai/errors.go
├── services/ai/filter_tokenizer.go
├── ui/src/stores/ai.js
├── ui/src/components/ai/AIQueryInput.svelte
├── ui/src/components/ai/AIFilterDisplay.svelte
├── ui/src/components/ai/AIQueryResults.svelte
├── ui/src/components/ai/AIQueryPanel.svelte
├── ui/src/components/ai/AISettingsForm.svelte
├── ui/src/pages/settings/AI.svelte
├── ui/src/scss/_ai.scss
├── docs/AI_QUERY_FEATURE.md

MODIFIED FILES (5):
├── core/settings.go
├── apis/base.go
├── ui/src/App.svelte
├── ui/src/pages/settings/Index.svelte
├── README.md
```

---

## Execution Order

```
PR #1 (Setup) 
    ↓
PR #2 (Settings) ←── Unit tests
    ↓
PR #3 (OpenAI Client) ←── Unit tests + Mock server
    ↓
PR #4 (Prompt Builder) ←── Unit tests
    ↓
PR #5 (Validation) ←── Unit tests
    ↓
PR #6 (API Endpoint) ←── Integration tests
    ↓
PR #7 (UI Panel) ←── Manual testing
    ↓
PR #8 (UI Settings) ←── Manual testing
    ↓
PR #9 (Docs) ←── Full test suite
```

---

**Document Status:** Ready for implementation  
**Total Estimated Time:** 35-45 hours  
**Total PRs:** 9  
**Total Test Files:** 7
