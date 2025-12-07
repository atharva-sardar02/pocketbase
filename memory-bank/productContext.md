# Product Context: PocketBase AI Query Assistant

## Why This Project Exists

PocketBase is a powerful backend-as-a-service platform, but several pain points exist for different user types:
- **Non-technical users** find the filter syntax intimidating
- **Developers** want SQL-like access without raw database access
- **Administrators** need visibility into system health beyond text logs
- **Data managers** need efficient ways to bulk import data

This project addresses all these needs through AI-powered features and developer tools.

## Problems It Solves

### Problem 1: Technical Query Barrier (V1)
**Current State:** Users must learn PocketBase filter syntax (`status = "active" && total > 100`)  
**Pain Point:** Non-technical users avoid querying data or rely on developers  
**Solution:** Natural language queries ("show me active orders over $100")

### Problem 2: Slow Query Construction (V1)
**Current State:** Power users manually construct complex filters  
**Pain Point:** Time-consuming, error-prone, requires syntax knowledge  
**Solution:** Describe intent in plain English, get valid filter instantly

### Problem 3: No SQL Access (V2)
**Current State:** Developers can't use familiar SQL syntax  
**Pain Point:** Must learn PocketBase-specific filter syntax  
**Solution:** SQL Terminal with full SQL support that creates real records

### Problem 4: Limited System Visibility (V3)
**Current State:** Only text-based logs available for monitoring  
**Pain Point:** Hard to spot trends, no visual performance overview  
**Solution:** Real-time metrics dashboard with charts and graphs

### Problem 5: Manual Data Entry (V3)
**Current State:** Records must be created one-by-one or via API  
**Pain Point:** Bulk data import requires custom scripts  
**Solution:** Data Import Wizard supporting CSV and JSON files

## How It Should Work

### V1: AI Query Flow
1. User opens Admin UI → Clicks "AI Query" in sidebar
2. Types natural language query → "pending orders over $100 from last week"
3. System generates PocketBase filter → Validates → Executes
4. User sees filter expression (copyable) + matching records

### V2: SQL Terminal Flow
1. User opens SQL Terminal from sidebar
2. Chooses mode: AI Mode or SQL Mode
3. **AI Mode:** Types "create a products table" → Gets SQL generated
4. **SQL Mode:** Types raw SQL → Executes directly
5. Results displayed in table → Can export to CSV/JSON

### V3: Dashboard Flow
1. User opens Dashboard from sidebar
2. Sees overview cards: Requests, Latency, Error Rate, DB Size
3. Charts show trends: Requests over time, Latency percentiles
4. Bar chart shows top endpoints
5. Table shows collection record counts
6. Auto-refreshes every 30 seconds

### V3: Import Wizard Flow
1. User opens Import from sidebar
2. **Step 1:** Selects collection, uploads CSV/JSON file
3. **Step 2:** Previews data (first 5 rows)
4. **Step 3:** Maps CSV columns to collection fields
5. **Step 4:** Imports with progress bar, sees success/error counts

## User Personas

### Persona 1: Non-Technical Admin/Business User
- **Profile:** Uses PocketBase Admin UI to manage application data
- **Technical Level:** Low - not familiar with filter syntax or SQL
- **Goal:** Find and analyze data quickly without learning syntax
- **Uses:** AI Query (V1), Dashboard for monitoring (V3)
- **Key Needs:**
  - Plain English queries
  - Visual system health overview
  - Easy data import from spreadsheets

### Persona 2: Developer/Power User
- **Profile:** Building applications on PocketBase, knows SQL
- **Technical Level:** High
- **Goal:** Efficient database management, familiar SQL syntax
- **Uses:** SQL Terminal (V2), AI Query (V1), Data Import (V3)
- **Key Needs:**
  - SQL access for familiar operations
  - Bulk operations (multi-statement SQL, data import)
  - API access for automation

### Persona 3: PocketBase Administrator
- **Profile:** Manages PocketBase instance, responsible for operations
- **Technical Level:** High
- **Goal:** Monitor system health, manage data, control costs
- **Uses:** Dashboard (V3), All features
- **Key Needs:**
  - Real-time performance metrics
  - Error rate visibility
  - Database size monitoring
  - LLM cost awareness

### Persona 4: Data Manager
- **Profile:** Migrates data from other systems, manages bulk operations
- **Technical Level:** Medium
- **Goal:** Efficiently import large datasets
- **Uses:** Data Import Wizard (V3)
- **Key Needs:**
  - CSV/JSON import support
  - Field mapping flexibility
  - Error handling for partial imports
  - Progress visibility

## User Experience Goals

1. **Intuitive:** No training required - features self-explanatory
2. **Fast:** Query results in 1-3 seconds, dashboard loads instantly
3. **Transparent:** Users see generated queries, import details
4. **Forgiving:** Helpful error messages, retry capabilities
5. **Visual:** Charts and graphs provide insights at a glance
6. **Efficient:** Bulk operations save time vs. manual entry
7. **Secure:** All features respect PocketBase auth and rules

## Feature Matrix by Persona

| Feature | Business User | Developer | Admin | Data Manager |
|---------|---------------|-----------|-------|--------------|
| AI Query | ⭐ Primary | ✅ Uses | ✅ Uses | ➖ Rarely |
| SQL Terminal | ➖ Rarely | ⭐ Primary | ✅ Uses | ✅ Uses |
| Dashboard | ✅ Uses | ✅ Uses | ⭐ Primary | ➖ Rarely |
| Data Import | ✅ Uses | ✅ Uses | ✅ Uses | ⭐ Primary |

## Success Metrics

### V1 Metrics
- Query success rate (% of queries returning valid results)
- Time saved (manual query construction vs. AI)
- Learning rate (users transitioning to manual filters)

### V2 Metrics
- SQL Terminal adoption (% of developers using it)
- Operations performed (CREATE, INSERT, SELECT counts)
- Export usage (CSV/JSON downloads)

### V3 Metrics
- Dashboard daily active users
- Import volume (records imported per week)
- Import success rate (% successful vs. failed)
- Average latency visibility (are admins catching issues earlier?)
