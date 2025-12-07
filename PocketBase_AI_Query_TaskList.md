# PocketBase AI Query Assistant — Task List

## Project Overview

**Repository:** Fork of https://github.com/pocketbase/pocketbase  
**Feature:** AI-powered natural language query assistant  
**Total PRs:** 22 (V1: 9 ✅, V2: 8 ✅, V3: 5 🚧)

---

## Version Summary

| Version | PRs | Status | Features |
|---------|-----|--------|----------|
| V1 | #1-9 | ✅ Complete | Single-collection AI Query, Filter generation, Settings UI |
| V2 | #10-17 | ✅ Complete | Multi-table queries, Dual output, SQL Terminal |
| V3 | #18-22 | 🚧 In Progress | Real-time Metrics Dashboard, Data Import Wizard |  

---

## File Structure Overview

### V1 Files (✅ Complete)

```
pocketbase/                          # Forked repository root
├── apis/
│   ├── ai_query.go                  # ✅ API endpoint handler
│   └── ai_query_test.go             # ✅ API endpoint tests
├── core/
│   ├── ai_settings.go               # ✅ AI settings struct & validation
│   ├── ai_settings_test.go          # ✅ Settings tests
│   └── settings.go                  # ✅ MODIFIED — Added AI settings
├── services/
│   └── ai/
│       ├── openai_client.go         # ✅ OpenAI API client
│       ├── openai_client_test.go    # ✅ Client tests (mocked)
│       ├── prompt_builder.go        # ✅ System prompt construction
│       ├── prompt_builder_test.go   # ✅ Prompt tests
│       ├── prompt_template.go       # ✅ Prompt template
│       ├── schema_extractor.go      # ✅ Collection schema extraction
│       ├── schema_extractor_test.go # ✅ Schema extraction tests
│       ├── filter_validator.go      # ✅ Filter syntax validation
│       ├── filter_validator_test.go # ✅ Validation tests
│       ├── filter_tokenizer.go      # ✅ Filter parsing
│       └── errors.go                # ✅ Custom error types
├── ui/src/
│   ├── components/ai/
│   │   ├── AIQueryPanel.svelte      # ✅ Main sidebar panel
│   │   ├── AIQueryInput.svelte      # ✅ Query input component
│   │   ├── AIQueryResults.svelte    # ✅ Results display
│   │   ├── AIFilterDisplay.svelte   # ✅ Filter with copy button
│   │   └── AISettingsForm.svelte    # ✅ Settings form component
│   ├── pages/settings/
│   │   └── AI.svelte                # ✅ AI settings page
│   ├── stores/
│   │   └── ai.js                    # ✅ AI-related state store
│   └── App.svelte                   # ✅ MODIFIED — Added sidebar entry
└── docs/
    └── AI_QUERY_FEATURE.md          # ✅ Feature documentation
```

### V2 Files (✅ Complete)

```
pocketbase/
├── apis/
│   └── sql_terminal.go              # ✅ SQL Terminal API endpoints
├── services/
│   ├── ai/
│   │   ├── schema_extractor.go      # ✅ Extract ALL collections + relationships
│   │   └── prompt_template.go       # ✅ Add SQL syntax rules
│   └── sql/                         # ✅ COMPLETE
│       ├── parser.go                # ✅ SQL statement parser
│       ├── executor.go              # ✅ SQL execution via PocketBase APIs
│       └── mapper.go                # ✅ SQL type → PocketBase field mapper
├── ui/src/
│   ├── components/ai/
│   │   ├── AIQueryPanel.svelte      # ✅ Dual output tabs
│   │   ├── QueryTabs.svelte         # ✅ Filter/SQL tab switcher
│   │   └── EditableCodeBlock.svelte # ✅ Editable query with syntax highlight
│   ├── components/sql/              # ✅ COMPLETE
│   │   ├── SQLEditor.svelte         # ✅ Code editor component
│   │   ├── SchemaExplorer.svelte    # ✅ Collections sidebar browser
│   │   ├── ResultsTable.svelte      # ✅ Dynamic results display
│   │   └── QueryHistory.svelte      # ✅ Command history dropdown
│   ├── pages/
│   │   └── SQLTerminal.svelte       # ✅ Main SQL Terminal page
│   ├── stores/
│   │   └── sql.js                   # ✅ SQL terminal state
│   └── App.svelte                   # ✅ SQL Terminal in sidebar
└── docs/
    └── SQL_TERMINAL_FEATURE.md      # ✅ SQL Terminal documentation
```

### V3 Files (🚧 Planned)

