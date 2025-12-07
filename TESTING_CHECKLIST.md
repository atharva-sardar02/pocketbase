# AI Query Feature - Detailed Testing Checklist

## Pre-Testing Setup

### 1. Start PocketBase
- [ ] Start the server: `cd examples/base && .\base.exe serve`
- [ ] Verify server starts without errors
- [ ] Access Admin UI at `http://127.0.0.1:8090/_/`
- [ ] Log in as superuser

### 2. Create Test Collections
- [ ] Create a test collection (e.g., "orders") with fields:
  - `status` (select: pending, completed, cancelled)
  - `total` (number)
  - `created` (date)
  - `customer_name` (text)
- [ ] Add at least 5-10 test records with varied data
- [ ] Create a second collection (e.g., "products") with:
  - `name` (text)
  - `price` (number)
  - `in_stock` (bool)
  - `category` (select: electronics, clothing, food)

---

## Settings Configuration Testing

### 3. AI Settings Page Access
- [ ] Navigate to Settings → AI Query
- [ ] Verify page loads without errors
- [ ] Verify all form fields are visible

### 4. Enable/Disable Toggle
- [ ] Toggle "Enable AI Query" OFF → Save
- [ ] Verify AI Query disappears from main sidebar
- [ ] Toggle "Enable AI Query" ON → Save
- [ ] Verify AI Query reappears in sidebar

### 5. Provider Configuration - OpenAI
- [ ] Select "OpenAI" provider
- [ ] Verify Base URL auto-fills to `https://api.openai.com/v1`
- [ ] Enter valid OpenAI API key
- [ ] Select model: `gpt-4o-mini`
- [ ] Set temperature to 0.1
- [ ] Click "Test Connection"
- [ ] Verify success message appears
- [ ] Save settings
- [ ] Reload page and verify settings persisted

### 6. Provider Configuration - Ollama (if available)
- [ ] Ensure Ollama is running locally on port 11434
- [ ] Select "Ollama" provider
- [ ] Verify Base URL auto-fills to `http://localhost:11434/v1`
- [ ] Verify API Key field is hidden
- [ ] Select model: `llama2` or `llama3`
- [ ] Click "Test Connection"
- [ ] Verify success message appears
- [ ] Save settings

### 7. Provider Switching
- [ ] Switch from OpenAI to Anthropic
- [ ] Verify Base URL updates automatically
- [ ] Verify model dropdown updates
- [ ] Switch to Custom provider
- [ ] Verify manual input fields appear
- [ ] Switch back to OpenAI

### 8. Settings Validation
- [ ] Try to save with empty Base URL → Verify error
- [ ] Try to save with invalid URL → Verify error
- [ ] Try to save OpenAI without API key → Verify error
- [ ] Try to save with temperature > 1.0 → Verify error
- [ ] Try to save with temperature < 0.0 → Verify error

---

## Admin UI - AI Query Panel Testing

### 9. Panel Access
- [ ] Click "AI Query" in main sidebar
- [ ] Verify panel loads at `/ai-query`
- [ ] Verify collection dropdown is visible
- [ ] Verify query input textarea is visible
- [ ] Verify search button is visible

### 10. Collection Selection
- [ ] Open collection dropdown
- [ ] Verify all collections appear in list
- [ ] Select "orders" collection
- [ ] Verify collection is selected

### 11. Simple Queries
- [ ] Query: "show me all orders"
- [ ] Verify filter is generated (should be empty or `id != ""`)
- [ ] Verify results appear (if execute enabled)
- [ ] Query: "orders with status equals pending"
- [ ] Verify filter: `status = "pending"`
- [ ] Verify results match filter

### 12. Comparison Queries
- [ ] Query: "orders with total greater than 100"
- [ ] Verify filter: `total > 100`
- [ ] Verify results are correct
- [ ] Query: "orders with total less than 50"
- [ ] Verify filter: `total < 50`
- [ ] Query: "orders created in the last 7 days"
- [ ] Verify filter uses datetime macro: `created >= @now - 604800`

