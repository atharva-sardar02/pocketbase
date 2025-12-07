# AI Query Assistant Feature

## Overview

The AI Query Assistant is a feature extension to PocketBase that enables users to query database collections using natural language instead of learning PocketBase's filter syntax. This feature makes data exploration accessible to non-technical users while providing power users with a faster way to construct complex queries.

## Features

### Core Features (V1)
- **Natural Language Query Interface** - Query collections using plain English in the Admin UI
- **Schema-Aware Prompting** - Automatically understands your collection's field names and types
- **Filter Expression Display** - Shows generated PocketBase filter syntax with copy functionality
- **API Endpoint** - Programmatic access for building AI-powered search into applications
- **Multiple LLM Provider Support** - Configure OpenAI, Ollama, Anthropic, or custom providers
- **Security Integration** - Respects existing collection API rules and authentication
- **Validation & Error Handling** - Validates generated filters and provides helpful error messages

### V2 Enhancements
- **Dual Output Mode** - Get both PocketBase filter AND SQL for any query
- **Editable Code Blocks** - Edit generated filter/SQL before executing
- **Tab Interface** - Switch between Filter and SQL views
- **Multi-Collection Schema** - Supports queries across related collections
- **SQL Terminal** - Full SQL interface (see [SQL_TERMINAL_FEATURE.md](./SQL_TERMINAL_FEATURE.md))
- **Auto SQL Detection** - Identifies when SQL is required (JOINs, aggregations)

## Setup Instructions

### Prerequisites

- PocketBase instance (this fork)
- Go 1.21+ (for building from source)
- Node.js 18+ (for building UI)
- LLM Provider API key (OpenAI, Anthropic, or access to Ollama instance)

### Installation

1. **Build the application** (if building from source):
   ```powershell
   # Build UI
   cd ui
   npm install
   npm run build

   # Build backend
   cd ../examples/base
   $env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build
   ```

2. **Start PocketBase**:
   ```powershell
   .\base.exe serve
   ```

3. **Access Admin UI**:
   - Navigate to `http://127.0.0.1:8090/_/`
   - Log in as superuser

## Configuration Guide

### Enable AI Query Feature

1. Navigate to **Settings** → **AI Query** in the Admin UI
2. Toggle **Enable AI Query** to ON
3. Configure your LLM provider:

#### OpenAI Configuration

- **Provider**: OpenAI
- **API Base URL**: `https://api.openai.com/v1` (auto-filled)
- **API Key**: Your OpenAI API key (starts with `sk-`)
- **Model**: Choose from:
  - `gpt-4o-mini` (recommended, cost-effective)
  - `gpt-4o` (more capable)
  - `gpt-3.5-turbo` (legacy)
- **Temperature**: `0.1` (recommended for consistent filter generation)

#### Ollama Configuration (Local)

- **Provider**: Ollama
- **API Base URL**: `http://localhost:11434/v1` (auto-filled)
- **API Key**: Not required (leave empty)
- **Model**: Choose from:
  - `llama2`
  - `llama3`
  - `mistral`
- **Temperature**: `0.1`

#### Anthropic Configuration

- **Provider**: Anthropic
- **API Base URL**: `https://api.anthropic.com/v1` (auto-filled)
- **API Key**: Your Anthropic API key
- **Model**: Choose from:
  - `claude-3-5-sonnet-20241022`
  - `claude-3-opus-20240229`
- **Temperature**: `0.1`

#### Custom Provider Configuration

- **Provider**: Custom
- **API Base URL**: Your custom API endpoint (must be OpenAI-compatible)
- **API Key**: Your API key (if required)
- **Model**: Enter model name manually
- **Temperature**: `0.1`

### Test Connection

After configuring your settings, click **Test Connection** to verify:
- API endpoint is reachable
- API key is valid
- Model is available

### Save Settings

Click **Save** to persist your configuration. Settings are encrypted and stored in PocketBase's `_params` table.

## Usage Guide

### Admin UI - Natural Language Queries

1. **Access AI Query Panel**:
   - Click **AI Query** in the main sidebar (robot icon)
   - Or navigate to `/ai-query` in the Admin UI

2. **Select Collection**:
   - Choose the collection you want to query from the dropdown

3. **Enter Natural Language Query**:
   - Type your query in plain English
   - Examples:
     - "show me all orders from last week that are still pending"
     - "find users created in the last 30 days"
     - "products with price greater than 100 and status is active"
     - "recent posts with more than 10 likes"