```
pocketbase/
├── apis/
│   ├── metrics.go                   # NEW — Dashboard metrics API endpoints
│   └── import.go                    # NEW — Data import API endpoints
├── ui/src/
│   ├── pages/
│   │   ├── Dashboard.svelte         # NEW — Metrics dashboard page
│   │   └── ImportWizard.svelte      # NEW — Data import wizard page
│   ├── components/
│   │   ├── dashboard/               # NEW DIRECTORY
│   │   │   ├── MetricCard.svelte    # NEW — Stat card with icon
│   │   │   ├── RequestsChart.svelte # NEW — Time-series line chart
│   │   │   ├── LatencyChart.svelte  # NEW — Multi-line percentile chart
│   │   │   ├── EndpointsChart.svelte# NEW — Horizontal bar chart
│   │   │   └── CollectionsTable.svelte # NEW — Record counts table
│   │   └── import/                  # NEW DIRECTORY
│   │       ├── FileUpload.svelte    # NEW — Drag-drop file zone
│   │       ├── DataPreview.svelte   # NEW — Preview table
│   │       ├── FieldMapper.svelte   # NEW — Column mapping UI
│   │       └── ImportProgress.svelte# NEW — Progress + errors
│   ├── stores/
│   │   ├── dashboard.js             # NEW — Dashboard state management
│   │   └── import.js                # NEW — Import wizard state
│   └── scss/
│       ├── _dashboard.scss          # NEW — Dashboard styles
│       └── _import.scss             # NEW — Import styles
└── docs/
    ├── DASHBOARD_FEATURE.md         # NEW — Dashboard documentation
    └── IMPORT_FEATURE.md            # NEW — Import documentation
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
- [x] All tests pass
- [x] Documentation is complete and accurate
- [ ] Demo video recorded successfully (optional)
- [x] Feature works end-to-end

---

# V2 PRs: Multi-Table SQL Queries & SQL Terminal

---

## PR #10: Multi-Collection Schema Extraction

**Branch:** `feat/ai-query-multi-schema`  
**Estimated Time:** 4-5 hours  
**Dependencies:** V1 Complete (PR #9)  

### Description
Extend schema extraction to include ALL collections and their relationships, enabling multi-table queries.

### Tasks

- [ ] **10.1** Modify `services/ai/schema_extractor.go`:
  - [ ] `ExtractAllSchemas(app *pocketbase.PocketBase) string` — extracts all collection schemas
  - [ ] Include relationship mappings between collections
  - [ ] Format schema for LLM understanding of JOINs
- [ ] **10.2** Create relationship detection:
  - [ ] Parse relation fields to identify foreign keys
  - [ ] Build relationship map (e.g., `orders.customer → customers.id`)
- [ ] **10.3** Update prompt template with multi-table examples
- [ ] **10.4** Write unit tests for multi-collection schema

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/ai/schema_extractor.go` | MODIFY | Add multi-collection extraction |
| `services/ai/schema_extractor_test.go` | MODIFY | Add multi-collection tests |
| `services/ai/prompt_template.go` | MODIFY | Add SQL syntax and multi-table examples |

### Tests
> ✅ **Unit Tests Required**

```go
func TestExtractAllSchemas(t *testing.T)
func TestExtractRelationships(t *testing.T)
func TestSchemaFormatForJoins(t *testing.T)
```

---

## PR #11: Dual Output Backend (Filter + SQL)

**Branch:** `feat/ai-query-dual-output`  
**Estimated Time:** 5-6 hours  
**Dependencies:** PR #10  

### Description
Modify AI Query API to return BOTH PocketBase filter AND SQL for queries where both are possible.

### Tasks

- [ ] **11.1** Update `apis/ai_query.go`:
  - [ ] Generate both filter and SQL outputs
  - [ ] Add `canUseFilter` field to response
  - [ ] Add `sql` field to response
  - [ ] Detect when query requires SQL-only (JOINs, aggregates)
- [ ] **11.2** Update prompt template to request dual output
- [ ] **11.3** Implement query complexity detection:
  - [ ] Simple (single table, basic conditions) → Filter works
  - [ ] Complex (JOINs, GROUP BY, aggregates) → SQL only
- [ ] **11.4** Update response schema
- [ ] **11.5** Write integration tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `apis/ai_query.go` | MODIFY | Add dual output logic |
| `apis/ai_query_test.go` | MODIFY | Add dual output tests |
| `services/ai/prompt_template.go` | MODIFY | Request dual output from LLM |

### Tests
> ✅ **Integration Tests Required**

```go
func TestAIQueryAPI_DualOutput_SimpleQuery(t *testing.T)
func TestAIQueryAPI_DualOutput_ComplexQuery(t *testing.T)
func TestAIQueryAPI_SQLOnlyForJoins(t *testing.T)
```

---

## PR #12: Editable Query UI with Tabs

**Branch:** `feat/ai-query-editable-ui`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #11  

### Description
Update AI Query panel with tabbed interface (Filter/SQL) and editable query blocks.

### Tasks

- [ ] **12.1** Create `ui/src/components/ai/QueryTabs.svelte`:
  - [ ] Tab component with Filter/SQL options
  - [ ] Active tab state management
  - [ ] Disable tab when option not available
- [ ] **12.2** Create `ui/src/components/ai/EditableCodeBlock.svelte`:
  - [ ] Textarea with syntax highlighting (basic)
  - [ ] Edit mode toggle
  - [ ] Re-execute button after editing
- [ ] **12.3** Modify `AIQueryPanel.svelte`:
  - [ ] Integrate tabs component
  - [ ] Handle dual response (filter + SQL)
  - [ ] Show appropriate output in each tab
- [ ] **12.4** Modify `AIFilterDisplay.svelte`:
  - [ ] Make filter editable
  - [ ] Add execute button
