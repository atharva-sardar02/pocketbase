<script>
    import { link } from "svelte-spa-router";
    
    export let data = [];
    export let loading = false;

    function formatNumber(num) {
        if (num === null || num === undefined) return "0";
        return num.toLocaleString();
    }

    function getTypeIcon(type) {
        switch (type) {
            case "auth":
                return "ri-user-line";
            case "view":
                return "ri-eye-line";
            default:
                return "ri-table-line";
        }
    }

    function getTypeBadgeClass(type) {
        switch (type) {
            case "auth":
                return "badge-auth";
            case "view":
                return "badge-view";
            default:
                return "badge-base";
        }
    }
</script>

<div class="collections-table-wrapper" class:loading>
    {#if loading}
        <div class="table-loader">
            <div class="loader"></div>
        </div>
    {/if}
    
    <table class="collections-table">
        <thead>
            <tr>
                <th>Collection</th>
                <th>Type</th>
                <th class="text-right">Records</th>
                <th class="text-right">Actions</th>
            </tr>
        </thead>
        <tbody>
            {#if data.length === 0 && !loading}
                <tr>
                    <td colspan="4" class="empty-state">
                        <i class="ri-database-2-line"></i>
                        <span>No collections found</span>
                    </td>
                </tr>
            {:else}
                {#each data as item}
                    <tr>
                        <td class="collection-name">
                            <i class={getTypeIcon(item.type)}></i>
                            {item.name}
                        </td>
                        <td>
                            <span class="type-badge {getTypeBadgeClass(item.type)}">
                                {item.type}
                            </span>
                        </td>
                        <td class="text-right record-count">
                            {formatNumber(item.recordCount)}
                        </td>
                        <td class="text-right">
                            <a 
                                href="/collections?collection={item.name}" 
                                class="btn btn-sm btn-secondary"
                                use:link
                            >
                                <i class="ri-arrow-right-line"></i>
                            </a>
                        </td>
                    </tr>
                {/each}
            {/if}
        </tbody>
    </table>
</div>

<style>
    .collections-table-wrapper {
        position: relative;
        overflow: auto;
        max-height: 400px;
    }
    .collections-table-wrapper.loading table {
        opacity: 0.5;
    }
    .table-loader {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        z-index: 10;
    }
    .collections-table {
        width: 100%;
        border-collapse: collapse;
    }
    .collections-table th,
    .collections-table td {
        padding: 12px 16px;
        text-align: left;
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .collections-table th {
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
        background: var(--baseAlt1Color);
        position: sticky;
        top: 0;
        z-index: 1;
    }
    .collections-table tbody tr:hover {
        background: var(--baseAlt1Color);
    }
    .collection-name {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 500;
    }
    .collection-name i {
        color: var(--txtHintColor);
    }
    .type-badge {
        display: inline-block;
        padding: 3px 8px;
        font-size: 11px;
        font-weight: 500;
        border-radius: 4px;
        text-transform: uppercase;
    }
    .badge-auth {
        background: rgba(var(--successColorRaw), 0.1);
        color: var(--successColor);
    }
    .badge-view {
        background: rgba(var(--infoColorRaw), 0.1);
        color: var(--infoColor);
    }
    .badge-base {
        background: var(--baseAlt1Color);
        color: var(--txtSecondaryColor);
    }
    .record-count {
        font-family: var(--monospaceFontFamily);
        font-weight: 500;
    }
    .text-right {
        text-align: right;
    }
    .empty-state {
        text-align: center;
        padding: 40px !important;
        color: var(--txtHintColor);
    }
    .empty-state i {
        font-size: 24px;
        margin-bottom: 8px;
        display: block;
    }
</style>
