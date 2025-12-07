import { writable, derived } from "svelte/store";

// Wizard step (1-4)
export const importStep = writable(1);

// Selected collection
export const selectedCollection = writable("");

// File data
export const uploadedFile = writable(null);
export const fileFormat = writable("csv"); // "csv" or "json"
export const fileDelimiter = writable(",");

// Preview data from API
export const previewHeaders = writable([]);
export const previewRows = writable([]);
export const previewTotalRows = writable(0);
export const previewErrors = writable([]);

// Field mapping: { sourceColumn: targetField }
export const fieldMapping = writable({});

// Available collection fields from validation
export const collectionFields = writable({});

// Import progress
export const importInProgress = writable(false);
export const importResults = writable({
    totalRows: 0,
    successCount: 0,
    failureCount: 0,
    errors: []
});

// Loading states
export const previewLoading = writable(false);
export const validateLoading = writable(false);
export const executeLoading = writable(false);

// Error state
export const importError = writable("");

// Computed: check if mapping is valid (at least one field mapped)
export const hasMappedFields = derived(fieldMapping, ($mapping) => {
    return Object.values($mapping).some(v => v && v !== "-");
});

// Computed: get mapped field count
export const mappedFieldCount = derived(fieldMapping, ($mapping) => {
    return Object.values($mapping).filter(v => v && v !== "-").length;
});

// Reset all import data
export function resetImportData() {
    importStep.set(1);
    selectedCollection.set("");
    uploadedFile.set(null);
    fileFormat.set("csv");
    fileDelimiter.set(",");
    previewHeaders.set([]);
    previewRows.set([]);
    previewTotalRows.set(0);
    previewErrors.set([]);
    fieldMapping.set({});
    collectionFields.set({});
    importInProgress.set(false);
    importResults.set({
        totalRows: 0,
        successCount: 0,
        failureCount: 0,
        errors: []
    });
    previewLoading.set(false);
    validateLoading.set(false);
    executeLoading.set(false);
    importError.set("");
}

// Reset to step 1 (start over)
export function startOver() {
    resetImportData();
}

// Go to next step
export function nextStep() {
    importStep.update(n => Math.min(n + 1, 4));
}

// Go to previous step
export function prevStep() {
    importStep.update(n => Math.max(n - 1, 1));
}

// Go to specific step
export function goToStep(step) {
    if (step >= 1 && step <= 4) {
        importStep.set(step);
    }
}

// Initialize field mapping from headers
export function initializeMapping(headers) {
    const mapping = {};
    headers.forEach(header => {
        mapping[header] = ""; // Empty means not mapped
    });
    fieldMapping.set(mapping);
}

// Auto-detect mapping based on field name similarity
export function autoDetectMapping(headers, fields) {
    const mapping = {};
    const fieldNames = Object.keys(fields);
    
    headers.forEach(header => {
        const headerLower = header.toLowerCase().replace(/[_\s-]/g, "");
        
        // Try exact match first
        const exactMatch = fieldNames.find(f => f.toLowerCase() === headerLower);
        if (exactMatch) {
            mapping[header] = exactMatch;
            return;
        }
        
        // Try partial match
        const partialMatch = fieldNames.find(f => 
            f.toLowerCase().includes(headerLower) || 
            headerLower.includes(f.toLowerCase())
        );
        if (partialMatch) {
            mapping[header] = partialMatch;
            return;
        }
        
        // No match
        mapping[header] = "";
    });
    
    fieldMapping.set(mapping);
}