- [ ] **12.5** Update stores for dual state

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/components/ai/QueryTabs.svelte` | CREATE | Tab switcher component |
| `ui/src/components/ai/EditableCodeBlock.svelte` | CREATE | Editable code block |
| `ui/src/components/ai/AIQueryPanel.svelte` | MODIFY | Add tabs and dual output |
| `ui/src/components/ai/AIFilterDisplay.svelte` | MODIFY | Make editable |
| `ui/src/stores/ai.js` | MODIFY | Add SQL state |

### Tests
> ❌ **No automated tests** — Manual UI testing

**Manual Test Checklist:**
- [ ] Filter tab shows PocketBase filter syntax
- [ ] SQL tab shows SQL query
- [ ] Tabs switch correctly
- [ ] SQL tab disabled when filter-only query
- [ ] Editing filter and re-executing works
- [ ] Editing SQL and re-executing works

---

## PR #13: SQL Parser & Type Mapper

**Branch:** `feat/sql-parser`  
**Estimated Time:** 6-7 hours  
**Dependencies:** PR #10  

### Description
Create SQL parser to understand SQL statements and map SQL types to PocketBase field types.

### Tasks

- [ ] **13.1** Create `services/sql/parser.go`:
  - [ ] `ParseSQL(sql string) (*SQLStatement, error)`
  - [ ] Detect statement type (SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP)
  - [ ] Extract table names
  - [ ] Extract column definitions (for CREATE TABLE)
  - [ ] Extract WHERE clauses
- [ ] **13.2** Create `services/sql/mapper.go`:
  - [ ] `MapSQLType(sqlType string) string` — returns PocketBase field type
  - [ ] Handle TEXT → text, INTEGER → number, REAL → number, etc.
  - [ ] Handle REFERENCES → relation
  - [ ] Handle CHECK(IN(...)) → select
- [ ] **13.3** Create SQL statement structs:
  ```go
  type SQLStatement struct {
      Type       string   // SELECT, INSERT, CREATE, etc.
      Tables     []string
      Columns    []ColumnDef
      Where      string
      Values     []interface{}
  }
  
  type ColumnDef struct {
      Name       string
      Type       string
      Required   bool
      Reference  string // For relations
      Options    []string // For select fields
  }
  ```
- [ ] **13.4** Write comprehensive parser tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/sql/parser.go` | CREATE | SQL statement parser |
| `services/sql/parser_test.go` | CREATE | Parser tests |
| `services/sql/mapper.go` | CREATE | SQL → PocketBase type mapper |
| `services/sql/mapper_test.go` | CREATE | Mapper tests |

### Tests
> ✅ **Unit Tests Required**

```go
func TestParseSQL_Select(t *testing.T)
func TestParseSQL_Insert(t *testing.T)
func TestParseSQL_Update(t *testing.T)
func TestParseSQL_Delete(t *testing.T)
func TestParseSQL_CreateTable(t *testing.T)
func TestParseSQL_AlterTable(t *testing.T)
func TestParseSQL_DropTable(t *testing.T)
func TestMapSQLType_Text(t *testing.T)
func TestMapSQLType_Number(t *testing.T)
func TestMapSQLType_Relation(t *testing.T)
func TestMapSQLType_Select(t *testing.T)
```

---

## PR #14: SQL Executor (PocketBase API Integration)

**Branch:** `feat/sql-executor`  
**Estimated Time:** 6-7 hours  
**Dependencies:** PR #13  

### Description
Execute parsed SQL statements using PocketBase APIs to create real collections and records.

### Tasks

- [ ] **14.1** Create `services/sql/executor.go`:
  - [ ] `ExecuteSQL(app *pocketbase.PocketBase, stmt *SQLStatement) (*ExecutionResult, error)`
  - [ ] Route to appropriate handler based on statement type
- [ ] **14.2** Implement CREATE TABLE handler:
  - [ ] Convert parsed columns to PocketBase fields
  - [ ] Create collection via PocketBase Collection API
  - [ ] Return created collection info
- [ ] **14.3** Implement ALTER TABLE handler:
  - [ ] Add/modify/drop fields in existing collection
- [ ] **14.4** Implement DROP TABLE handler:
  - [ ] Delete collection via PocketBase API
- [ ] **14.5** Implement INSERT handler:
  - [ ] Create record via PocketBase Records API
  - [ ] Return created record ID
- [ ] **14.6** Implement UPDATE handler:
  - [ ] Update records via PocketBase Records API
  - [ ] Return affected row count
- [ ] **14.7** Implement DELETE handler:
  - [ ] Delete records via PocketBase Records API
  - [ ] Return affected row count
- [ ] **14.8** Implement SELECT handler:
  - [ ] Execute query directly against SQLite
  - [ ] Return results with column names
- [ ] **14.9** Add query timeout and result limits

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `services/sql/executor.go` | CREATE | SQL execution engine |
| `services/sql/executor_test.go` | CREATE | Executor tests |

### Tests
> ✅ **Integration Tests Required**

```go
func TestExecuteSQL_CreateTable(t *testing.T)
func TestExecuteSQL_Insert(t *testing.T)
func TestExecuteSQL_Update(t *testing.T)
func TestExecuteSQL_Delete(t *testing.T)
func TestExecuteSQL_Select(t *testing.T)
func TestExecuteSQL_SelectWithJoin(t *testing.T)
func TestExecuteSQL_Timeout(t *testing.T)
```

---

## PR #15: SQL Terminal API Endpoints

**Branch:** `feat/sql-terminal-api`  
**Estimated Time:** 5-6 hours  
**Dependencies:** PR #14  

