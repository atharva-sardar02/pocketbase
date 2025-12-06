<script>
    import { createEventDispatcher } from "svelte";
    import { link } from "svelte-spa-router";
    import { collections } from "@/stores/collections";
    import { aiQuery, aiCollection, aiLoading } from "@/stores/ai";
    import AutoExpandTextarea from "@/components/base/AutoExpandTextarea.svelte";
    import Select from "@/components/base/Select.svelte";

    const dispatch = createEventDispatcher();

    let queryText = "";
    let selectedCollection = "";

    $: if ($collections && $collections.length > 0 && !selectedCollection) {
        // Default to first non-system collection
        const firstCollection = $collections.find((c) => !c.system) || $collections[0];
        selectedCollection = firstCollection?.name || "";
    }

    $: collectionOptions = $collections
        ? $collections.map((c) => ({
              value: c.name,
              label: c.name,
          }))
        : [];

    function handleSubmit() {
        if (!$aiLoading && queryText.trim() && selectedCollection) {
            aiQuery.set(queryText.trim());
            aiCollection.set(selectedCollection);
            dispatch("submit", {
                query: queryText.trim(),
                collection: selectedCollection,
            });
        }
    }

    function handleKeydown(e) {
        // Ctrl+Enter or Cmd+Enter to submit
        if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
            e.preventDefault();
            handleSubmit();
        }
    }
</script>

<div class="ai-query-input">
    <div class="form-field m-b-sm">
        <label class="form-label">Collection</label>
        {#if $collections && $collections.length > 0}
            <Select
                options={collectionOptions}
                value={selectedCollection}
                on:change={(e) => {
                    selectedCollection = e.detail.value;
                    aiCollection.set(selectedCollection);
                }}
                disabled={$aiLoading}
            />
        {:else}
            <div class="alert alert-info">
                <i class="ri-information-line" aria-hidden="true"></i>
                <div class="content">
                    <p class="m-0">
                        No collections found. Please
                        <a href="/collections" use:link>create a collection</a> first.
                    </p>
                </div>
            </div>
        {/if}
    </div>

    <div class="form-field m-b-sm">
        <label class="form-label">Query</label>
        <AutoExpandTextarea
            value={queryText}
            placeholder="Enter your query in natural language... (e.g., 'show me active orders over $100')"
            on:input={(e) => (queryText = e.target.value)}
            on:keydown={handleKeydown}
            disabled={$aiLoading}
            rows="3"
        />
        <div class="form-hint">
            Press <kbd>Ctrl+Enter</kbd> or <kbd>Cmd+Enter</kbd> to search
        </div>
    </div>

    <button
        type="button"
        class="btn btn-primary"
        on:click={handleSubmit}
        disabled={$aiLoading || !queryText.trim() || !selectedCollection}
    >
        {#if $aiLoading}
            <i class="ri-loader-4-line spin" aria-hidden="true"></i>
            <span class="txt">Searching...</span>
        {:else}
            <i class="ri-search-line" aria-hidden="true"></i>
            <span class="txt">Search</span>
        {/if}
    </button>
</div>

<style>
    .ai-query-input {
        padding: 15px;
    }

    kbd {
        background: var(--baseAlt1Color);
        border: 1px solid var(--borderColor);
        border-radius: 3px;
        padding: 2px 6px;
        font-size: 0.85em;
        font-family: monospace;
    }
</style>

