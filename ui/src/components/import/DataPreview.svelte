<script>
    import {
        previewHeaders,
        previewRows,
        previewTotalRows,
        previewErrors,
        previewLoading
    } from "@/stores/import";

    $: hasData = $previewHeaders.length > 0;
    $: hasErrors = $previewErrors.length > 0;
</script>

<div class="data-preview">
    {#if $previewLoading}
        <div class="loading-state">
            <i class="ri-loader-4-line spinning"></i>
            <span>Parsing file...</span>
        </div>
    {:else if hasErrors}
        <div class="error-state">
            <i class="ri-error-warning-line"></i>
            <div class="error-content">
                <strong>Failed to parse file</strong>
                {#each $previewErrors as error}
                    <p>{error}</p>
                {/each}
            </div>
        </div>
    {:else if hasData}
        <div class="preview-info">
            <span class="row-count">
                <i class="ri-table-line"></i>
                {$previewTotalRows} row{$previewTotalRows !== 1 ? "s" : ""} detected
            </span>
            <span class="col-count">
                <i class="ri-layout-column-line"></i>
                {$previewHeaders.length} column{$previewHeaders.length !== 1 ? "s" : ""}
            </span>
        </div>

        <div class="table-container">
            <table class="preview-table">
                <thead>
                    <tr>
                        <th class="row-num">#</th>
                        {#each $previewHeaders as header}
                            <th>{header}</th>
                        {/each}
                    </tr>
                </thead>
                <tbody>
                    {#each $previewRows as row, rowIndex}
                        <tr>
                            <td class="row-num">{rowIndex + 1}</td>
                            {#each row as cell}
                                <td title={cell}>{cell || "-"}</td>
                            {/each}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        {#if $previewRows.length < $previewTotalRows}
            <p class="preview-note">
                Showing first {$previewRows.length} of {$previewTotalRows} rows
            </p>
        {/if}
    {:else}
        <div class="empty-state">
            <i class="ri-file-list-3-line"></i>
            <span>Upload a file to preview data</span>
        </div>
    {/if}
</div>

<style>
    .data-preview {
        width: 100%;
    }

    .loading-state,
    .empty-state,
    .error-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 20px;
        text-align: center;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .loading-state i,
    .empty-state i {
        font-size: 36px;
        color: var(--txtDisabledColor);
        margin-bottom: 12px;
    }

    .loading-state span,
    .empty-state span {
        color: var(--txtHintColor);
    }

    .spinning {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    .error-state {
        flex-direction: row;
        gap: 12px;
        background: rgba(var(--dangerColorRaw), 0.1);
        border-color: var(--dangerColor);
        text-align: left;
    }

    .error-state i {
        font-size: 24px;
        color: var(--dangerColor);
    }

    .error-content strong {
        display: block;
        color: var(--dangerColor);
        margin-bottom: 4px;
    }

    .error-content p {
        margin: 0;
        font-size: 13px;
        color: var(--txtPrimaryColor);
    }

    .preview-info {
        display: flex;
        gap: 20px;
        margin-bottom: 12px;
        padding: 10px 14px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
    }

    .preview-info span {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: var(--txtHintColor);
    }

    .preview-info i {
        font-size: 16px;
    }

    .table-container {
        overflow-x: auto;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .preview-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 13px;
    }

    .preview-table th,
    .preview-table td {
        padding: 10px 12px;
        text-align: left;
        border-bottom: 1px solid var(--baseAlt2Color);
        max-width: 200px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .preview-table th {
        background: var(--baseAlt1Color);
        font-weight: 600;
        color: var(--txtPrimaryColor);
        position: sticky;
        top: 0;
    }

    .preview-table td {
        background: var(--baseColor);
        color: var(--txtPrimaryColor);
    }

    .preview-table tr:last-child td {
        border-bottom: none;
    }

    .preview-table .row-num {
        width: 40px;
        text-align: center;
        color: var(--txtHintColor);
        background: var(--baseAlt1Color);
    }

    .preview-note {
        margin: 10px 0 0;
        font-size: 12px;
        color: var(--txtHintColor);
        text-align: center;
    }
</style>
