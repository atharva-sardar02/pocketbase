<script>
    import { createEventDispatcher } from "svelte";
    import {
        importInProgress,
        importResults,
        executeLoading,
        importError
    } from "@/stores/import";

    const dispatch = createEventDispatcher();

    $: totalRows = $importResults.totalRows;
    $: successCount = $importResults.successCount;
    $: failureCount = $importResults.failureCount;
    $: errors = $importResults.errors || [];
    $: progress = totalRows > 0 ? Math.round((successCount + failureCount) / totalRows * 100) : 0;
    $: isComplete = !$executeLoading && (successCount + failureCount) === totalRows && totalRows > 0;
    $: hasErrors = errors.length > 0;

    let showErrors = false;

    function toggleErrors() {
        showErrors = !showErrors;
    }

    function handleStartOver() {
        dispatch("startOver");
    }

    function handleViewCollection() {
        dispatch("viewCollection");
    }
</script>

<div class="import-progress">
    {#if $executeLoading || $importInProgress}
        <div class="progress-container">
            <div class="progress-icon">
                <i class="ri-loader-4-line spinning"></i>
            </div>
            <h3>Importing records...</h3>
            <div class="progress-bar-container">
                <div class="progress-bar" style="width: {progress}%"></div>
            </div>
            <div class="progress-stats">
                <span class="stat success">
                    <i class="ri-check-line"></i>
                    {successCount} imported
                </span>
                <span class="stat error">
                    <i class="ri-close-line"></i>
                    {failureCount} failed
                </span>
                <span class="stat total">
                    of {totalRows} total
                </span>
            </div>
        </div>
    {:else if isComplete}
        <div class="complete-container" class:has-errors={hasErrors}>
            <div class="complete-icon" class:success={!hasErrors} class:warning={hasErrors}>
                {#if hasErrors}
                    <i class="ri-error-warning-line"></i>
                {:else}
                    <i class="ri-check-double-line"></i>
                {/if}
            </div>
            
            <h3>
                {#if hasErrors}
                    Import completed with errors
                {:else}
                    Import completed successfully!
                {/if}
            </h3>

            <div class="result-summary">
                <div class="result-card success">
                    <span class="result-value">{successCount}</span>
                    <span class="result-label">Records imported</span>
                </div>
                {#if hasErrors}
                    <div class="result-card error">
                        <span class="result-value">{failureCount}</span>
                        <span class="result-label">Records failed</span>
                    </div>
                {/if}
            </div>

            {#if hasErrors}
                <div class="errors-section">
                    <button 
                        type="button" 
                        class="btn btn-sm btn-secondary"
                        on:click={toggleErrors}
                    >
                        <i class={showErrors ? "ri-eye-off-line" : "ri-eye-line"}></i>
                        {showErrors ? "Hide" : "Show"} error details ({errors.length})
                    </button>

                    {#if showErrors}
                        <div class="errors-list">
                            {#each errors as error, index}
                                <div class="error-item">
                                    <div class="error-header">
                                        <span class="error-row">Row {error.row}</span>
                                        <span class="error-message">{error.message}</span>
                                    </div>
                                    {#if error.data}
                                        <pre class="error-data">{JSON.stringify(error.data, null, 2)}</pre>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}

            <div class="complete-actions">
                <button 
                    type="button" 
                    class="btn btn-secondary"
                    on:click={handleStartOver}
                >
                    <i class="ri-restart-line"></i>
                    Import more data
                </button>
                <button 
                    type="button" 
                    class="btn btn-primary"
                    on:click={handleViewCollection}
                >
                    <i class="ri-database-2-line"></i>
                    View collection
                </button>
            </div>
        </div>
    {:else if $importError}
        <div class="error-container">
            <div class="error-icon">
                <i class="ri-error-warning-line"></i>
            </div>
            <h3>Import failed</h3>
            <p class="error-text">{$importError}</p>
            <button 
                type="button" 
                class="btn btn-secondary"
                on:click={handleStartOver}
            >
                <i class="ri-restart-line"></i>
                Try again
            </button>
        </div>
    {:else}
        <div class="ready-container">
            <div class="ready-icon">
                <i class="ri-upload-2-line"></i>
            </div>
            <h3>Ready to import</h3>
            <p>Click "Start Import" to begin importing your data.</p>
        </div>
    {/if}
</div>

<style>
    .import-progress {
        width: 100%;
    }

    .progress-container,
    .complete-container,
    .error-container,
    .ready-container {
        text-align: center;
        padding: 40px 20px;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .progress-icon,
    .complete-icon,
    .error-icon,
    .ready-icon {
        width: 64px;
        height: 64px;
        margin: 0 auto 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        font-size: 32px;
    }

    .progress-icon {
        background: rgba(var(--primaryColorRaw), 0.1);
        color: var(--primaryColor);
    }

    .complete-icon.success {
        background: rgba(var(--successColorRaw), 0.1);
        color: var(--successColor);
    }

    .complete-icon.warning {
        background: rgba(var(--warningColorRaw), 0.1);
        color: var(--warningColor);
    }

    .error-icon {
        background: rgba(var(--dangerColorRaw), 0.1);
        color: var(--dangerColor);
    }

    .ready-icon {
        background: rgba(var(--primaryColorRaw), 0.1);
        color: var(--primaryColor);
    }

    .spinning {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    h3 {
        margin: 0 0 16px;
        font-size: 18px;
        color: var(--txtPrimaryColor);
    }

    .progress-bar-container {
        width: 100%;
        max-width: 400px;
        height: 8px;
        margin: 0 auto 16px;
        background: var(--baseAlt2Color);
        border-radius: 4px;
        overflow: hidden;
    }

    .progress-bar {
        height: 100%;
        background: var(--primaryColor);
        transition: width 0.3s ease;
    }

    .progress-stats {
        display: flex;
        justify-content: center;
        gap: 20px;
        font-size: 14px;
    }

    .stat {
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .stat.success {
        color: var(--successColor);
    }

    .stat.error {
        color: var(--dangerColor);
    }

    .stat.total {
        color: var(--txtHintColor);
    }

    .result-summary {
        display: flex;
        justify-content: center;
        gap: 20px;
        margin-bottom: 24px;
    }

    .result-card {
        padding: 16px 32px;
        border-radius: var(--baseRadius);
        text-align: center;
    }

    .result-card.success {
        background: rgba(var(--successColorRaw), 0.1);
    }

    .result-card.error {
        background: rgba(var(--dangerColorRaw), 0.1);
    }

    .result-value {
        display: block;
        font-size: 32px;
        font-weight: 700;
    }

    .result-card.success .result-value {
        color: var(--successColor);
    }

    .result-card.error .result-value {
        color: var(--dangerColor);
    }

    .result-label {
        font-size: 13px;
        color: var(--txtHintColor);
    }

    .errors-section {
        margin-top: 20px;
    }

    .errors-list {
        margin-top: 12px;
        max-height: 300px;
        overflow-y: auto;
        text-align: left;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .error-item {
        padding: 12px;
        border-bottom: 1px solid var(--baseAlt2Color);
    }

    .error-item:last-child {
        border-bottom: none;
    }

    .error-header {
        display: flex;
        gap: 10px;
        align-items: flex-start;
    }

    .error-row {
        flex-shrink: 0;
        padding: 2px 8px;
        background: rgba(var(--dangerColorRaw), 0.1);
        color: var(--dangerColor);
        font-size: 11px;
        font-weight: 600;
        border-radius: 4px;
    }

    .error-message {
        font-size: 13px;
        color: var(--txtPrimaryColor);
    }

    .error-data {
        margin: 8px 0 0;
        padding: 8px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        font-size: 11px;
        overflow-x: auto;
    }

    .complete-actions {
        display: flex;
        justify-content: center;
        gap: 12px;
        margin-top: 24px;
    }

    .error-text {
        margin: 0 0 20px;
        color: var(--txtHintColor);
    }
</style>