### Description
Create API endpoints for SQL Terminal functionality.

### Tasks

- [ ] **15.1** Create `apis/sql_terminal.go`:
  - [ ] `POST /api/sql/execute` — Execute raw SQL
  - [ ] `POST /api/sql/ai` — AI mode (natural language → SQL)
  - [ ] `GET /api/sql/history` — Get query history (optional)
- [ ] **15.2** Implement execute endpoint:
  - [ ] Parse SQL using parser
  - [ ] Execute using executor
  - [ ] Return structured response
- [ ] **15.3** Implement AI mode endpoint:
  - [ ] Build SQL-focused prompt
  - [ ] Call LLM to generate SQL
  - [ ] Optionally execute generated SQL
  - [ ] Return SQL + results
- [ ] **15.4** Add authentication (require logged-in user)
- [ ] **15.5** Add confirmation requirement for destructive operations
- [ ] **15.6** Register routes in `apis/base.go`
- [ ] **15.7** Write integration tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `apis/sql_terminal.go` | CREATE | SQL Terminal API endpoints |
| `apis/sql_terminal_test.go` | CREATE | API tests |
| `apis/base.go` | MODIFY | Register SQL routes |

### Tests
> ✅ **Integration Tests Required**

```go
func TestSQLTerminal_Execute_Select(t *testing.T)
func TestSQLTerminal_Execute_CreateTable(t *testing.T)
func TestSQLTerminal_Execute_Insert(t *testing.T)
func TestSQLTerminal_AI_GenerateSQL(t *testing.T)
func TestSQLTerminal_Unauthorized(t *testing.T)
```

---

## PR #16: SQL Terminal UI

**Branch:** `feat/sql-terminal-ui`  
**Estimated Time:** 8-10 hours  
**Dependencies:** PR #15  

### Description
Build the SQL Terminal frontend page with code editor, schema browser, and results display.

### Tasks

- [ ] **16.1** Create `ui/src/pages/SQLTerminal.svelte`:
  - [ ] Main page layout with sidebar and editor
  - [ ] AI Mode / SQL Mode toggle
  - [ ] Integration with API endpoints
- [ ] **16.2** Create `ui/src/components/sql/SQLEditor.svelte`:
  - [ ] Textarea with basic syntax highlighting
  - [ ] Line numbers
  - [ ] Keyboard shortcuts (Ctrl+Enter to run)
  - [ ] Auto-complete for table/column names (basic)
- [ ] **16.3** Create `ui/src/components/sql/SchemaExplorer.svelte`:
  - [ ] Tree view of collections
  - [ ] Expandable to show fields
  - [ ] Click to insert table/field name
- [ ] **16.4** Create `ui/src/components/sql/ResultsTable.svelte`:
  - [ ] Dynamic column headers from query
  - [ ] Scrollable data rows
  - [ ] Export to CSV button
  - [ ] Export to JSON button
- [ ] **16.5** Create `ui/src/components/sql/QueryHistory.svelte`:
  - [ ] Dropdown of recent queries
  - [ ] Click to restore query
  - [ ] Stored in localStorage
- [ ] **16.6** Create `ui/src/stores/sql.js`:
  - [ ] Current query
  - [ ] Query results
  - [ ] Loading state
  - [ ] History
  - [ ] Mode (AI/SQL)
- [ ] **16.7** Modify `ui/src/App.svelte`:
  - [ ] Add SQL Terminal to sidebar navigation
  - [ ] Add route