4. **Execute Query**:
   - Click **Search** button or press `Ctrl+Enter` (Windows) / `Cmd+Enter` (Mac)
   - The AI will generate a PocketBase filter expression
   - Results will be displayed (if `execute` is enabled)

5. **Use Generated Filter**:
   - **Copy Filter**: Click the copy button to copy the filter expression
   - **Apply Filter**: Click "Apply Filter" to navigate to the collection with the filter pre-applied
   - **View Results**: Review the matching records in the results panel

### API Usage

#### Endpoint

```
POST /api/ai/query
```

#### Authentication

Requires authentication via Bearer token (same as other PocketBase API endpoints).

#### Request Body

```json
{
  "collection": "orders",
  "query": "show me all pending orders from last week",
  "execute": true,
  "page": 1,
  "perPage": 30,
  "mode": "dual"
}
```

**Fields:**
- `collection` (required): Collection name or ID
- `query` (required): Natural language query string
- `execute` (optional): Whether to execute the filter and return results (default: `false`)
- `page` (optional): Page number for pagination (default: `1`)
- `perPage` (optional): Records per page (default: `30`)
- `mode` (optional, V2): Query mode - `"filter"`, `"dual"`, or `"sql"` (default: `"filter"`)

#### Response

**Success Response (200 OK) - Filter Mode:**

```json
{
  "filter": "status = \"pending\" && created >= @now - 604800",
  "results": [
    {
      "id": "abc123",
      "status": "pending",
      "total": 150.00,
      "created": "2025-12-01 10:00:00.000Z"
    }
  ],
  "totalItems": 42,
  "page": 1,
  "perPage": 30
}
```

**Success Response (200 OK) - Dual Mode (V2):**

```json
{
  "filter": "status = \"pending\" && created >= @now - 604800",
  "sql": "SELECT * FROM orders WHERE status = 'pending' AND created >= datetime('now', '-7 days')",
  "requiresSQL": false,
  "canUseFilter": true,
  "results": [...],
  "totalItems": 42,
  "page": 1,
  "perPage": 30
}
```

**V2 Response Fields:**
- `sql`: Generated SQL query equivalent
- `requiresSQL`: `true` if query requires SQL (JOINs, aggregations)
- `canUseFilter`: `true` if PocketBase filter can be used

**Error Response (400 Bad Request):**

```json
{
  "filter": "",
  "error": "The query could not be expressed as a filter."
}
```

#### Example cURL Request

```bash
curl -X POST http://127.0.0.1:8090/api/ai/query \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "orders",
    "query": "pending orders over $100",
    "execute": true
  }'
```

#### Example JavaScript (using PocketBase JS SDK)

```javascript
const pb = new PocketBase('http://127.0.0.1:8090');

// Authenticate first
await pb.admins.authWithPassword('admin@example.com', 'password');

// Make AI query
const response = await pb.send('/api/ai/query', {
  method: 'POST',
  body: {
    collection: 'orders',
    query: 'show me all pending orders from last week',
    execute: true,
    page: 1,
    perPage: 30
  }
});

console.log('Generated Filter:', response.filter);
console.log('Results:', response.results);
console.log('Total Items:', response.totalItems);
```

## API Reference

### Request Format

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `collection` | string | Yes | - | Collection name or ID to query |
| `query` | string | Yes | - | Natural language query |
| `execute` | boolean | No | `false` | Execute filter and return results |
| `page` | integer | No | `1` | Page number for pagination |
| `perPage` | integer | No | `30` | Records per page |

### Response Format

| Field | Type | Description |
|-------|------|-------------|
| `filter` | string | Generated PocketBase filter expression |
| `sql` | string | Generated SQL query (V2, dual/sql mode only) |
| `requiresSQL` | boolean | Whether SQL is required for this query (V2) |
| `canUseFilter` | boolean | Whether filter can be used (V2) |
| `results` | array | Matching records (if `execute: true`) |
| `totalItems` | integer | Total matching records (if `execute: true`) |
| `page` | integer | Current page number (if `execute: true`) |
| `perPage` | integer | Records per page (if `execute: true`) |
| `error` | string | Error message (if query failed) |

### Error Codes