### 13. Combined Queries
- [ ] Query: "pending orders over 100 dollars"
- [ ] Verify filter: `status = "pending" && total > 100`
- [ ] Verify results match both conditions
- [ ] Query: "orders that are pending or cancelled"
- [ ] Verify filter: `status = "pending" || status = "cancelled"`

### 14. Filter Display & Actions
- [ ] After generating a filter, verify it displays in code block
- [ ] Click "Copy Filter" button
- [ ] Paste in notepad → Verify filter was copied correctly
- [ ] Click "Apply Filter" button
- [ ] Verify navigates to collection view with filter applied
- [ ] Verify filter appears in collection filter input

### 15. Results Display
- [ ] Execute query with `execute: true`
- [ ] Verify results list appears
- [ ] Verify total count is displayed
- [ ] Verify pagination info (if applicable)
- [ ] Click "View in Collection" link
- [ ] Verify navigates to collection with filter applied

### 16. Keyboard Shortcuts
- [ ] Type query in textarea
- [ ] Press `Ctrl+Enter` (Windows) or `Cmd+Enter` (Mac)
- [ ] Verify query executes
- [ ] Verify same behavior as clicking Search button

### 17. Error Handling - UI
- [ ] Query with empty collection selected → Verify error message
- [ ] Query with empty query text → Verify error message
- [ ] Query with invalid collection name → Verify error message
- [ ] Disable AI in settings → Try to use AI Query → Verify error
- [ ] Enter invalid API key → Test connection → Verify error

---

## API Endpoint Testing

### 18. Authentication
- [ ] Make request without auth token → Verify 401 error
- [ ] Make request with invalid token → Verify 401 error
- [ ] Make request with valid token → Verify success

### 19. Basic API Request
```bash
curl -X POST http://127.0.0.1:8090/api/ai/query \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "orders",
    "query": "show me all pending orders",
    "execute": false
  }'
```
- [ ] Verify response contains `filter` field
- [ ] Verify filter is valid PocketBase syntax
- [ ] Verify no `results` field (since execute=false)

### 20. API with Execution
```bash
curl -X POST http://127.0.0.1:8090/api/ai/query \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "orders",
    "query": "pending orders over 100",
    "execute": true,
    "page": 1,
    "perPage": 10
  }'
```
- [ ] Verify response contains `filter` field
- [ ] Verify response contains `results` array
- [ ] Verify response contains `totalItems`
- [ ] Verify response contains `page` and `perPage`
- [ ] Verify results match the filter

### 21. API Pagination
- [ ] Request page 1, perPage 5 → Verify 5 results
- [ ] Request page 2, perPage 5 → Verify next 5 results
- [ ] Verify totalItems is consistent across pages

### 22. API Error Cases
- [ ] Missing collection → Verify 400 error
- [ ] Missing query → Verify 400 error
- [ ] Invalid collection name → Verify 404 error
- [ ] AI disabled → Verify 400 error with message
- [ ] Query that can't be expressed as filter → Verify 400 error

### 23. API with Collection Rules
- [ ] Create collection with `listRule: "status = 'active'"`
- [ ] Make API request as non-superuser
- [ ] Verify results respect the listRule
- [ ] Verify generated filter is combined with listRule

---

## Advanced Query Testing

### 24. Date/Time Queries
- [ ] "orders from today"
- [ ] "orders from last week"
- [ ] "orders from last month"
- [ ] "orders created after 2025-01-01"
- [ ] Verify datetime macros are used correctly

### 25. Text Search Queries
- [ ] "orders with customer name containing 'John'"
- [ ] Verify filter uses `~` operator: `customer_name ~ "John"`
- [ ] "products with name starting with 'iPhone'"
- [ ] Verify appropriate text matching

### 26. Select Field Queries
- [ ] "orders with status pending"
- [ ] "orders with status not equal to cancelled"
- [ ] Verify exact match for select values

### 27. Boolean Queries
- [ ] "products that are in stock"
- [ ] Verify filter: `in_stock = true`
- [ ] "products that are not in stock"
- [ ] Verify filter: `in_stock = false`