- [ ] **16.8** Create styles for SQL components
- [ ] **16.9** Add confirmation dialogs for destructive operations

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/pages/SQLTerminal.svelte` | CREATE | Main SQL Terminal page |
| `ui/src/components/sql/SQLEditor.svelte` | CREATE | Code editor |
| `ui/src/components/sql/SchemaExplorer.svelte` | CREATE | Schema browser |
| `ui/src/components/sql/ResultsTable.svelte` | CREATE | Results display |
| `ui/src/components/sql/QueryHistory.svelte` | CREATE | History dropdown |
| `ui/src/stores/sql.js` | CREATE | SQL state management |
| `ui/src/App.svelte` | MODIFY | Add sidebar entry |
| `ui/src/scss/_sql.scss` | CREATE | SQL component styles |

### Tests
> ❌ **No automated tests** — Manual UI testing

**Manual Test Checklist:**
- [ ] SQL Terminal accessible from sidebar
- [ ] SQL Mode: Can type and execute raw SQL
- [ ] AI Mode: Can type natural language and get SQL
- [ ] CREATE TABLE creates real collection (visible in Admin UI)
- [ ] INSERT creates real records
- [ ] UPDATE modifies records
- [ ] DELETE removes records
- [ ] SELECT returns results in table
- [ ] Schema browser shows all collections
- [ ] Query history saves and restores queries
- [ ] Export to CSV works
- [ ] Export to JSON works
- [ ] Confirmation dialog for DROP/DELETE
- [ ] Error messages display clearly

---

## PR #17: V2 Documentation & Polish

**Branch:** `feat/v2-docs`  
**Estimated Time:** 3-4 hours  
**Dependencies:** PR #12, PR #16  

### Description
Complete V2 documentation and final testing.

### Tasks

- [ ] **17.1** Create `docs/SQL_TERMINAL_FEATURE.md`:
  - [ ] Feature overview
  - [ ] Setup instructions
  - [ ] Usage guide (AI Mode vs SQL Mode)
  - [ ] SQL → PocketBase mapping reference
  - [ ] Security considerations
- [ ] **17.2** Update `docs/AI_QUERY_FEATURE.md`:
  - [ ] Add multi-table query examples
  - [ ] Document dual output feature
  - [ ] Add editable query documentation
- [ ] **17.3** Update `README.md`:
  - [ ] Add SQL Terminal to features
  - [ ] Add screenshots
- [ ] **17.4** Update `CHANGELOG.md`:
  - [ ] Add V2 release notes
- [ ] **17.5** Run full test suite
- [ ] **17.6** Build final release binary
- [ ] **17.7** End-to-end testing of all features

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `docs/SQL_TERMINAL_FEATURE.md` | CREATE | SQL Terminal documentation |
| `docs/AI_QUERY_FEATURE.md` | MODIFY | Add V2 features |
| `README.md` | MODIFY | Add SQL Terminal |
| `CHANGELOG.md` | MODIFY | Add V2 release notes |

### Verification
- [x] All tests pass
- [x] Documentation is complete
- [x] Multi-table queries work end-to-end
- [x] SQL Terminal creates real collections
- [x] All changes visible in Admin UI

---

# V3 PRs: Real-time Dashboard & Data Import Wizard

---

## PR #18: Metrics Backend API

**Branch:** `feat/metrics-api`  
**Estimated Time:** 3-4 hours  
**Dependencies:** V2 Complete (PR #17)  

### Description
Create backend API endpoints for the metrics dashboard, aggregating data from the existing `_logs` table.

### Tasks

- [ ] **18.1** Create `apis/metrics.go`:
  - [ ] `GET /api/metrics/overview` — Total requests, avg latency, error rate, DB size
  - [ ] `GET /api/metrics/requests` — Time-series of requests per minute (last 24h)
  - [ ] `GET /api/metrics/latency` — Avg/p50/p95/p99 latency over time
  - [ ] `GET /api/metrics/errors` — Error count by status code over time
  - [ ] `GET /api/metrics/endpoints` — Top 10 endpoints by request count
  - [ ] `GET /api/metrics/collections` — Record counts per collection
- [ ] **18.2** Implement metrics aggregation from `_logs` table:
  - [ ] Query `data.execTime` for latency metrics
  - [ ] Query `data.status` for error rate calculation
  - [ ] Query `data.url` for endpoint statistics
- [ ] **18.3** Implement database size calculation:
  - [ ] Use SQLite PRAGMA: `page_count * page_size`
- [ ] **18.4** Implement collection stats:
  - [ ] Loop through all collections
  - [ ] Use `CountRecords()` for each
- [ ] **18.5** Add authentication (superuser only)
- [ ] **18.6** Register routes in `apis/base.go`
- [ ] **18.7** Write unit tests for aggregation functions

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `apis/metrics.go` | CREATE | Metrics API endpoints |
| `apis/base.go` | MODIFY | Register metrics routes |

### Tests
> ✅ **Unit Tests Required** — `apis/metrics_test.go`

```go
func TestMetricsAPI_Overview(t *testing.T)
func TestMetricsAPI_Requests(t *testing.T)
func TestMetricsAPI_Latency(t *testing.T)
func TestMetricsAPI_Errors(t *testing.T)
func TestMetricsAPI_Endpoints(t *testing.T)
func TestMetricsAPI_Collections(t *testing.T)
func TestMetricsAPI_Unauthorized(t *testing.T)
```

---

## PR #19: Dashboard UI

**Branch:** `feat/dashboard-ui`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #18  

### Description
Build the metrics dashboard frontend with Chart.js visualizations.

### Tasks

- [ ] **19.1** Create `ui/src/pages/Dashboard.svelte`:
  - [ ] Main page layout with grid
  - [ ] Integration with metrics API
  - [ ] Auto-refresh every 30 seconds
  - [ ] Time range selector (1h, 6h, 24h)
- [ ] **19.2** Create `ui/src/components/dashboard/MetricCard.svelte`:
  - [ ] Stat card with value, label, icon
  - [ ] Trend indicator (up/down)
- [ ] **19.3** Create `ui/src/components/dashboard/RequestsChart.svelte`:
  - [ ] Line chart using Chart.js (already in PocketBase)
  - [ ] Requests over time
- [ ] **19.4** Create `ui/src/components/dashboard/LatencyChart.svelte`:
  - [ ] Multi-line chart for p50/p95/p99
  - [ ] Different colors per percentile
- [ ] **19.5** Create `ui/src/components/dashboard/EndpointsChart.svelte`:
  - [ ] Horizontal bar chart
  - [ ] Top 10 endpoints
- [ ] **19.6** Create `ui/src/components/dashboard/CollectionsTable.svelte`:
  - [ ] Table with collection name, record count
  - [ ] Sortable columns
- [ ] **19.7** Create `ui/src/stores/dashboard.js`:
  - [ ] Metrics data stores
  - [ ] Loading state
  - [ ] Refresh interval
- [ ] **19.8** Modify `ui/src/App.svelte`:
  - [ ] Add Dashboard to sidebar with chart icon
  - [ ] Add route `/dashboard`
- [ ] **19.9** Create `ui/src/scss/_dashboard.scss`:
  - [ ] Grid layout styles
  - [ ] Card styles
  - [ ] Chart container styles

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/pages/Dashboard.svelte` | CREATE | Main dashboard page |
| `ui/src/components/dashboard/MetricCard.svelte` | CREATE | Stat card component |
| `ui/src/components/dashboard/RequestsChart.svelte` | CREATE | Requests line chart |
| `ui/src/components/dashboard/LatencyChart.svelte` | CREATE | Latency chart |
| `ui/src/components/dashboard/EndpointsChart.svelte` | CREATE | Endpoints bar chart |
| `ui/src/components/dashboard/CollectionsTable.svelte` | CREATE | Collections table |
| `ui/src/stores/dashboard.js` | CREATE | Dashboard state |
| `ui/src/scss/_dashboard.scss` | CREATE | Dashboard styles |
| `ui/src/App.svelte` | MODIFY | Add sidebar entry |
| `ui/src/routes.js` | MODIFY | Add dashboard route |

