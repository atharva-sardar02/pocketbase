import { writable, derived } from "svelte/store";

// Time period for metrics (1h, 6h, 24h, 7d)
export const dashboardPeriod = writable("24h");

// Refresh interval in seconds (0 = no auto-refresh)
export const dashboardRefreshInterval = writable(30);

// Loading states
export const dashboardLoading = writable(false);
export const overviewLoading = writable(false);
export const requestsLoading = writable(false);
export const latencyLoading = writable(false);
export const errorsLoading = writable(false);
export const endpointsLoading = writable(false);
export const collectionsLoading = writable(false);

// Data stores
export const overviewData = writable({
    totalRequests: 0,
    avgLatency: 0,
    errorRate: 0,
    databaseSize: 0,
    totalErrors: 0,
    period: "24h"
});

export const requestsData = writable([]);
export const latencyData = writable([]);
export const errorsData = writable([]);
export const endpointsData = writable([]);
export const collectionsData = writable([]);

// Error states
export const dashboardError = writable("");

// Computed loading state - true if any data is loading
export const isAnyLoading = derived(
    [overviewLoading, requestsLoading, latencyLoading, errorsLoading, endpointsLoading, collectionsLoading],
    ([$overview, $requests, $latency, $errors, $endpoints, $collections]) => {
        return $overview || $requests || $latency || $errors || $endpoints || $collections;
    }
);

// Format bytes to human readable
export function formatBytes(bytes) {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

// Format latency (milliseconds) to human readable
export function formatLatency(ms) {
    if (ms < 1) return "< 1ms";
    if (ms < 1000) return `${Math.round(ms)}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
}

// Format number with commas
export function formatNumber(num) {
    if (num === null || num === undefined) return "0";
    return num.toLocaleString();
}

// Format percentage
export function formatPercent(value) {
    if (value === null || value === undefined) return "0%";
    return `${value.toFixed(2)}%`;
}

// Reset all dashboard data
export function resetDashboardData() {
    overviewData.set({
        totalRequests: 0,
        avgLatency: 0,
        errorRate: 0,
        databaseSize: 0,
        totalErrors: 0,
        period: "24h"
    });
    requestsData.set([]);
    latencyData.set([]);
    errorsData.set([]);
    endpointsData.set([]);
    collectionsData.set([]);
    dashboardError.set("");
}
