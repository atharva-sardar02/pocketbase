# Technical Context: PocketBase AI Query Assistant

## Technology Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** PocketBase (forked repository)
- **Database:** SQLite (existing PocketBase database)
- **HTTP Client:** Standard `net/http` package (no external LLM libraries)
- **Build:** `go build` with `CGO_ENABLED=0` for static binary

### Frontend
- **Framework:** Svelte 4
- **Build Tool:** Vite
- **State Management:** Svelte stores
- **Styling:** SCSS (following PocketBase patterns)
- **Build:** `npm run build` outputs to `/ui/dist`

### External Services
- **LLM Provider:** OpenAI API (default: `gpt-4o-mini`)
- **Alternative:** Ollama (local deployment)
- **API Format:** OpenAI-compatible REST API

## Development Environment

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- Git (for forking/cloning repository)

### Setup Commands
```powershell
# Frontend setup
cd ui
npm install
npm run build

# Backend build
cd ../examples/base
$env:CGO_ENABLED="0"
go build
.\base.exe serve
```

### Development Workflow
1. Make changes to Go backend code
2. Make changes to Svelte frontend code
3. Run `npm run build` in `/ui` directory
4. Rebuild Go binary (includes embedded UI)
5. Test in browser at `http://127.0.0.1:8090/_/`

## Project Structure

### Key Directories
```
pocketbase/
├── apis/                    # API endpoint handlers
│   └── ai_query.go         # NEW: AI query endpoint
├── core/                    # Core functionality
│   ├── ai_settings.go      # NEW: AI settings struct
│   └── settings.go          # MODIFY: Add AI settings
├── services/                # Service layer
│   └── ai/                  # NEW: AI services directory
│       ├── openai_client.go
│       ├── prompt_builder.go
│       ├── schema_extractor.go
│       └── filter_validator.go
├── ui/                      # Frontend (Svelte)
│   └── src/
│       ├── components/ai/   # NEW: AI components
│       ├── pages/settings/  # NEW: AI settings page
│       └── stores/          # NEW: AI state store
└── examples/base/           # Entry point
    └── main.go
```

## Technical Constraints

### Go Backend Constraints
1. **No External LLM Libraries:** Use raw `net/http` for full control
2. **Context Timeouts:** All LLM calls must use `context.WithTimeout()` (30s default)
3. **Error Handling:** Go's explicit error handling required
4. **Type Safety:** Strict typing - no dynamic types
5. **Dependency Management:** Use `go mod tidy` after adding imports

### Svelte Frontend Constraints
1. **Build Process:** Must rebuild UI after every change, then rebuild Go binary
2. **Component Patterns:** Follow existing PocketBase component structures
3. **State Management:** Use Svelte stores (existing pattern)
4. **API Client:** Use existing `ApiClient` utility in `/ui/src/utils/ApiClient.js`
5. **Styling:** Follow PocketBase SCSS patterns

### PocketBase Integration Constraints
1. **Hook System:** Must understand `OnServe()` and router binding
2. **Settings Storage:** Use existing `_params` table (no schema changes)
3. **Authentication:** Use existing PocketBase auth system
4. **API Rules:** Must respect collection `listRule` for security
5. **Router:** Register routes in `apis/base.go` following existing patterns

## Dependencies

### Backend Dependencies
- **Standard Library Only:** No external Go packages for LLM communication
- **PocketBase Core:** Uses existing PocketBase packages
- **Testing:** `testing` package + PocketBase test utilities

### Frontend Dependencies
- **Svelte 4:** Existing framework
- **Vite:** Build tool (existing)
- **No New Dependencies:** Use existing UI libraries

## API Specifications

### AI Query Endpoint
**Route:** `POST /api/ai/query`  
**Authentication:** Required (Bearer token)  
**Content-Type:** `application/json`

**Request:**
```json
{
  "collection": "orders",
  "query": "pending orders over $100 from last week",
  "execute": true,
  "page": 1,
  "perPage": 20
}
```

**Response (Success):**
```json
{
  "filter": "status='pending' && total>100 && created>=@now-604800",
  "results": [...],
  "totalItems": 42,
  "page": 1,
  "perPage": 20
}
```

**Response (Error):**
```json
{
  "error": "Unknown field: invalid_field. Available fields: status, total, created"
}
```