### 28. Complex Nested Queries
- [ ] "(pending orders over 100) or (completed orders under 50)"
- [ ] Verify proper parentheses and logical operators
- [ ] "orders that are pending and (total > 100 or created this week)"
- [ ] Verify complex nested logic

---

## Edge Cases & Error Scenarios

### 29. Invalid Queries
- [ ] Query: "show me the meaning of life" → Should return error or "INVALID_QUERY"
- [ ] Query: "delete all orders" → Should not generate destructive filter
- [ ] Query: "orders with nonexistent_field = 123" → Should validate and error

### 30. Empty Results
- [ ] Query that matches no records
- [ ] Verify empty results array returned
- [ ] Verify totalItems is 0
- [ ] Verify no error, just empty results

### 31. Large Result Sets
- [ ] Query that matches many records (100+)
- [ ] Verify pagination works correctly
- [ ] Verify performance is acceptable

### 32. Special Characters
- [ ] Query with quotes: "orders with status 'pending'"
- [ ] Query with special characters in field values
- [ ] Verify proper escaping in generated filter

### 33. LLM Provider Errors
- [ ] Temporarily use invalid API key
- [ ] Make query → Verify helpful error message
- [ ] Restore valid API key → Verify works again

---

## Integration Testing

### 34. End-to-End Workflow
1. [ ] Configure AI settings with OpenAI
2. [ ] Create test collection with data
3. [ ] Use AI Query panel to find specific records
4. [ ] Copy generated filter
5. [ ] Apply filter to collection view
6. [ ] Verify same results appear
7. [ ] Use API endpoint with same query
8. [ ] Verify consistent results

### 35. Settings Persistence
- [ ] Configure settings
- [ ] Save and restart PocketBase
- [ ] Verify settings are still configured
- [ ] Verify AI Query still works

### 36. Multiple Collections
- [ ] Query "orders" collection
- [ ] Switch to "products" collection
- [ ] Query "products" collection
- [ ] Verify filters are collection-specific
- [ ] Verify no cross-collection contamination

---

## Performance Testing

### 37. Response Times
- [ ] Measure time for simple query (< 2 seconds expected)
- [ ] Measure time for complex query (< 5 seconds expected)
- [ ] Measure time with execution enabled
- [ ] Verify acceptable performance

### 38. Concurrent Requests
- [ ] Make 3-5 simultaneous API requests
- [ ] Verify all complete successfully
- [ ] Verify no race conditions

---

## Security Testing

### 39. Authentication Enforcement
- [ ] Verify unauthenticated requests are rejected
- [ ] Verify expired tokens are rejected
- [ ] Verify invalid tokens are rejected

### 40. Collection Access Control
- [ ] Create collection with `listRule` restricting access
- [ ] Make API request as non-superuser
- [ ] Verify only accessible records returned
- [ ] Verify filter respects collection rules

### 41. API Key Security
- [ ] Verify API keys are encrypted in database
- [ ] Verify API keys are not exposed in API responses
- [ ] Verify API keys are masked in UI

---

## Final Verification

### 42. Documentation Accuracy
- [ ] Verify all documented features work as described
- [ ] Verify API examples in docs are correct
- [ ] Verify troubleshooting guide addresses real issues

### 43. UI/UX Polish
- [ ] Verify all UI elements are styled correctly
- [ ] Verify loading states work properly
- [ ] Verify error messages are user-friendly
- [ ] Verify success messages appear appropriately

### 44. Browser Compatibility
- [ ] Test in Chrome
- [ ] Test in Firefox
- [ ] Test in Edge
- [ ] Verify consistent behavior

---

## Test Results Summary

**Date:** _______________
**Tester:** _______________
**Environment:** _______________

**Total Tests:** 44
**Passed:** ___
**Failed:** ___
**Skipped:** ___

**Critical Issues Found:**
1. 
2. 
3. 

**Minor Issues Found:**
1. 
2. 
3. 

**Notes:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
