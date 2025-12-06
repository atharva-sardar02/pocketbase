# System Patterns: PocketBase AI Query Assistant

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    PocketBase Admin UI                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  AI Query Sidebar Panel (Svelte Components)          │  │
│  │  - AIQueryInput.svelte                                │  │
│  │  - AIFilterDisplay.svelte                             │  │
│  │  - AIQueryResults.svelte                              │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP POST /api/ai/query
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    PocketBase Go Backend                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  API Endpoint: apis/ai_query.go                      │  │
│  │  - Authentication check                              │  │
│  │  - Request validation                                │  │
│  │  - Collection schema loading                         │  │
│  └────────────────────────┬──────────────────────────────┘  │
│                           │                                  │
│  ┌────────────────────────▼──────────────────────────────┐  │
│  │  Services Layer: services/ai/                        │  │
│  │  - schema_extractor.go  → Extract collection schema  │  │
│  │  - prompt_builder.go    → Build LLM prompts          │  │
│  │  - openai_client.go     → Call LLM API               │  │
│  │  - filter_validator.go  → Validate generated filter  │  │
│  └────────────────────────┬──────────────────────────────┘  │
│                           │                                  │
│  ┌────────────────────────▼──────────────────────────────┐  │
│  │  Core Layer: core/                                   │  │
│  │  - ai_settings.go     → Settings struct & validation │  │
│  │  - settings.go        → Main settings integration    │  │
│  └────────────────────────┬──────────────────────────────┘  │
│                           │                                  │
│  ┌────────────────────────▼──────────────────────────────┐  │
│  │  Data Layer: SQLite (_params table)                   │  │
│  │  - AI settings storage (encrypted API keys)           │  │
│  └───────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              External LLM Provider (OpenAI/Ollama)          │
│  - OpenAI API (default: gpt-4o-mini)                        │
│  - Ollama (local deployment option)                         │
└─────────────────────────────────────────────────────────────┘
```

## Key Technical Decisions

### 1. Service Layer Pattern
**Pattern:** Separate service layer for AI functionality  
**Location:** `services/ai/`  
**Rationale:** Keeps AI logic isolated, testable, and reusable

**Components:**
- `schema_extractor.go` - Collection schema → prompt-friendly format
- `prompt_builder.go` - Constructs system and user prompts
- `openai_client.go` - HTTP client for LLM API calls
- `filter_validator.go` - Validates LLM-generated filters
- `filter_tokenizer.go` - Basic filter parsing for validation

### 2. Settings Integration Pattern
**Pattern:** Extend existing PocketBase settings system  
**Location:** `core/ai_settings.go` + `core/settings.go`  
**Storage:** `_params` table (existing)  
**Rationale:** No schema changes needed, follows existing patterns

**Structure:**
```go
type AISettings struct {
    Enabled     bool
    Provider    string  // "openai", "ollama", "anthropic"
    BaseURL     string
    APIKey      string  // encrypted
    Model       string
    Temperature float64
}
```

### 3. API Endpoint Pattern
**Pattern:** RESTful endpoint following PocketBase conventions  
**Location:** `apis/ai_query.go`  
**Route:** `POST /api/ai/query`  
**Authentication:** Uses existing PocketBase auth system  
**Security:** Respects collection API rules

**Request Flow:**
1. Authenticate user
2. Validate request (collection, query required)
3. Load AI settings (check if enabled)
4. Load collection schema
5. Extract schema → Build prompts → Call LLM
6. Validate generated filter
7. Optionally execute filter
8. Return response

### 4. Validation Layer Pattern
**Pattern:** Multi-layer validation before execution  
**Rationale:** Prevent hallucinated field names, invalid syntax, security issues

**Validation Steps:**
1. **LLM Output Validation:** Parse filter expression
2. **Field Existence Check:** Verify all fields exist in collection
3. **Operator Compatibility:** Check operators match field types
4. **Syntax Validation:** Ensure valid PocketBase filter syntax
5. **Security Check:** Respect collection API rules (listRule)

### 5. Error Handling Pattern
**Pattern:** Custom error types with user-friendly messages  
**Location:** `services/ai/errors.go`

**Error Types:**
- `AIClientError` - Base error for AI client issues
- `AIRateLimitError` - 429 responses (rate limited)
- `AIAuthError` - 401 responses (invalid API key)
- `AITimeoutError` - Context deadline exceeded
- `ValidationError` - Invalid filter generated

**Error Flow:**
- Backend errors → Structured JSON response with error message
- Frontend errors → Display in UI with helpful guidance

## Component Relationships

### Backend Components

```
apis/ai_query.go
    ├── Uses: core/ai_settings.go (load settings)
    ├── Uses: services/ai/schema_extractor.go (get schema)
    ├── Uses: services/ai/prompt_builder.go (build prompts)
    ├── Uses: services/ai/openai_client.go (call LLM)
    ├── Uses: services/ai/filter_validator.go (validate filter)
    └── Returns: JSON response with filter + results
