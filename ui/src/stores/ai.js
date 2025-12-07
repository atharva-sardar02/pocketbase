import { writable } from "svelte/store";

// Current query text
export const aiQuery = writable("");

// Generated filter expression
export const aiFilter = writable("");

// Generated SQL query (V2)
export const aiSQL = writable("");

// Whether SQL is required for this query (V2)
export const aiRequiresSQL = writable(false);

// Whether filter can be used (V2)
export const aiCanUseFilter = writable(true);

// Query mode: "filter", "dual", or "sql" (V2)
export const aiMode = writable("dual");

// Active output tab: "filter" or "sql" (V2)
export const aiActiveTab = writable("filter");

// Query results (array of records)
export const aiResults = writable([]);

// Total items count
export const aiTotalItems = writable(0);

// Current page
export const aiPage = writable(1);

// Items per page
export const aiPerPage = writable(30);

// Loading state
export const aiLoading = writable(false);

// Error message
export const aiError = writable("");

// Selected collection for query
export const aiCollection = writable("");

// Reset all AI state
export function resetAIState() {
    aiQuery.set("");
    aiFilter.set("");
    aiSQL.set("");
    aiRequiresSQL.set(false);
    aiCanUseFilter.set(true);
    aiActiveTab.set("filter");
    aiResults.set([]);
    aiTotalItems.set(0);
    aiPage.set(1);
    aiPerPage.set(30);
    aiLoading.set(false);
    aiError.set("");
}

// Set response from dual output API
export function setDualResponse(response) {
    aiFilter.set(response.filter || "");
    aiSQL.set(response.sql || "");
    aiRequiresSQL.set(response.requiresSQL || false);
    aiCanUseFilter.set(response.canUseFilter !== false);
    
    // Auto-select tab based on response
    if (response.requiresSQL && !response.canUseFilter) {
        aiActiveTab.set("sql");
    } else {
        aiActiveTab.set("filter");
    }
    
    if (response.results) {
        aiResults.set(response.results);
        aiTotalItems.set(response.totalItems || 0);
        aiPage.set(response.page || 1);
        aiPerPage.set(response.perPage || 30);
    }
}

