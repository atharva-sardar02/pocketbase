<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let columns = [];
    export let rows = [];
    export let totalRows = 0;
    export let rowsAffected = 0;
    export let loading = false;
    export let error = "";
    export let success = "";

    const dispatch = createEventDispatcher();

    let sortColumn = "";
    let sortDirection = "asc";
    let tableWrapper;

    $: displayRows = sortRows(rows, sortColumn, sortDirection);
    $: hasData = columns.length > 0 || rows.length > 0;

    function sortRows(data, column, direction) {
        if (!column || !data.length) return data;
        return [...data].sort((a, b) => {
            const aVal = a[column];
            const bVal = b[column];
            if (aVal === bVal) return 0;
            if (aVal === null || aVal === undefined) return 1;
            if (bVal === null || bVal === undefined) return -1;
            const cmp = aVal < bVal ? -1 : 1;
            return direction === "asc" ? cmp : -cmp;
        });
    }

    function toggleSort(column) {
        if (sortColumn === column) {
            sortDirection = sortDirection === "asc" ? "desc" : "asc";
        } else {
            sortColumn = column;
            sortDirection = "asc";
        }
    }

    function formatValue(value) {
        if (value === null || value === undefined) return "NULL";
        if (typeof value === "boolean") return value ? "true" : "false";
        if (typeof value === "object") return JSON.stringify(value);
        return String(value);
    }

    function isNullValue(value) {
        return value === null || value === undefined;
    }

    function exportCSV() {
        if (!rows.length) return;
        
        const headers = columns.join(",");
        const csvRows = rows.map(row => 
            columns.map(col => {
                const val = row[col];
                if (val === null || val === undefined) return "";
                const str = String(val);
                if (str.includes(",") || str.includes('"') || str.includes("\n")) {
                    return `"${str.replace(/"/g, '""')}"`;
                }
                return str;
            }).join(",")
        );
        
        const csv = [headers, ...csvRows].join("\n");
        downloadFile(csv, "query_results.csv", "text/csv");
    }

    function exportJSON() {
        if (!rows.length) return;
        const json = JSON.stringify(rows, null, 2);
        downloadFile(json, "query_results.json", "application/json");
    }

    function downloadFile(content, filename, type) {
        const blob = new Blob([content], { type });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    function copyToClipboard() {
        if (!rows.length) return;
        const json = JSON.stringify(rows, null, 2);
        navigator.clipboard.writeText(json);
        dispatch("copy");
    }
</script>

<div class="results-table-container">
    <div class="results-header">
        <div class="results-info">
            {#if loading}
                <i class="ri-loader-4-line animate-spin"></i>
                <span>Executing...</span>
            {:else if error}
                <i class="ri-error-warning-line error-icon"></i>
                <span class="error-text">Error</span>
            {:else if success}
                <i class="ri-check-line success-icon"></i>
                <span>{success}</span>
            {:else if hasData}
                <span>{totalRows || rows.length} row{(totalRows || rows.length) !== 1 ? "s" : ""}</span>
                {#if rowsAffected > 0}
                    <span class="affected">({rowsAffected} affected)</span>
                {/if}
            {:else}
                <span class="hint">No results</span>
            {/if}
        </div>
        {#if hasData && !loading && !error}
            <div class="results-actions">
                <button
                    type="button"
                    class="btn btn-xs btn-transparent"
                    on:click={copyToClipboard}
                    use:tooltip={"Copy JSON"}
                >
                    <i class="ri-file-copy-line"></i>
                </button>
                <button
                    type="button"
                    class="btn btn-xs btn-transparent"
                    on:click={exportCSV}
                    use:tooltip={"Export CSV"}
                >
                    <i class="ri-file-excel-line"></i>
                </button>
                <button
                    type="button"
                    class="btn btn-xs btn-transparent"
                    on:click={exportJSON}
                    use:tooltip={"Export JSON"}
                >
                    <i class="ri-file-code-line"></i>
                </button>
            </div>
        {/if}
    </div>

    <div class="results-body">
        {#if loading}
            <div class="loading-state">
                <div class="loader"></div>
            </div>
        {:else if error}
            <div class="error-state">
                <i class="ri-error-warning-line"></i>
                <pre class="error-message">{error}</pre>
            </div>
        {:else if !hasData}
            <div class="empty-state">
                <i class="ri-table-line"></i>
                <span>Execute a query to see results</span>
            </div>
        {:else}
            {#key columns.join(',')}
                <div class="table-wrapper" bind:this={tableWrapper}>
                    <table class="data-table">
                        <thead>
                            <tr>
                                {#each columns as column}
                                    <th on:click={() => toggleSort(column)}>
                                        <span class="column-name">{column}</span>
                                        {#if sortColumn === column}
                                            <i class="ri-arrow-{sortDirection === 'asc' ? 'up' : 'down'}-s-line"></i>
                                        {/if}
                                    </th>
                                {/each}
                            </tr>
                        </thead>
                        <tbody>
                            {#each displayRows as row}
                                <tr>
                                    {#each columns as column}
                                        <td class:null-value={isNullValue(row[column])}>
                                            {formatValue(row[column])}
                                        </td>
                                    {/each}
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            {/key}
        {/if}
    </div>
</div>

<style>
    .results-table-container {
        display: flex;
        flex-direction: column;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        overflow: hidden;
        height: 100%;
    }
    .results-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 8px 12px;
        background: var(--baseAlt1Color);
        border-bottom: 1px solid var(--baseAlt2Color);
        flex-shrink: 0;
    }
    .results-info {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        color: var(--txtSecondaryColor);
    }
    .results-info i {
        font-size: 14px;
    }
    .error-icon {
        color: var(--dangerColor);
    }
    .error-text {
        color: var(--dangerColor);
    }
    .success-icon {
        color: var(--successColor);
    }
    .affected {
        color: var(--txtHintColor);
    }
    .hint {
        color: var(--txtHintColor);
    }
    .results-actions {
        display: flex;
        gap: 4px;
    }
    .results-body {
        flex: 1;
        min-height: 100px;
        overflow: auto;
    }
    .loading-state,
    .empty-state,
    .error-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 10px;
        padding: 40px 20px;
        color: var(--txtHintColor);
        font-size: 13px;
    }
    .empty-state i,
    .error-state i {
        font-size: 32px;
    }
    .error-state {
        color: var(--dangerColor);
    }
    .error-message {
        margin: 0;
        padding: 10px;
        background: var(--dangerAltColor);
        border-radius: var(--baseRadius);
        font-size: 12px;
        max-width: 100%;
        overflow-x: auto;
        white-space: pre-wrap;
        word-break: break-word;
    }
    .loader {
        width: 24px;
        height: 24px;
        border: 2px solid var(--baseAlt3Color);
        border-top-color: var(--primaryColor);
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }
    .table-wrapper {
        width: 100%;
        margin: 0;
        padding: 0;
    }
    .data-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 13px;
        margin: 0;
        padding: 0;
        table-layout: auto;
    }
    .data-table th,
    .data-table td {
        padding: 10px 16px;
        text-align: left;
        border-bottom: 1px solid var(--baseAlt1Color);
        white-space: nowrap;
        min-width: 80px;
    }
    .data-table th {
        background: var(--baseAlt1Color);
        font-weight: 600;
        color: var(--txtSecondaryColor);
        position: sticky;
        top: 0;
        cursor: pointer;
        user-select: none;
        z-index: 1;
        border-right: 1px solid var(--baseAlt2Color);
    }
    .data-table th:last-child {
        border-right: none;
    }
    .data-table th:hover {
        background: var(--baseAlt2Color);
    }
    .data-table th i {
        font-size: 12px;
        margin-left: 4px;
        vertical-align: middle;
    }
    .data-table td {
        font-family: var(--monospaceFontFamily);
        color: var(--txtPrimaryColor);
        max-width: 400px;
        overflow: hidden;
        text-overflow: ellipsis;
        border-right: 1px solid var(--baseAlt1Color);
    }
    .data-table td:last-child {
        border-right: none;
    }
    .data-table td.null-value {
        color: var(--txtDisabledColor);
        font-style: italic;
    }
    .data-table tbody tr:hover {
        background: var(--baseAlt1Color);
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>