### Tests
> ❌ **No automated tests** — Manual UI testing

**Manual Test Checklist:**
- [ ] Dashboard accessible from sidebar
- [ ] Overview cards show correct data
- [ ] Requests chart displays time-series
- [ ] Latency chart shows p50/p95/p99 lines
- [ ] Top endpoints bar chart renders
- [ ] Collections table shows record counts
- [ ] Auto-refresh updates data
- [ ] Time range selector works
- [ ] Loading states display correctly

---

## PR #20: Import Backend API

**Branch:** `feat/import-api`  
**Estimated Time:** 3-4 hours  
**Dependencies:** V2 Complete (PR #17)  

### Description
Create backend API endpoints for the data import wizard.

### Tasks

- [ ] **20.1** Create `apis/import.go`:
  - [ ] `POST /api/import/preview` — Parse file, return headers + sample rows
  - [ ] `POST /api/import/validate` — Validate mapping against schema
  - [ ] `POST /api/import/execute` — Perform bulk import
- [ ] **20.2** Implement CSV parser:
  - [ ] Handle comma and tab delimiters
  - [ ] Extract headers from first row
  - [ ] Parse sample rows (first 5)
  - [ ] Count total rows
- [ ] **20.3** Implement JSON parser:
  - [ ] Handle array of objects
  - [ ] Extract keys as headers
  - [ ] Parse sample objects (first 5)
- [ ] **20.4** Implement field mapping validation:
  - [ ] Verify target collection exists
  - [ ] Verify mapped fields exist in schema
  - [ ] Check type compatibility
- [ ] **20.5** Implement bulk import execution:
  - [ ] Iterate through rows
  - [ ] Create records via PocketBase API
  - [ ] Track success/failure counts
  - [ ] Collect error details for failed rows
- [ ] **20.6** Add authentication (authenticated users)
- [ ] **20.7** Register routes in `apis/base.go`
- [ ] **20.8** Write unit tests

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `apis/import.go` | CREATE | Import API endpoints |
| `apis/base.go` | MODIFY | Register import routes |

### Tests
> ✅ **Unit Tests Required** — `apis/import_test.go`

```go
func TestImportAPI_Preview_CSV(t *testing.T)
func TestImportAPI_Preview_JSON(t *testing.T)
func TestImportAPI_Validate_Success(t *testing.T)
func TestImportAPI_Validate_InvalidField(t *testing.T)
func TestImportAPI_Execute_Success(t *testing.T)
func TestImportAPI_Execute_PartialFailure(t *testing.T)
func TestImportAPI_Unauthorized(t *testing.T)
```

---

## PR #21: Import Wizard UI

**Branch:** `feat/import-ui`  
**Estimated Time:** 4-5 hours  
**Dependencies:** PR #20  

### Description
Build the 4-step data import wizard frontend.

### Tasks

- [ ] **21.1** Create `ui/src/pages/ImportWizard.svelte`:
  - [ ] 4-step wizard layout
  - [ ] Step indicator/navigation
  - [ ] Back/Next/Import buttons
- [ ] **21.2** Create `ui/src/components/import/FileUpload.svelte`:
  - [ ] Drag-and-drop zone
  - [ ] File input fallback
  - [ ] File type validation (CSV, JSON)
  - [ ] File size display
- [ ] **21.3** Create `ui/src/components/import/DataPreview.svelte`:
  - [ ] Preview table with headers
  - [ ] First 5 rows of data
  - [ ] Total row count
- [ ] **21.4** Create `ui/src/components/import/FieldMapper.svelte`:
  - [ ] Two-column mapping UI
  - [ ] Dropdown for target field selection
  - [ ] Skip option for unwanted columns
  - [ ] Auto-detect button
  - [ ] Required field indicators
- [ ] **21.5** Create `ui/src/components/import/ImportProgress.svelte`:
  - [ ] Progress bar
  - [ ] Success/failure counts
  - [ ] Error log expandable section
  - [ ] Retry button for failed rows
- [ ] **21.6** Create `ui/src/stores/import.js`:
  - [ ] Current step
  - [ ] Selected collection
  - [ ] Uploaded file data
  - [ ] Field mapping
  - [ ] Import progress
  - [ ] Errors list
- [ ] **21.7** Modify `ui/src/App.svelte`:
  - [ ] Add Import to sidebar with upload icon
  - [ ] Add route `/import`
- [ ] **21.8** Create `ui/src/scss/_import.scss`:
  - [ ] Wizard step styles
  - [ ] Drop zone styles
  - [ ] Progress bar styles

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `ui/src/pages/ImportWizard.svelte` | CREATE | Main wizard page |
| `ui/src/components/import/FileUpload.svelte` | CREATE | File upload component |
| `ui/src/components/import/DataPreview.svelte` | CREATE | Data preview table |
| `ui/src/components/import/FieldMapper.svelte` | CREATE | Field mapping UI |
| `ui/src/components/import/ImportProgress.svelte` | CREATE | Progress display |
| `ui/src/stores/import.js` | CREATE | Import state |
| `ui/src/scss/_import.scss` | CREATE | Import styles |
| `ui/src/App.svelte` | MODIFY | Add sidebar entry |
| `ui/src/routes.js` | MODIFY | Add import route |

### Tests
> ❌ **No automated tests** — Manual UI testing

**Manual Test Checklist:**
- [ ] Import wizard accessible from sidebar
- [ ] Step 1: Can select collection from dropdown
- [ ] Step 1: Can drag-drop CSV file
- [ ] Step 1: Can drag-drop JSON file
- [ ] Step 2: Preview shows headers and data
- [ ] Step 2: Total row count is accurate
- [ ] Step 3: Auto-detect maps matching columns
- [ ] Step 3: Can manually map columns to fields
- [ ] Step 3: Can skip unwanted columns
- [ ] Step 4: Progress bar updates during import
- [ ] Step 4: Success count increments
- [ ] Step 4: Errors are logged with details
- [ ] Imported records appear in collection

---

## PR #22: V3 Documentation & Polish

**Branch:** `feat/v3-docs`  
**Estimated Time:** 2-3 hours  
**Dependencies:** PR #19, PR #21  

### Description
Complete V3 documentation and final testing.

### Tasks

- [ ] **22.1** Create `docs/DASHBOARD_FEATURE.md`:
  - [ ] Feature overview
  - [ ] Metrics explained
  - [ ] API reference
  - [ ] Screenshots
- [ ] **22.2** Create `docs/IMPORT_FEATURE.md`:
  - [ ] Feature overview
  - [ ] Supported formats (CSV, JSON)
  - [ ] Step-by-step guide
  - [ ] Troubleshooting
- [ ] **22.3** Update `README.md`:
  - [ ] Add Dashboard to features
  - [ ] Add Import to features
  - [ ] Add screenshots
- [ ] **22.4** Update `CHANGELOG.md`:
  - [ ] Add V3 release notes
- [ ] **22.5** Run full test suite
- [ ] **22.6** Build final release binary
- [ ] **22.7** End-to-end testing of all V3 features

### Files Created/Modified

| File | Action | Description |
|------|--------|-------------|
| `docs/DASHBOARD_FEATURE.md` | CREATE | Dashboard documentation |
| `docs/IMPORT_FEATURE.md` | CREATE | Import documentation |
| `README.md` | MODIFY | Add V3 features |
| `CHANGELOG.md` | MODIFY | Add V3 release notes |

### Verification
- [ ] All tests pass
- [ ] Documentation is complete
- [ ] Dashboard shows accurate metrics
- [ ] Import wizard successfully imports data
- [ ] All features accessible from sidebar

---

## Summary: Test Coverage Matrix

### V1 Tests (✅ Complete)

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

### V2 Tests (✅ Complete)

| PR | Unit Tests | Integration Tests | Manual Tests |
|----|------------|-------------------|--------------|
| PR #10: Multi-Schema | ✅ `schema_extractor_test.go` | ❌ | ❌ |
| PR #11: Dual Output | ❌ | ✅ `ai_query_test.go` | ❌ |
| PR #12: Editable UI | ❌ | ❌ | ✅ Full UI testing |
| PR #13: SQL Parser | ✅ `parser_test.go`, `mapper_test.go` | ❌ | ❌ |
| PR #14: SQL Executor | ✅ `executor_test.go` | ✅ Integration tests | ❌ |
| PR #15: SQL Terminal API | ❌ | ✅ `sql_terminal_test.go` | ✅ curl/Postman |
| PR #16: SQL Terminal UI | ❌ | ❌ | ✅ Full UI testing |
| PR #17: V2 Docs | ❌ | ✅ Full suite run | ✅ E2E testing |

### V3 Tests (🚧 Planned)

| PR | Unit Tests | Integration Tests | Manual Tests |
|----|------------|-------------------|--------------|
| PR #18: Metrics API | ✅ `metrics_test.go` | ✅ API tests | ✅ curl/Postman |
| PR #19: Dashboard UI | ❌ | ❌ | ✅ Full UI testing |
| PR #20: Import API | ✅ `import_test.go` | ✅ API tests | ✅ curl/Postman |
| PR #21: Import UI | ❌ | ❌ | ✅ Full UI testing |
| PR #22: V3 Docs | ❌ | ✅ Full suite run | ✅ E2E testing |

---

## Quick Reference: All Files

### V1 Files (✅ Complete)

```
NEW FILES (18):
├── apis/ai_query.go
├── apis/ai_query_test.go
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
├── services/ai/filter_tokenizer.go
├── services/ai/errors.go
├── ui/src/stores/ai.js
├── ui/src/components/ai/AIQueryInput.svelte
├── ui/src/components/ai/AIFilterDisplay.svelte
├── ui/src/components/ai/AIQueryResults.svelte
├── ui/src/components/ai/AIQueryPanel.svelte
├── ui/src/components/ai/AISettingsForm.svelte
├── ui/src/pages/settings/AI.svelte
├── ui/src/scss/_ai.scss
├── docs/AI_QUERY_FEATURE.md

MODIFIED FILES (4):
├── core/settings.go
├── apis/base.go
├── ui/src/App.svelte
├── README.md
```

### V2 Files (✅ Complete)

```
NEW FILES (14):
├── apis/sql_terminal.go
├── services/sql/parser.go
├── services/sql/executor.go
├── services/sql/mapper.go
├── ui/src/stores/sql.js
├── ui/src/components/ai/QueryTabs.svelte
├── ui/src/components/ai/EditableCodeBlock.svelte
├── ui/src/components/sql/SQLEditor.svelte
├── ui/src/components/sql/SchemaExplorer.svelte
├── ui/src/components/sql/ResultsTable.svelte
├── ui/src/components/sql/QueryHistory.svelte
├── ui/src/pages/SQLTerminal.svelte
├── ui/src/scss/_sql.scss
├── docs/SQL_TERMINAL_FEATURE.md

MODIFIED FILES (6):
├── services/ai/schema_extractor.go
├── services/ai/prompt_template.go
├── apis/ai_query.go
├── apis/base.go
├── ui/src/App.svelte
├── ui/src/components/ai/AIQueryPanel.svelte
```

### V3 Files (🚧 Planned)

```
NEW FILES (16):
├── apis/metrics.go
├── apis/metrics_test.go
├── apis/import.go
├── apis/import_test.go
├── ui/src/pages/Dashboard.svelte
├── ui/src/pages/ImportWizard.svelte
├── ui/src/components/dashboard/MetricCard.svelte
├── ui/src/components/dashboard/RequestsChart.svelte
├── ui/src/components/dashboard/LatencyChart.svelte
├── ui/src/components/dashboard/EndpointsChart.svelte
├── ui/src/components/dashboard/CollectionsTable.svelte
├── ui/src/components/import/FileUpload.svelte
├── ui/src/components/import/DataPreview.svelte
├── ui/src/components/import/FieldMapper.svelte
├── ui/src/components/import/ImportProgress.svelte
├── ui/src/stores/dashboard.js
├── ui/src/stores/import.js
├── ui/src/scss/_dashboard.scss
├── ui/src/scss/_import.scss
├── docs/DASHBOARD_FEATURE.md
├── docs/IMPORT_FEATURE.md

MODIFIED FILES (4):
├── apis/base.go
├── ui/src/App.svelte
├── ui/src/routes.js
├── README.md
```

---

## Execution Order

### V1 Execution (✅ Complete)

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

### V2 Execution (✅ Complete)

```
                    PR #10 (Multi-Schema)
                    ↓              ↓
        PR #11 (Dual Output)   PR #13 (SQL Parser)
                    ↓              ↓
        PR #12 (Editable UI)   PR #14 (SQL Executor)
                    ↓              ↓
                    ↓          PR #15 (SQL Terminal API)
                    ↓              ↓
                    ↓          PR #16 (SQL Terminal UI)
                    ↓              ↓
                    └──────────────┴──→ PR #17 (V2 Docs)
```

### V3 Execution (🚧 In Progress)

```
        PR #18 (Metrics API)        PR #20 (Import API)
                ↓                           ↓
        PR #19 (Dashboard UI)       PR #21 (Import UI)
                ↓                           ↓
                └───────────────────────────┴──→ PR #22 (V3 Docs)
```

---

## Time Estimates

### V1 (✅ Complete)

| Phase | Hours |
|-------|-------|
| Setup (PR #1) | 2-3 |
| Backend (PRs #2-6) | 25-30 |
| Frontend (PRs #7-8) | 10-13 |
| Docs (PR #9) | 3-4 |
| **V1 Total** | **38 hours** |

### V2 (✅ Complete)

| Phase | Hours |
|-------|-------|
| Enhanced AI Query (PRs #10-12) | 13-16 |
| SQL Terminal Backend (PRs #13-15) | 17-20 |
| SQL Terminal UI (PR #16) | 8-10 |
| Documentation (PR #17) | 3-4 |
| **V2 Total** | **45 hours** |

### V3 (🚧 In Progress)

| Phase | Hours |
|-------|-------|
| Dashboard Backend (PR #18) | 3-4 |
| Dashboard UI (PR #19) | 4-5 |
| Import Backend (PR #20) | 3-4 |
| Import UI (PR #21) | 4-5 |
| Documentation (PR #22) | 2-3 |
| **V3 Total** | **17-23 hours** |

### Grand Total

| Version | Status | Hours |
|---------|--------|-------|
| V1 | ✅ Complete | 38 |
| V2 | ✅ Complete | 45 |
| V3 | 🚧 In Progress | 17-23 |
| **Total** | | **100-106 hours** |

---

**Document Status:** V1 Complete, V2 Complete, V3 In Progress  
**Total PRs:** 22 (V1: 9 ✅, V2: 8 ✅, V3: 5 🚧)  
**Total Test Files:** 16