```

### Frontend Components

```
ui/src/App.svelte
    └── Routes to: ui/src/components/ai/AIQueryPanel.svelte
            ├── Uses: ui/src/components/ai/AIQueryInput.svelte
            ├── Uses: ui/src/components/ai/AIFilterDisplay.svelte
            ├── Uses: ui/src/components/ai/AIQueryResults.svelte
            └── Uses: ui/src/stores/ai.js (state management)
```

### Settings Components

```
ui/src/pages/settings/AI.svelte
    └── Uses: ui/src/components/ai/AISettingsForm.svelte
            └── Calls: apis/ai_settings.go (save/load settings)
```

## Design Patterns in Use

### 1. Dependency Injection
- Services receive dependencies (settings, app instance) via constructor
- Enables testing with mocks

### 2. Context Pattern (Go)
- All LLM calls use `context.WithTimeout()` for cancellation
- Default timeout: 30 seconds

### 3. Builder Pattern
- `prompt_builder.go` constructs prompts step-by-step
- Allows flexible prompt customization

### 4. Validator Pattern
- Separate validation layer before execution
- Returns detailed error messages

### 5. Store Pattern (Svelte)
- Centralized state management in `ui/src/stores/ai.js`
- Components subscribe to store updates

## Integration Points with PocketBase

### 1. Settings System
- Extends `core/settings.go` with `AI AISettings` field
- Uses existing `_params` table for storage
- Follows existing encryption patterns for API keys

### 2. Authentication System
- Uses existing PocketBase auth middleware
- Respects user roles and permissions

### 3. Collection API Rules
- Validates against collection's `listRule` before execution
- Only returns data user has access to

### 4. Router System
- Registers route in `apis/base.go` using PocketBase's router
- Follows existing API endpoint patterns

### 5. Admin UI Framework
- Uses Svelte 4 (existing framework)
- Follows PocketBase component patterns
- Uses existing `ApiClient` utility

## Security Patterns

### 1. Input Sanitization
- User queries are sanitized before sending to LLM
- Prevents prompt injection attacks

### 2. Output Validation
- All LLM outputs are validated before execution
- Prevents code injection via hallucinated filters

### 3. API Key Encryption
- API keys stored encrypted at rest
- Follows existing OAuth secret encryption pattern

### 4. Access Control
- Feature can be enabled/disabled globally
- Respects collection API rules
- Requires authentication for API access

## Testing Patterns

### Unit Tests
- Each service has corresponding `*_test.go` file
- Mock HTTP servers for LLM client tests
- Table-driven tests for validation logic

### Integration Tests
- Full API endpoint tests in `apis/ai_query_test.go`
- Uses PocketBase test utilities (`tests.TestApp`)

### Manual Tests
- UI components tested manually (Svelte testing is limited)
- End-to-end workflow verification



