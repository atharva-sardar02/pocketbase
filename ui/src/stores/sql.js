import { writable, derived } from "svelte/store";

// Current SQL query text
export const sqlQuery = writable("");

// AI natural language query (for AI mode)
export const sqlAIQuery = writable("");

// Mode: "sql" for direct SQL, "ai" for AI-generated SQL
export const sqlMode = writable("sql");

// Query results
export const sqlResults = writable([]);

// Result columns (for dynamic table headers)
export const sqlColumns = writable([]);

// Total rows affected/returned
export const sqlTotalRows = writable(0);

// Rows affected by write operations
export const sqlRowsAffected = writable(0);

// Multi-statement support
export const sqlIsMulti = writable(false);
export const sqlMultiResults = writable([]); // Array of individual statement results
export const sqlTotalStatements = writable(0);
export const sqlSuccessfulCount = writable(0);
export const sqlFailedCount = writable(0);

// Loading state
export const sqlLoading = writable(false);

// Error message
export const sqlError = writable("");

// Success message
export const sqlSuccess = writable("");

// Query history (stored in localStorage)
const HISTORY_KEY = "sql_terminal_history";
const MAX_HISTORY = 50;

function loadHistory() {
    try {
        const stored = localStorage.getItem(HISTORY_KEY);
        return stored ? JSON.parse(stored) : [];
    } catch {
        return [];
    }
}

function saveHistory(history) {
    try {
        localStorage.setItem(HISTORY_KEY, JSON.stringify(history.slice(0, MAX_HISTORY)));
    } catch {
        // Ignore storage errors
    }
}

export const sqlHistory = writable(loadHistory());

// Add query to history
export function addToHistory(query, type = "sql") {
    if (!query || !query.trim()) return;
    
    sqlHistory.update(history => {
        // Remove duplicate if exists
        const filtered = history.filter(h => h.query !== query);
        // Add to front
        const newHistory = [
            { query, type, timestamp: Date.now() },
            ...filtered
        ].slice(0, MAX_HISTORY);
        saveHistory(newHistory);
        return newHistory;
    });
}

// Clear history
export function clearHistory() {
    sqlHistory.set([]);
    localStorage.removeItem(HISTORY_KEY);
}

// Schema data (collections and fields)
export const sqlSchema = writable([]);

// Schema loading state
export const sqlSchemaLoading = writable(false);

// Confirmation required for destructive operations
export const sqlNeedsConfirmation = writable(false);
export const sqlConfirmMessage = writable("");
export const sqlPendingQuery = writable("");

// Reset all SQL state
export function resetSQLState() {
    sqlQuery.set("");
    sqlAIQuery.set("");
    sqlResults.set([]);
    sqlColumns.set([]);
    sqlTotalRows.set(0);
    sqlRowsAffected.set(0);
    sqlLoading.set(false);
    sqlError.set("");
    sqlSuccess.set("");
    sqlNeedsConfirmation.set(false);
    sqlConfirmMessage.set("");
    sqlPendingQuery.set("");
    // Multi-statement reset
    sqlIsMulti.set(false);
    sqlMultiResults.set([]);
    sqlTotalStatements.set(0);
    sqlSuccessfulCount.set(0);
    sqlFailedCount.set(0);
}

// Check if query is potentially destructive
export function isDestructiveQuery(query) {
    if (!query) return false;
    const upper = query.toUpperCase().trim();
    return (
        upper.startsWith("DELETE") ||
        upper.startsWith("DROP") ||
        upper.startsWith("TRUNCATE") ||
        upper.startsWith("ALTER")
    );
}

// Derived store: has results to display
export const hasResults = derived(
    [sqlResults, sqlColumns],
    ([$results, $columns]) => $results.length > 0 || $columns.length > 0
);

// Derived store: is in AI mode
export const isAIMode = derived(sqlMode, $mode => $mode === "ai");

