---
name: Multi-Table SQL Queries
overview: Extend the AI Query feature to support complex multi-table queries by adding a hybrid system that uses PocketBase filters for simple queries and raw SQL for complex JOINs, aggregates, and subqueries.
todos:
  - id: pr10-schema
    content: "PR #10: Multi-collection schema extraction with relationships"
    status: pending
  - id: pr11-dual-output
    content: "PR #11: Dual output backend - generate both Filter and SQL"
    status: pending
    dependencies:
      - pr10-schema
  - id: pr12-editable-ui
    content: "PR #12: Editable query UI with Filter/SQL tabs"
    status: pending
    dependencies:
      - pr11-dual-output
  - id: pr13-sql-parser
    content: "PR #13: SQL parser and type mapper"
    status: pending
    dependencies:
      - pr10-schema
  - id: pr14-sql-executor
    content: "PR #14: SQL executor using PocketBase APIs"
    status: pending
    dependencies:
      - pr13-sql-parser
  - id: pr15-sql-api
    content: "PR #15: SQL Terminal API endpoints"
    status: pending
    dependencies:
      - pr14-sql-executor
  - id: pr16-sql-ui
    content: "PR #16: SQL Terminal UI with editor, schema browser, results"
    status: pending
    dependencies:
      - pr15-sql-api
  - id: pr17-docs
    content: "PR #17: Update PRD, task list, and documentation"
    status: pending
    dependencies:
      - pr12-editable-ui
      - pr16-sql-ui
---

# Multi-Table SQL Queries + SQL Terminal

## Overview

Two major enhancements to PocketBase AI:

1. **Enhanced AI Query** - Multi-table support with dual output (Filter + SQL tabs), editable queries
2. **SQL Terminal** - Full database console with AI assistance, creates real PocketBase collections

**Key Principle**: All operations create/modify **real PocketBase collections** (not raw SQLite tables), so changes appear immediately in Admin UI with full PocketBase features.

---

## Feature 1: Enhanced AI Query (Multi-Table Support)

### Changes from Current

| Current | Enhanced |

|---------|----------|

| Single collection only | Multiple collections with JOINs |

| Filter output only | Dual output: Filter + SQL (tabbed) |

| Simple queries | Complex queries (aggregates, GROUP BY, subqueries) |

| Read-only output | Editable - modify query before executing |

### Dual Output UI with Tabs

```
┌─────────────────────────────────────────────────┐
│  AI Query Results                               │
├─────────────────────────────────────────────────┤
│  ┌────────┐ ┌─────┐                             │
│  │ Filter │ │ SQL │  ← Click to switch          │
│  └────────┘ └─────┘                             │
├─────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────┐    │
│  │ status = "active" && total > 100        │    │  ← Editable textarea
│  └─────────────────────────────────────────┘    │
│  [Copy] [Execute] [Apply to Collection]         │
├─────────────────────────────────────────────────┤
│  Results: 42 records                            │
│  ... record list ...                            │
└─────────────────────────────────────────────────┘
```

Users can:

- Switch between Filter and SQL tabs
- **Edit** the generated query before executing
- Re-execute modified queries
- Changes reflect immediately in Admin UI

---

## Feature 2: SQL Terminal (New Feature)

A full database management console that creates **real PocketBase collections**.

### UI Design

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
│  Executed in 45ms                                               │
└─────────────────────────────────────────────────────────────────┘
```

### Two Modes

| Mode | Input | Behavior |

|------|-------|----------|

| **AI Mode** | Natural language: "Create a products table with name, price, category" | AI generates SQL → executes |

| **SQL Mode** | Raw SQL: `CREATE TABLE products (...)` | Direct execution |

### SQL → PocketBase Mapping

| SQL Command | PocketBase Action | Result |

|-------------|-------------------|--------|

| `CREATE TABLE products (name TEXT, price REAL)` | Creates PocketBase collection | Full collection with API rules, Admin UI |

| `ALTER TABLE products ADD COLUMN stock INTEGER` | Modifies collection schema | Field added, visible in Admin |

| `DROP TABLE products` | Deletes collection | Removed from Admin UI |

| `INSERT INTO products (name, price) VALUES ('Widget', 9.99)` | Creates record via PocketBase API | Record with proper ID, timestamps |

| `UPDATE products SET price = 19.99 WHERE name = 'Widget'` | Updates record via PocketBase API | Updated timestamp, hooks triggered |

| `DELETE FROM products WHERE stock = 0` | Deletes records via PocketBase API | Proper cleanup, hooks triggered |

| `SELECT * FROM products JOIN categories ON...` | Direct SQL query | Read-only, returns results |

### Field Type Mapping

| SQL Type | PocketBase Field Type |

|----------|----------------------|

| `TEXT`, `VARCHAR` | text |

| `INTEGER`, `INT` | number |

| `REAL`, `FLOAT`, `DOUBLE` | number |

| `BOOLEAN`, `BOOL` | bool |

| `DATE`, `DATETIME`, `TIMESTAMP` | date |

| `TEXT PRIMARY KEY` | autoId (id field) |

| `TEXT REFERENCES table(id)` | relation |

| `TEXT CHECK(value IN ('a','b','c'))` | select |

| `TEXT UNIQUE` | text with unique constraint |

### Features

- **Schema Browser** - See all collections in sidebar, click to explore
- **Syntax highlighting** - SQL keywords colored
- **Auto-complete** - Collection and field name suggestions
- **Command history** - Previous queries saved
- **Export results** - CSV/JSON export
- **Live refresh** - Collections sidebar updates after changes
- **Confirmation dialogs** - For destructive operations (DROP, DELETE)

---

## How Changes Appear in Frontend

### Flow: SQL Terminal → Admin UI

```
User types: CREATE TABLE products (name TEXT, price NUMBER)
                    ↓
