import { writable } from "svelte/store";

// Current query text
export const aiQuery = writable("");

// Generated filter expression
export const aiFilter = writable("");

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
    aiResults.set([]);
    aiTotalItems.set(0);
    aiPage.set(1);
    aiPerPage.set(30);
    aiLoading.set(false);
    aiError.set("");
}