| Status Code | Description |
|-------------|-------------|
| `400` | Bad Request - Invalid input, AI disabled, or query cannot be expressed as filter |
| `401` | Unauthorized - Missing or invalid authentication token |
| `403` | Forbidden - User doesn't have access to collection |
| `404` | Not Found - Collection doesn't exist |
| `500` | Internal Server Error - LLM API error or validation failure |

## Troubleshooting Guide

### AI Query Feature Not Appearing

**Problem**: AI Query option doesn't appear in sidebar.

**Solutions**:
1. Ensure AI Query is enabled in Settings → AI Query
2. Verify you're logged in as a superuser or authenticated user
3. Clear browser cache and reload the Admin UI
4. Rebuild the UI: `cd ui && npm run build`

### "AI Query feature is not enabled" Error

**Problem**: API returns error that AI is not enabled.

**Solutions**:
1. Navigate to Settings → AI Query in Admin UI
2. Toggle "Enable AI Query" to ON
3. Save settings
4. Verify settings persisted by reloading the page

### "Failed to generate filter from query" Error

**Problem**: LLM API call failed.

**Possible Causes**:
- Invalid API key
- Network connectivity issues
- LLM provider API is down
- Rate limiting (too many requests)

**Solutions**:
1. Test connection in Settings → AI Query
2. Verify API key is correct and has sufficient credits
3. Check network connectivity
4. Wait a few minutes and retry (if rate limited)
5. Check LLM provider status page

### "Generated filter is invalid" Error

**Problem**: LLM generated a filter that doesn't match the collection schema.

**Possible Causes**:
- LLM hallucinated field names that don't exist
- Complex query that couldn't be properly translated
- LLM model limitations

**Solutions**:
1. Simplify your query
2. Use more specific field names in your query
3. Try a different LLM model (e.g., `gpt-4o` instead of `gpt-4o-mini`)
4. Manually edit the generated filter if close to correct

### "Collection not found" Error

**Problem**: API returns 404 for collection.

**Solutions**:
1. Verify collection name or ID is correct
2. Ensure collection exists in your PocketBase instance
3. Check collection name spelling (case-sensitive)

### "Only superusers can perform this action" Error

**Problem**: Collection has no `listRule` and user is not superuser.

**Solutions**:
1. Log in as superuser, or
2. Add a `listRule` to the collection to allow authenticated users

### Test Connection Fails

**Problem**: Test Connection button shows error.

**Solutions**:
1. Verify API Base URL is correct
2. Check API key is valid (for providers requiring it)
3. Ensure Ollama is running (if using local Ollama)
4. Check firewall/network settings
5. Verify model name is available for your provider

### Filter Doesn't Match Expected Results

**Problem**: Generated filter doesn't return expected records.

**Solutions**:
1. Review the generated filter expression
2. Test the filter manually in the collection view
3. Simplify your natural language query
4. Check collection schema to ensure fields exist
5. Verify field types match filter operators (e.g., don't use `>` on text fields)

## Best Practices

### Query Writing Tips

1. **Be Specific**: Include field names when possible
   - ✅ "orders with status equals pending"
   - ❌ "pending things"

2. **Use Clear Operators**: Explicitly state comparisons
   - ✅ "price greater than 100"
   - ❌ "expensive items"

3. **Specify Time Ranges**: Use relative or absolute dates
   - ✅ "created in the last 7 days"
   - ✅ "created after 2025-01-01"

4. **Combine Conditions Clearly**: Use "and" and "or" explicitly
   - ✅ "status is active and price is greater than 50"
   - ❌ "active expensive items"

### Performance Considerations

1. **Use `execute: false`** when you only need the filter expression
2. **Set appropriate `perPage`** values (default 30 is usually good)
3. **Use pagination** for large result sets
4. **Cache filter expressions** in your application when possible

### Security Best Practices

1. **Protect API Keys**: Never commit API keys to version control
2. **Use Environment Variables**: Store API keys securely
3. **Monitor Usage**: Track API costs and usage patterns
4. **Set Collection Rules**: Use `listRule` to control access
5. **Review Generated Filters**: Always validate filters before using in production

## Architecture

### Components

**Core AI Services:**
- **`core/ai_settings.go`**: AI settings data structure and validation
- **`services/ai/openai_client.go`**: LLM API client (OpenAI-compatible)
- **`services/ai/schema_extractor.go`**: Collection schema extraction (V2: multi-collection support)
- **`services/ai/prompt_builder.go`**: System and user prompt construction (V2: dual/SQL prompts)
- **`services/ai/prompt_template.go`**: Prompt templates (V2: dual output, SQL terminal)
- **`services/ai/filter_validator.go`**: Filter validation against collection schema
- **`apis/ai_query.go`**: REST API endpoint handler (V2: dual mode support)

**V2 SQL Services:**
- **`services/sql/parser.go`**: SQL statement parser
- **`services/sql/mapper.go`**: SQL type to PocketBase field mapper
- **`services/sql/executor.go`**: SQL execution engine
- **`apis/sql_terminal.go`**: SQL Terminal API endpoints

**UI Components:**
- **`ui/src/components/ai/`**: AI Query components
- **`ui/src/components/sql/`**: SQL Terminal components (V2)
- **`ui/src/pages/SQLTerminal.svelte`**: SQL Terminal page (V2)
- **`ui/src/pages/settings/AI.svelte`**: Settings page

### Data Flow

1. User submits natural language query (UI or API)
2. System extracts collection schema
3. Schema + query → LLM prompt
4. LLM generates PocketBase filter expression
5. Filter is validated against collection schema
6. Filter is optionally executed and results returned

### Security Model

- Respects existing PocketBase authentication
- Enforces collection `listRule` permissions
- Validates all generated filters before execution
- API keys encrypted at rest
- No data sent to LLM provider except schema and query

## Limitations

### Current Limitations

**V1 (Filter Mode):**
1. **Single Collection Queries**: Cannot query across multiple collections in one query
2. **No Aggregations**: Cannot generate aggregation queries (COUNT, SUM, etc.)
3. **No Sorting**: Generated filters don't include sorting (use PocketBase's `sort` parameter)
4. **English Only**: Natural language queries work best in English
5. **Schema-Dependent**: Accuracy depends on clear field names and types

