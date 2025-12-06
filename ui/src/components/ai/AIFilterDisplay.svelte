<script>
    import { link } from "svelte-spa-router";
    import { aiFilter, aiCollection } from "@/stores/ai";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import { addSuccessToast } from "@/stores/toasts";

    function copyFilter() {
        if ($aiFilter) {
            navigator.clipboard.writeText($aiFilter);
            addSuccessToast("Filter copied to clipboard");
        }
    }

    function applyFilter() {
        if ($aiFilter && $aiCollection) {
            // Navigate to collection with filter applied
            const filterParam = encodeURIComponent($aiFilter);
            window.location.href = `/collections?collection=${$aiCollection}&filter=${filterParam}`;
        }
    }
</script>

{#if $aiFilter}
    <div class="ai-filter-display">
        <div class="filter-header">
            <h5 class="m-0">Generated Filter</h5>
            <div class="filter-actions">
                <button type="button" class="btn btn-sm btn-transparent" on:click={copyFilter}>
                    <i class="ri-file-copy-line" aria-hidden="true"></i>
                    <span class="txt">Copy</span>
                </button>
                <button type="button" class="btn btn-sm btn-primary" on:click={applyFilter}>
                    <i class="ri-arrow-right-line" aria-hidden="true"></i>
                    <span class="txt">Apply Filter</span>
                </button>
            </div>
        </div>
        <div class="filter-content">
            <CodeBlock content={$aiFilter} language="javascript" />
        </div>
    </div>
{/if}

<style>
    .ai-filter-display {
        margin: 15px;
        border: 1px solid var(--borderColor);
        border-radius: 4px;
        background: var(--baseColor);
    }

    .filter-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 15px;
        border-bottom: 1px solid var(--borderColor);
    }

    .filter-actions {
        display: flex;
        gap: 8px;
    }

    .filter-content {
        padding: 0;
    }
</style>

