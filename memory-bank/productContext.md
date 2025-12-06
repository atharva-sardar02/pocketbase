# Product Context: PocketBase AI Query Assistant

## Why This Project Exists

PocketBase is a powerful backend-as-a-service platform, but its filter syntax can be intimidating for non-technical users. Business users, content managers, and administrators often need to query data but lack the technical knowledge to write filter expressions. This creates a barrier to data exploration and analysis.

## Problems It Solves

### Problem 1: Technical Barrier
**Current State:** Users must learn PocketBase filter syntax (`status = "active" && total > 100`)  
**Pain Point:** Non-technical users avoid querying data or rely on developers  
**Solution:** Natural language queries ("show me active orders over $100")

### Problem 2: Slow Query Construction
**Current State:** Power users manually construct complex filters  
**Pain Point:** Time-consuming, error-prone, requires syntax knowledge  
**Solution:** Describe intent in plain English, get valid filter instantly

### Problem 3: Learning Curve
**Current State:** No way to learn filter syntax gradually  
**Pain Point:** Users either know it or don't  
**Solution:** Show generated filter expressions so users can learn over time

## How It Should Work

### User Experience Flow

1. **User opens Admin UI** → Navigates to collection (e.g., "orders")
2. **User clicks "AI Query" in sidebar** → AI Query panel opens
3. **User types natural language query** → "show me pending orders over $100 from last week"
4. **System processes query:**
   - Extracts collection schema (field names, types, relations)
   - Builds optimized prompt with schema context
   - Calls LLM to generate filter expression
   - Validates generated filter against schema
   - Executes filter (if requested)
5. **User sees results:**
   - Generated filter expression (with copy button)
   - Matching records (if executed)
   - Option to refine query or apply filter to collection view

### Key User Interactions

**Primary Interaction: Natural Language Input**
- Text input box for typing queries
- Collection selector (defaults to current collection)
- Search button with loading state
- Keyboard shortcut (Ctrl+Enter)

**Secondary Interaction: Filter Display**
- Code block showing generated PocketBase filter
- Copy to clipboard button
- "Apply Filter" button (navigates to collection with filter)

**Tertiary Interaction: Results Preview**
- Results count
- Basic record preview (id, first few fields)
- "View in Collection" link

**Configuration Interaction: Settings**
- Enable/Disable toggle
- LLM provider selection
- API key configuration
- Model selection
- Temperature adjustment
- Test connection button

## User Personas

### Persona 1: Non-Technical Admin/Business User
- **Profile:** Uses PocketBase Admin UI to manage application data
- **Technical Level:** Low - not familiar with filter syntax or SQL
- **Goal:** Find and analyze data quickly without learning syntax
- **Key Needs:**
  - Plain English queries
  - Automatic field name understanding
  - Helpful error messages
  - Learning opportunity (see generated filters)

### Persona 2: Developer/Power User
- **Profile:** Building applications on PocketBase, knows filter syntax
- **Technical Level:** High
- **Goal:** Faster query construction for complex filters
- **Key Needs:**
  - API access for building AI-powered search
  - Copy-paste filter expressions
  - Respect for existing security model

### Persona 3: PocketBase Administrator
- **Profile:** Manages PocketBase instance, responsible for configuration and security
- **Technical Level:** High
- **Goal:** Control costs, privacy, and feature availability
- **Key Needs:**
  - LLM provider configuration
  - Enable/disable control
  - Usage monitoring
  - Cost management

## User Experience Goals

1. **Intuitive:** No training required - users understand how to use it immediately
2. **Fast:** Query results appear in 1-3 seconds (LLM latency)
3. **Transparent:** Users see what filter was generated and can learn from it
4. **Forgiving:** Helpful error messages guide users to successful queries
5. **Flexible:** Works with multiple LLM providers for different deployment needs
6. **Secure:** Respects existing authentication and authorization rules

## Success Metrics (Future)

- User adoption rate (percentage of users who try AI Query)
- Query success rate (percentage of queries that return valid results)
- Time saved (average time to construct query manually vs. AI)
- Learning rate (users who transition from AI queries to manual filters)
- Cost per query (monitoring LLM API costs)