**V2 (Dual/SQL Mode) - Improvements:**
- ✅ Multi-collection queries supported via SQL
- ✅ Aggregations supported via SQL (COUNT, SUM, AVG, etc.)
- ✅ JOINs supported via SQL

**V2 Remaining Limitations:**
1. **Complex subqueries**: Limited support for nested SELECT statements
2. **Transactions**: Individual statements only, no multi-statement transactions
3. **English Only**: Natural language queries work best in English

### Known Issues

- Complex nested queries may not always translate correctly
- Very long queries may exceed LLM context limits
- Some edge cases in filter validation may need refinement

## V2 Features

### Dual Output Mode

Enable dual output mode to get both PocketBase filter and SQL for any query:

```javascript
const response = await pb.send('/api/ai/query', {
  method: 'POST',
  body: {
    collection: 'orders',
    query: 'show me all orders with customer names',
    mode: 'dual',  // Enable dual output
    execute: true
  }
});

console.log('Filter:', response.filter);
console.log('SQL:', response.sql);
console.log('Requires SQL:', response.requiresSQL);
```

### SQL Terminal

For more advanced SQL operations, use the SQL Terminal:
- See [SQL_TERMINAL_FEATURE.md](./SQL_TERMINAL_FEATURE.md)

### When to Use Which Mode

| Use Case | Recommended Mode |
|----------|------------------|
| Simple single-collection queries | `filter` |
| Need both filter and SQL | `dual` |
| JOINs across collections | `dual` or `sql` |
| Aggregations (COUNT, SUM) | `sql` |
| DDL operations | SQL Terminal |

## Future Enhancements

Potential improvements for future versions:

- ~~Multi-collection queries with joins~~ ✅ (V2)
- ~~Aggregation query support (COUNT, SUM, AVG, etc.)~~ ✅ (V2)
- Query history and favorites ✅ (V2 - SQL Terminal)
- Custom prompt templates
- Multi-language support
- Query suggestions and autocomplete
- Usage analytics and cost tracking

## Support

For issues, questions, or contributions:

- Check this documentation first
- Review the [PRD document](../PocketBase_AI_Query_Assistant_PRD.md)
- Check existing GitHub issues
- Create a new issue with detailed information

## Related Documentation

- [SQL Terminal Feature](./SQL_TERMINAL_FEATURE.md) (V2)
- [PocketBase Filter Syntax](https://pocketbase.io/docs/api-records/#filtering)
- [PocketBase API Documentation](https://pocketbase.io/docs)
- [Product Requirements Document](../PocketBase_AI_Query_Assistant_PRD.md)
- [Task List](../PocketBase_AI_Query_TaskList.md)
