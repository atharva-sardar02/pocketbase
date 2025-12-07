<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";
    import { sqlHistory, clearHistory } from "@/stores/sql";

    const dispatch = createEventDispatcher();

    let isOpen = false;
    let searchQuery = "";

    $: filteredHistory = filterHistory($sqlHistory, searchQuery);

    function filterHistory(history, query) {
        if (!query) return history;
        const lower = query.toLowerCase();
        return history.filter(h => h.query.toLowerCase().includes(lower));
    }

    function selectQuery(item) {
        dispatch("select", { query: item.query, type: item.type });
        isOpen = false;
    }

    function formatTime(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        const diff = now - date;
        
        if (diff < 60000) return "just now";
        if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
        if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
        return date.toLocaleDateString();
    }

    function truncateQuery(query, maxLength = 60) {
        const singleLine = query.replace(/\s+/g, " ").trim();
        if (singleLine.length <= maxLength) return singleLine;
        return singleLine.substring(0, maxLength) + "...";
    }

    function handleClear() {
        if (confirm("Clear all query history?")) {
            clearHistory();
        }
    }

    function handleClickOutside(event) {
        if (isOpen && !event.target.closest(".query-history")) {
            isOpen = false;
        }
    }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="query-history">
    <button
        type="button"
        class="history-toggle"
        class:active={isOpen}
        on:click|stopPropagation={() => isOpen = !isOpen}
        use:tooltip={"Query history"}
        disabled={$sqlHistory.length === 0}
    >
        <i class="ri-history-line"></i>
        History
        {#if $sqlHistory.length > 0}
            <span class="badge">{$sqlHistory.length}</span>
        {/if}
        <i class="ri-arrow-down-s-line"></i>
    </button>

    {#if isOpen}
        <div class="history-dropdown" on:click|stopPropagation>
            <div class="dropdown-header">
                <input
                    type="text"
                    bind:value={searchQuery}
                    placeholder="Search history..."
                    class="search-input"
                />
                <button
                    type="button"
                    class="btn btn-xs btn-transparent btn-danger"
                    on:click={handleClear}
                    use:tooltip={"Clear history"}
                >
                    <i class="ri-delete-bin-line"></i>
                </button>
            </div>

            <div class="dropdown-body">
                {#if filteredHistory.length === 0}
                    <div class="empty-state">
                        {#if searchQuery}
                            No matches found
                        {:else}
                            No history yet
                        {/if}
                    </div>
                {:else}
                    <ul class="history-list">
                        {#each filteredHistory as item}
                            <li>
                                <button
                                    type="button"
                                    class="history-item"
                                    on:click={() => selectQuery(item)}
                                >
                                    <div class="item-query">
                                        <code>{truncateQuery(item.query)}</code>
                                    </div>
                                    <div class="item-meta">
                                        <span class="item-type">{item.type || "sql"}</span>
                                        <span class="item-time">{formatTime(item.timestamp)}</span>
                                    </div>
                                </button>
                            </li>
                        {/each}
                    </ul>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .query-history {
        position: relative;
        display: inline-block;
    }
    .history-toggle {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 6px 12px;
        background: var(--baseAlt1Color);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        color: var(--txtSecondaryColor);
        font-size: 13px;
        cursor: pointer;
        transition: all 0.15s;
    }
    .history-toggle:hover:not(:disabled) {
        background: var(--baseAlt2Color);
        color: var(--txtPrimaryColor);
    }
    .history-toggle:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    .history-toggle.active {
        background: var(--baseAlt2Color);
        border-color: var(--primaryColor);
    }
    .badge {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 18px;
        height: 18px;
        padding: 0 5px;
        font-size: 11px;
        font-weight: 600;
        background: var(--primaryColor);
        color: var(--baseColor);
        border-radius: 9px;
    }
    .history-dropdown {
        position: absolute;
        top: 100%;
        left: 0;
        z-index: 100;
        width: 400px;
        max-width: 90vw;
        margin-top: 4px;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    .dropdown-header {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px;
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .search-input {
        flex: 1;
        padding: 6px 10px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseAlt1Color);
        font-size: 13px;
        color: var(--txtPrimaryColor);
    }
    .search-input:focus {
        outline: none;
        border-color: var(--primaryColor);
    }
    .dropdown-body {
        max-height: 300px;
        overflow-y: auto;
    }
    .empty-state {
        padding: 30px 20px;
        text-align: center;
        color: var(--txtHintColor);
        font-size: 13px;
    }
    .history-list {
        list-style: none;
        margin: 0;
        padding: 0;
    }
    .history-item {
        display: block;
        width: 100%;
        padding: 10px 12px;
        background: none;
        border: none;
        border-bottom: 1px solid var(--baseAlt1Color);
        text-align: left;
        cursor: pointer;
    }
    .history-item:hover {
        background: var(--baseAlt1Color);
    }
    .history-item:last-child {
        border-bottom: none;
    }
    .item-query {
        margin-bottom: 4px;
    }
    .item-query code {
        font-family: var(--monospaceFontFamily);
        font-size: 12px;
        color: var(--txtPrimaryColor);
        word-break: break-all;
    }
    .item-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 11px;
        color: var(--txtHintColor);
    }
    .item-type {
        text-transform: uppercase;
        font-weight: 600;
    }
</style>

