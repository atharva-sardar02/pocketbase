# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **fork of PocketBase** with an added **AI Query Assistant** feature. PocketBase is an open-source Go backend with embedded SQLite, realtime subscriptions, file/user management, and an admin dashboard. The AI extension enables natural language queries to database collections using LLM providers (OpenAI, Ollama, Anthropic, custom).

**Base Project**: https://github.com/pocketbase/pocketbase
**Go Version**: 1.24.0
**UI Framework**: Svelte 4 + Vite
**Database**: SQLite (via modernc.org/sqlite)

## Build & Development Commands

### Backend Development

```powershell
# Run development server from examples/base directory
cd examples\base
go run main.go serve

# Build static binary (Windows)
cd examples\base
$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build

# Build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build
```

### Frontend Development

```powershell
# Install dependencies
cd ui
npm install

# Start dev server (http://localhost:3000)
npm run dev

# Build for production (outputs to ui/dist)
npm run build
```

**Important**: After UI changes, you must:
1. Run `npm run build` in `/ui`
2. Rebuild the Go binary in `/examples/base`

The backend serves the pre-built UI from `ui/dist` at `http://localhost:8090/_/`

### Testing & Quality

```bash
# Run all tests
make test
# or
go test ./... -v --cover

# Run linter (golangci-lint required)
make lint
# or
golangci-lint run -c ./golangci.yml ./...

# Generate test coverage report
make test-report
```

## Architecture Overview

### Directory Structure

```
pocket-base-ai/
├── apis/                    # REST API endpoint handlers (one file per feature)
│   ├── base.go              # Router setup - binds all API routes
│   ├── ai_query.go          # AI Query endpoint (POST /api/ai/query)
│   ├── settings.go          # Settings endpoints (includes AI settings)
│   ├── record_crud.go       # Record CRUD operations
│   └── ...                  # Other API handlers (auth, collections, logs, etc.)
│
├── core/                    # Core application logic and interfaces
│   ├── app.go               # App interface definition
│   ├── base.go              # Base app implementation
│   ├── ai_settings.go       # AI configuration struct with validation
│   └── ...                  # Collection, Field, Record models
│
├── services/
│   └── ai/                  # AI Query service implementations
│       ├── openai_client.go      # LLM API client (OpenAI-compatible)
│       ├── schema_extractor.go   # Collection schema extraction
│       ├── prompt_builder.go     # System/user prompt construction
│       ├── prompt_template.go    # Pre-defined system prompt
│       ├── filter_validator.go   # Validates generated filters
│       └── filter_tokenizer.go   # Parses filter expressions
│
├── tools/                   # Utility packages
│   ├── router/              # HTTP routing (chi-based)
│   ├── auth/                # OAuth2 provider integrations
│   ├── filesystem/          # File storage (local or S3)
│   ├── mailer/              # Email sending
│   ├── hook/                # Event hooks system
│   └── ...                  # logger, security, subscriptions, etc.
│
├── ui/src/                  # Frontend (Svelte)
│   ├── components/
│   │   ├── ai/              # AI Query components
│   │   │   ├── AIQueryPanel.svelte       # Main AI query interface
│   │   │   ├── AIQueryInput.svelte       # Natural language input
│   │   │   ├── AIFilterDisplay.svelte    # Filter display with copy
│   │   │   ├── AIQueryResults.svelte     # Results display
│   │   │   └── AISettingsForm.svelte     # Settings configuration
│   │   ├── base/            # Core UI components
│   │   └── ...              # records, collections, settings, logs
│   ├── pages/
│   │   └── settings/
│   │       └── AI.svelte    # AI settings page
│   ├── stores/
│   │   ├── ai.js            # AI state management
│   │   └── ...              # Other Svelte stores
│   └── routes.js            # Route definitions
│
├── examples/base/           # Standalone executable entry point
│   └── main.go              # Main program with plugins (jsvm, migrations, etc.)
│
├── plugins/                 # Plugin system
│   ├── jsvm/                # JavaScript VM for hooks
│   ├── migratecmd/          # Migration commands
│   └── ghupdate/            # GitHub self-update
│
└── docs/
    └── AI_QUERY_FEATURE.md  # Comprehensive AI feature documentation
```

### Key Architectural Patterns

**Go Backend Patterns:**
- **Interface-based design**: `core.App` interface allows extensibility
- **Hook system**: Events use `hook.Handler` with priority for extensibility (see `OnServe()`)
- **One file per API feature**: Each API endpoint in separate file (e.g., `ai_query.go`, `record_crud.go`)
- **Validation**: Uses `go-ozzo/validation` package extensively
- **Transaction support**: `TxInfo()` for database transactions

**API Patterns:**
- RESTful endpoints grouped by feature
- Authentication via PocketBase's existing auth system
- Role-based access (superuser checks via `re.HasSuperuserAuth()`)
- Consistent error handling with `apis.NewApiError()`
- JSON request/response bodies

**Frontend Patterns:**
- Svelte for reactive components
- Stores for centralized state management (`/ui/src/stores/`)
- Component composition (PageWrapper, sidebar patterns)
- Router-based navigation with access conditions
- API communication via `ApiClient` utility (`/ui/src/utils/ApiClient.js`)

## AI Query Feature Integration

### Request Flow