### AI Settings Endpoint
**Route:** `GET/POST /api/settings/ai` (to be implemented)  
**Authentication:** Superuser only  
**Purpose:** Load and save AI configuration

## LLM Integration Details

### OpenAI API Format
**Endpoint:** `POST {baseUrl}/chat/completions`  
**Headers:**
```
Authorization: Bearer {apiKey}
Content-Type: application/json
```

**Request Body:**
```json
{
  "model": "gpt-4o-mini",
  "messages": [
    {"role": "system", "content": "{systemPrompt}"},
    {"role": "user", "content": "{userQuery}"}
  ],
  "temperature": 0.1
}
```

**Response:**
```json
{
  "choices": [{
    "message": {
      "content": "status = \"active\""
    }
  }]
}
```

### Ollama Compatibility
- Uses same OpenAI-compatible API format
- Base URL: `http://localhost:11434/v1`
- Model names: `llama2`, `mistral`, etc.

## PocketBase Filter Syntax

### Operators
- `=` - Equals
- `!=` - Not equals
- `>` `<` `>=` `<=` - Comparison
- `~` - Contains (LIKE)
- `!~` - Not contains
- `?=` - Any equals (arrays)
- `?~` - Any contains (arrays)

### Logical Operators
- `&&` - AND
- `||` - OR
- `()` - Grouping

### Datetime Macros
- `@now` - Current datetime
- `@second`, `@minute`, `@hour`, `@weekday`, `@day`, `@month`, `@year`
- Arithmetic: `@now - 604800` (subtract seconds)

### Example Translations
| Natural Language | PocketBase Filter |
|-----------------|-------------------|
| "active users" | `status = "active"` |
| "orders over $100" | `total > 100` |
| "posts from last week" | `created >= @now - 604800` |
| "titles containing 'hello'" | `title ~ "hello"` |

## Testing Strategy

### Unit Tests
- **Location:** `*_test.go` files alongside source
- **Framework:** Go `testing` package
- **Mocking:** `httptest.NewServer` for HTTP mocks
- **Coverage Target:** >80% for new code

### Integration Tests
- **Location:** `apis/ai_query_test.go`
- **Framework:** PocketBase test utilities (`tests.TestApp`)
- **Scope:** Full API endpoint testing

### Manual Tests
- **UI Components:** Manual browser testing
- **End-to-End:** Full workflow verification
- **Checklist:** Documented in task list

## Build & Deployment

### Development Build
```powershell
# Frontend
cd ui
npm run build

# Backend
cd ../examples/base
$env:CGO_ENABLED="0"
go build
```

### Production Build
- Static binary with embedded UI
- No external dependencies
- Single executable file

## Performance Considerations

### LLM Latency
- **Expected:** 0.5-2 seconds per query (OpenAI `gpt-4o-mini`)
- **Mitigation:** Show loading spinner, set 30s timeout

### Cost Management
- **Default Model:** `gpt-4o-mini` (~$0.00015 per query)
- **Settings:** Admin can choose model based on cost/quality tradeoff
- **Future:** Rate limiting, usage quotas (V2)

## Security Considerations

### Input Sanitization
- Sanitize user queries before sending to LLM
- Prevent prompt injection attacks

### Output Validation
- Validate all LLM outputs before execution
- Prevent code injection via hallucinated filters

### API Key Security
- Encrypt API keys at rest
- Never log API keys
- Follow existing OAuth secret encryption pattern

### Access Control
- Feature can be enabled/disabled globally
- Respects collection API rules
- Requires authentication for API access

## Known Technical Challenges

### Challenge 1: Go Learning Curve
**Risk:** Go's strict typing and error handling  
**Mitigation:** Budget 1-2 days for familiarization, use AI assistance

### Challenge 2: PocketBase Hook System
**Risk:** Must understand `OnServe()` and router binding  
**Mitigation:** Study `/apis/record.go` as reference

### Challenge 3: Svelte Build Process
**Risk:** Must rebuild UI after every change  
**Mitigation:** Create script to automate build process

### Challenge 4: LLM Hallucination
**Risk:** LLM invents field names or syntax  
**Mitigation:** Comprehensive validation layer

### Challenge 5: Prompt Injection
**Risk:** User could manipulate LLM via crafted queries  
**Mitigation:** Input sanitization, output validation, never execute raw output