SQL Terminal parses SQL, extracts schema
                    ↓
Calls PocketBase Collection API to create collection
                    ↓
Collection created with:
  - Proper ID field (auto-generated)
  - created/updated timestamps
  - Default API rules
  - Admin UI entry
                    ↓
Frontend receives success → Updates sidebar
                    ↓
User can immediately see "products" in:
  - SQL Terminal schema browser
  - Admin UI Collections list
  - AI Query collection dropdown
```

### Flow: AI Query Modification → Database

```
User runs query: "Show active customers"
                    ↓
AI generates: status = "active"
                    ↓
User switches to SQL tab, modifies to:
  UPDATE customers SET status = 'inactive' WHERE last_login < '2024-01-01'
                    ↓
User clicks [Execute]
                    ↓
SQL Terminal API parses UPDATE, executes via PocketBase Records API
                    ↓
Records updated with proper timestamps
                    ↓
User sees: "47 records updated"
                    ↓
Changes visible immediately in Admin UI Records view
```

---

## API Endpoints

### Enhanced AI Query

**`POST /api/ai/query`** (updated)

```json
// Request
{
  "query": "show orders with customer names where total > 100",
  "execute": true
}

// Response
{
  "filter": "total > 100",
  "filterCollection": "orders",
  "sql": "SELECT o.*, c.name as customer_name FROM orders o JOIN customers c ON o.customer = c.id WHERE o.total > 100",
  "canUseFilter": true,
  "results": [...],
  "columns": ["id", "customer", "total", "customer_name"]
}
```

### SQL Terminal - Execute

**`POST /api/sql/execute`**

```json
// Request - SELECT
{
  "sql": "SELECT * FROM customers WHERE city = 'NYC'"
}

// Response
{
  "type": "select",
  "results": [...],
  "columns": ["id", "name", "email", "city"],
  "rowCount": 15,
  "executionTime": 12
}
```
```json
// Request - CREATE TABLE
{
  "sql": "CREATE TABLE products (name TEXT NOT NULL, price REAL, category TEXT)"
}

// Response
{
  "type": "create",
  "collection": "products",
  "fields": [
    {"name": "name", "type": "text", "required": true},
    {"name": "price", "type": "number"},
    {"name": "category", "type": "text"}
  ],
  "message": "Collection 'products' created successfully"
}
```
```json
// Request - INSERT
{
  "sql": "INSERT INTO products (name, price) VALUES ('Widget', 9.99)"
}

// Response
{
  "type": "insert",
  "collection": "products",
  "recordId": "abc123xyz",
  "rowsAffected": 1
}
```

### SQL Terminal - AI Mode

**`POST /api/sql/ai`**

```json
// Request
{
  "query": "Create a products table with name, price, and category fields"
}