1. User enters natural language query in `AIQueryPanel.svelte`
2. Frontend sends `POST /api/ai/query` with: `{collection, query, execute, page, perPage}`
3. Backend handler in `apis/ai_query.go`:
   - Validates authentication and collection access
   - Extracts schema via `services/ai/schema_extractor.go`
   - Builds prompts via `services/ai/prompt_builder.go`
   - Sends to LLM via `services/ai/openai_client.go`
   - Validates filter via `services/ai/filter_validator.go`
   - Optionally executes filter and returns results
4. Frontend displays filter expression and results

### Settings Storage

AI settings are stored using PocketBase's existing settings mechanism:
- Stored in database via `core.AISettings` struct
- Accessed through `app.Settings().AI`
- Encrypted API keys using existing encryption patterns
- Configured in Settings → AI Query page

### LLM Provider Support

All providers use OpenAI-compatible API format:
- **OpenAI**: `https://api.openai.com/v1` (default: `gpt-4o-mini`)
- **Ollama**: `http://localhost:11434/v1` (local, no API key)
- **Anthropic**: `https://api.anthropic.com/v1`
- **Custom**: Any OpenAI-compatible endpoint

## Common Development Tasks

### Adding a New API Endpoint

1. Create new file in `/apis/` (e.g., `my_feature.go`)
2. Implement handler function that accepts `*core.RequestEvent`
3. Register route in `apis/base.go` within `BindAppRoutes()`
4. Add tests in `apis/my_feature_test.go`

Example:
```go
// apis/my_feature.go
func myFeatureHandler(re *core.RequestEvent) error {
    // Validate input
    data := struct {
        Field string `json:"field"`
    }{}
    if err := re.BindBody(&data); err != nil {
        return apis.NewBadRequestError("Invalid request", err)
    }

    // Business logic
    result := doSomething(data.Field)

    // Return response
    return re.JSON(200, result)
}

// apis/base.go - in BindAppRoutes()
api.POST("/my-feature", myFeatureHandler)
```

### Modifying UI Components

1. Edit Svelte components in `/ui/src/components/`
2. Test changes with `npm run dev` (runs on http://localhost:3000)
3. Build UI: `npm run build`
4. Rebuild backend to embed new UI

### Extending AI Query Service

Key files to modify:
- **Prompt Template**: `services/ai/prompt_template.go` - System prompt with filter syntax rules
- **Schema Extraction**: `services/ai/schema_extractor.go` - How collection schema is formatted
- **Validation**: `services/ai/filter_validator.go` - Filter validation logic
- **LLM Client**: `services/ai/openai_client.go` - LLM API communication

### Running Tests for AI Features

```bash
# Test AI query endpoint
go test ./apis -run TestAIQuery -v

# Test AI services
go test ./services/ai/... -v

# Test specific validation
go test ./services/ai -run TestFilterValidator -v
```

## PocketBase-Specific Conventions

### Filter Syntax

PocketBase uses a custom filter syntax (NOT SQL):
- Operators: `=`, `!=`, `>`, `<`, `>=`, `<=`, `~` (contains), `!~` (not contains)
- Logical: `&&` (AND), `||` (OR), `()` (grouping)
- Arrays: `?=` (any equals), `?~` (any contains)
- Datetime macros: `@now`, `@today`, `@month`, `@year`, etc.
- Example: `status = "active" && created >= @now - 604800`

### Settings Management

Settings are stored in `_params` table and accessed via:
```go
app.Settings().AI.Enabled  // Access AI settings
app.Save(settings)         // Save settings
```

### Hook System Usage

PocketBase uses hooks for extensibility:
```go
app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
    Func: func(e *core.ServeEvent) error {
        // Register routes
        e.Router.GET("/my-route", handler)
        return e.Next()  // Continue hook chain
    },
    Priority: 10,  // Lower = earlier execution
})
```

### Database Queries

Use PocketBase's DAO (Data Access Object):
```go
// Find records with filter
records, err := app.FindRecordsByFilter(
    "collectionName",
    "status = 'active'",
    "-created",  // Sort descending by created
    30,          // Limit
    0,           // Offset
)

// Find single record
record, err := app.FindFirstRecordByFilter("users", "email = {:email}", dbx.Params{"email": email})
```

## Important Notes

### Security Considerations

- AI Query respects existing collection `listRule` permissions
- Filter validation prevents SQL injection
- API keys stored encrypted in database
- All API endpoints require authentication
- Superuser checks for admin operations

### Testing Philosophy

- Mix of unit and integration tests
- Use standard Go `testing` package
- Test files alongside implementation files (`*_test.go`)
- Run full suite with `go test ./...`

### Build Targets

PocketBase supports limited platforms due to pure Go SQLite driver:
- darwin (amd64, arm64)
- linux (386, amd64, arm, arm64, loong64, ppc64le, riscv64, s390x)
- windows (386, amd64, arm64)
- freebsd (amd64, arm64)

### Configuration Files

- `go.mod` - Go dependencies
- `ui/package.json` - Frontend dependencies
- `Makefile` - Common commands (lint, test, test-report)
- `golangci.yml` - Linter configuration
- `.env.development.local` - Local UI development config (not in repo)

## Additional Resources

- **AI Feature Documentation**: `/docs/AI_QUERY_FEATURE.md`
- **Product Requirements**: `/PocketBase_AI_Query_Assistant_PRD.md`
- **PocketBase Docs**: https://pocketbase.io/docs
- **Filter Syntax Reference**: https://pocketbase.io/docs/api-records/#filtering
- **Contributing Guide**: `/CONTRIBUTING.md`
