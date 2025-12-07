<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let schema = [];
    export let loading = false;

    const dispatch = createEventDispatcher();

    let expandedCollections = {};
    let searchQuery = "";

    $: filteredSchema = filterSchema(schema, searchQuery);

    function filterSchema(collections, query) {
        if (!query) return collections;
        const lower = query.toLowerCase();
        return collections.filter(c => 
            c.name.toLowerCase().includes(lower) ||
            c.fields?.some(f => f.name.toLowerCase().includes(lower))
        );
    }

    function toggleCollection(name) {
        expandedCollections[name] = !expandedCollections[name];
    }

    function insertTableName(name) {
        dispatch("insert", { text: name });
    }

    function insertFieldName(tableName, fieldName) {
        dispatch("insert", { text: `${tableName}.${fieldName}` });
    }

    function getFieldIcon(type) {
        const icons = {
            text: "ri-text",
            number: "ri-hashtag",
            bool: "ri-toggle-line",
            date: "ri-calendar-line",
            select: "ri-list-check",
            file: "ri-file-line",
            relation: "ri-link",
            json: "ri-code-s-slash-line",
            url: "ri-links-line",
            email: "ri-mail-line",
            editor: "ri-article-line",
            autodate: "ri-time-line",
        };
        return icons[type] || "ri-question-line";
    }

    function refresh() {
        dispatch("refresh");
    }
</script>

<div class="schema-explorer">
    <div class="explorer-header">
        <span class="explorer-title">Schema</span>
        <button
            type="button"
            class="btn btn-xs btn-transparent"
            on:click={refresh}
            disabled={loading}
            use:tooltip={"Refresh schema"}
        >
            <i class="ri-refresh-line" class:animate-spin={loading}></i>
        </button>
    </div>

    <div class="explorer-search">
        <div class="search-input-wrapper">
            <i class="ri-search-line"></i>
            <input
                type="text"
                bind:value={searchQuery}
                placeholder="Filter..."
                class="search-input"
            />
            {#if searchQuery}
                <button
                    type="button"
                    class="clear-btn"
                    on:click={() => searchQuery = ""}
                >
                    <i class="ri-close-line"></i>
                </button>
            {/if}
        </div>
    </div>

    <div class="explorer-body">
        {#if loading}
            <div class="loading-state">
                <i class="ri-loader-4-line animate-spin"></i>
                <span>Loading schema...</span>
            </div>
        {:else if filteredSchema.length === 0}
            <div class="empty-state">
                {#if searchQuery}
                    <span>No collections match "{searchQuery}"</span>
                {:else}
                    <span>No collections found</span>
                {/if}
            </div>
        {:else}
            <ul class="collection-list">
                {#each filteredSchema as collection}
                    <li class="collection-item">
                        <button
                            type="button"
                            class="collection-header"
                            on:click={() => toggleCollection(collection.name)}
                        >
                            <i class="ri-arrow-right-s-line expand-icon" class:expanded={expandedCollections[collection.name]}></i>
                            <i class="ri-table-line collection-icon"></i>
                            <span class="collection-name" on:dblclick|stopPropagation={() => insertTableName(collection.name)}>
                                {collection.name}
                            </span>
                            {#if collection.type === "auth"}
                                <span class="collection-badge auth">auth</span>
                            {:else if collection.type === "view"}
                                <span class="collection-badge view">view</span>
                            {/if}
                        </button>

                        {#if expandedCollections[collection.name] && collection.fields}
                            <ul class="field-list">
                                {#each collection.fields as field}
                                    <li class="field-item">
                                        <button
                                            type="button"
                                            class="field-row"
                                            on:click={() => insertFieldName(collection.name, field.name)}
                                            use:tooltip={`Click to insert ${collection.name}.${field.name}`}
                                        >
                                            <i class={getFieldIcon(field.type)}></i>
                                            <span class="field-name">{field.name}</span>
                                            <span class="field-type">{field.type}</span>
                                        </button>
                                    </li>
                                {/each}
                            </ul>
                        {/if}
                    </li>
                {/each}
            </ul>
        {/if}
    </div>

    <div class="explorer-footer">
        <span class="hint">Double-click to insert</span>
    </div>
</div>

<style>
    .schema-explorer {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        overflow: hidden;
    }
    .explorer-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px 12px;
        background: var(--baseAlt1Color);
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .explorer-title {
        font-size: 12px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
    }
    .explorer-search {
        padding: 8px;
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .search-input-wrapper {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 6px 10px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
    }
    .search-input-wrapper i {
        color: var(--txtDisabledColor);
        font-size: 14px;
    }
    .search-input {
        flex: 1;
        border: none;
        background: transparent;
        font-size: 13px;
        color: var(--txtPrimaryColor);
        outline: none;
    }
    .search-input::placeholder {
        color: var(--txtDisabledColor);
    }
    .clear-btn {
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        color: var(--txtHintColor);
    }
    .clear-btn:hover {
        color: var(--txtPrimaryColor);
    }
    .explorer-body {
        flex: 1;
        overflow-y: auto;
        padding: 8px 0;
    }
    .loading-state,
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 30px 15px;
        color: var(--txtHintColor);
        font-size: 13px;
    }
    .collection-list {
        list-style: none;
        margin: 0;
        padding: 0;
    }
    .collection-item {
        border-bottom: 1px solid var(--baseAlt1Color);
    }
    .collection-item:last-child {
        border-bottom: none;
    }
    .collection-header {
        display: flex;
        align-items: center;
        gap: 6px;
        width: 100%;
        padding: 8px 12px;
        background: none;
        border: none;
        cursor: pointer;
        text-align: left;
        color: var(--txtPrimaryColor);
        font-size: 13px;
    }
    .collection-header:hover {
        background: var(--baseAlt1Color);
    }
    .expand-icon {
        transition: transform 0.15s;
        color: var(--txtHintColor);
    }
    .expand-icon.expanded {
        transform: rotate(90deg);
    }
    .collection-icon {
        color: var(--primaryColor);
    }
    .collection-name {
        flex: 1;
        font-weight: 500;
    }
    .collection-badge {
        font-size: 10px;
        padding: 1px 5px;
        border-radius: 3px;
        text-transform: uppercase;
        font-weight: 600;
    }
    .collection-badge.auth {
        background: var(--successAltColor);
        color: var(--successColor);
    }
    .collection-badge.view {
        background: var(--infoAltColor);
        color: var(--infoColor);
    }
    .field-list {
        list-style: none;
        margin: 0;
        padding: 0 0 8px 0;
        background: var(--baseAlt1Color);
    }
    .field-item {
        margin: 0;
    }
    .field-row {
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
        padding: 5px 12px 5px 32px;
        background: none;
        border: none;
        cursor: pointer;
        text-align: left;
        font-size: 12px;
        color: var(--txtSecondaryColor);
    }
    .field-row:hover {
        background: var(--baseAlt2Color);
        color: var(--txtPrimaryColor);
    }
    .field-row i {
        font-size: 12px;
        color: var(--txtHintColor);
    }
    .field-name {
        flex: 1;
    }
    .field-type {
        font-size: 10px;
        color: var(--txtDisabledColor);
        text-transform: uppercase;
    }
    .explorer-footer {
        padding: 6px 12px;
        background: var(--baseAlt1Color);
        border-top: 1px solid var(--baseAlt2Color);
    }
    .hint {
        font-size: 11px;
        color: var(--txtHintColor);
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>

