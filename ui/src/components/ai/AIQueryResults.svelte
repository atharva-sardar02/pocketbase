<script>
    import { link } from "svelte-spa-router";
    import { aiResults, aiTotalItems, aiCollection, aiPage, aiPerPage } from "@/stores/ai";
    import CommonHelper from "@/utils/CommonHelper";

    $: hasResults = $aiResults && $aiResults.length > 0;
    $: startIndex = ($aiPage - 1) * $aiPerPage + 1;
    $: endIndex = Math.min($aiPage * $aiPerPage, $aiTotalItems);
</script>

{#if hasResults}
    <div class="ai-query-results">
        <div class="results-header">
            <h5 class="m-0">
                Results
                {#if $aiTotalItems > 0}
                    <span class="results-count">
                        ({startIndex}-{endIndex} of {$aiTotalItems})
                    </span>
                {/if}
            </h5>
            {#if $aiCollection}
                <a
                    href="/collections?collection={$aiCollection}"
                    class="btn btn-sm btn-transparent"
                    use:link
                >
                    <i class="ri-external-link-line" aria-hidden="true"></i>
                    <span class="txt">View in Collection</span>
                </a>
            {/if}
        </div>
        <div class="results-list">
            {#each $aiResults as record}
                <div class="result-item">
                    <div class="result-id">
                        <strong>ID:</strong> {record.id}
                    </div>
                    <div class="result-fields">
                        {#each Object.entries(record).slice(0, 5) as [key, value]}
                            {#if key !== "id" && value !== null && value !== undefined}
                                <div class="result-field">
                                    <strong>{key}:</strong> {CommonHelper.displayValue(value)}
                                </div>
                            {/if}
                        {/each}
                    </div>
                </div>
            {/each}
        </div>
    </div>
{:else if $aiTotalItems === 0 && $aiResults.length === 0}
    <div class="ai-query-results empty">
        <div class="empty-state">
            <i class="ri-inbox-line" aria-hidden="true"></i>
            <p>No results found</p>
        </div>
    </div>
{/if}

<style>
    .ai-query-results {
        margin: 15px;
        border: 1px solid var(--borderColor);
        border-radius: 4px;
        background: var(--baseColor);
    }

    .results-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 15px;
        border-bottom: 1px solid var(--borderColor);
    }

    .results-count {
        font-weight: normal;
        color: var(--txtHintColor);
        font-size: 0.9em;
    }

    .results-list {
        padding: 15px;
        max-height: 400px;
        overflow-y: auto;
    }

    .result-item {
        padding: 12px;
        margin-bottom: 10px;
        border: 1px solid var(--borderColor);
        border-radius: 4px;
        background: var(--baseAlt1Color);
    }

    .result-item:last-child {
        margin-bottom: 0;
    }

    .result-id {
        margin-bottom: 8px;
        color: var(--txtPrimaryColor);
    }

    .result-fields {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
    }

    .result-field {
        font-size: 0.9em;
        color: var(--txtSecondaryColor);
    }

    .result-field strong {
        color: var(--txtPrimaryColor);
        margin-right: 4px;
    }

    .empty-state {
        padding: 40px;
        text-align: center;
        color: var(--txtHintColor);
    }

    .empty-state i {
        font-size: 48px;
        margin-bottom: 12px;
        opacity: 0.5;
    }
</style>