// Response
{
  "sql": "CREATE TABLE products (\n  name TEXT NOT NULL,\n  price REAL,\n  category TEXT\n)",
  "explanation": "Creates a 'products' collection with three fields"
}
```

---

## Files to Create/Modify

### Backend (Go)

| File | Action | Purpose |

|------|--------|---------|

| [`services/ai/schema_extractor.go`](services/ai/schema_extractor.go) | MODIFY | Extract ALL collections + relationships |

| [`services/ai/prompt_template.go`](services/ai/prompt_template.go) | MODIFY | Add SQL syntax, multi-table examples |

| [`services/sql/parser.go`](services/sql/parser.go) | CREATE | Parse SQL statements, extract intent |

| [`services/sql/executor.go`](services/sql/executor.go) | CREATE | Execute SQL via PocketBase APIs |

| [`services/sql/mapper.go`](services/sql/mapper.go) | CREATE | Map SQL types to PocketBase field types |

| [`apis/ai_query.go`](apis/ai_query.go) | MODIFY | Add dual output (filter + SQL) |

| [`apis/sql_terminal.go`](apis/sql_terminal.go) | CREATE | SQL Terminal API endpoints |

| [`core/ai_settings.go`](core/ai_settings.go) | MODIFY | Add SQL terminal settings |

### Frontend (Svelte)

| File | Action | Purpose |

|------|--------|---------|

| [`ui/src/components/ai/AIQueryPanel.svelte`](ui/src/components/ai/AIQueryPanel.svelte) | MODIFY | Add tabs, editable queries |

| [`ui/src/components/ai/QueryTabs.svelte`](ui/src/components/ai/QueryTabs.svelte) | CREATE | Filter/SQL tab switcher |

| [`ui/src/components/ai/EditableCodeBlock.svelte`](ui/src/components/ai/EditableCodeBlock.svelte) | CREATE | Editable query with syntax highlight |

| [`ui/src/pages/SQLTerminal.svelte`](ui/src/pages/SQLTerminal.svelte) | CREATE | Main SQL Terminal page |

| [`ui/src/components/sql/SQLEditor.svelte`](ui/src/components/sql/SQLEditor.svelte) | CREATE | Code editor component |

| [`ui/src/components/sql/SchemaExplorer.svelte`](ui/src/components/sql/SchemaExplorer.svelte) | CREATE | Collections sidebar browser |

| [`ui/src/components/sql/ResultsTable.svelte`](ui/src/components/sql/ResultsTable.svelte) | CREATE | Dynamic results display |

| [`ui/src/components/sql/QueryHistory.svelte`](ui/src/components/sql/QueryHistory.svelte) | CREATE | Command history |

| [`ui/src/stores/sql.js`](ui/src/stores/sql.js) | CREATE | SQL terminal state |

| [`ui/src/App.svelte`](ui/src/App.svelte) | MODIFY | Add SQL Terminal to sidebar |

---

## Implementation PRs

### PR #10: Multi-Collection Schema Extraction (4-5 hrs)

- Extract all collections with relationships
- Update prompt template for multi-table queries
- Add relationship detection

### PR #11: Dual Output Backend (5-6 hrs)

- Generate both filter AND SQL for queries
- Detect when filter is possible vs SQL-only
- Return both in API response

### PR #12: Editable Query UI (4-5 hrs)

- Create tabbed interface (Filter/SQL)
- Add editable code blocks
- Re-execute modified queries

### PR #13: SQL Parser & Mapper (6-7 hrs)

- Parse SQL statements (CREATE, ALTER, DROP, INSERT, UPDATE, DELETE, SELECT)
- Map SQL types to PocketBase field types
- Extract table/column information

### PR #14: SQL Executor (6-7 hrs)

- Execute DDL via PocketBase Collection APIs
- Execute DML via PocketBase Record APIs
- Execute SELECT via direct SQLite query
- Handle transactions

### PR #15: SQL Terminal API (5-6 hrs)

- `/api/sql/execute` endpoint
- `/api/sql/ai` endpoint
- Authentication, error handling

### PR #16: SQL Terminal UI (8-10 hrs)

- SQL Editor with syntax highlighting
- Schema Explorer sidebar
- Results table with export
- Query history
- AI/SQL mode toggle

### PR #17: Documentation & Polish (3-4 hrs)

- Update PRD and task list
- Add SQL Terminal documentation
- Test all flows end-to-end

---

## Total Estimated Time

| Component | Hours |

|-----------|-------|

| Enhanced AI Query (PRs #10-12) | 13-16 |

| SQL Terminal Backend (PRs #13-15) | 17-20 |

| SQL Terminal UI (PR #16) | 8-10 |

| Documentation (PR #17) | 3-4 |

| **Total** | **41-50 hours** |

---

## Example Workflows

### Creating a New Collection via SQL Terminal

```sql
-- User types in SQL Mode:
CREATE TABLE blog_posts (
  title TEXT NOT NULL,
  content TEXT,
  author TEXT REFERENCES users(id),
  status TEXT CHECK(status IN ('draft', 'published', 'archived')),
  published_at DATETIME
);

-- Result:
-- ✅ Collection 'blog_posts' created
-- Fields:
--   - title (text, required)
--   - content (text)
--   - author (relation → users)
--   - status (select: draft, published, archived)
--   - published_at (date)
```

### Complex Query in AI Query

```
User: "Show me total revenue by customer with their names, only customers who spent more than $1000"

Filter Tab: [Not available - requires JOIN and GROUP BY]

SQL Tab:
SELECT 
  c.name,
  c.email,
  SUM(o.total) as total_revenue
FROM customers c
JOIN orders o ON c.id = o.customer
GROUP BY c.id
HAVING total_revenue > 1000
ORDER BY total_revenue DESC

[Execute] → Shows results table with customer names and revenue
```