<script>
    export let activeTab = "filter";
    export let filterDisabled = false;
    export let sqlDisabled = false;

    function setTab(tab) {
        if (tab === "filter" && !filterDisabled) {
            activeTab = "filter";
        } else if (tab === "sql" && !sqlDisabled) {
            activeTab = "sql";
        }
    }
</script>

<div class="query-tabs">
    <button 
        class="tab-btn" 
        class:active={activeTab === "filter"}
        class:disabled={filterDisabled}
        on:click={() => setTab("filter")}
        disabled={filterDisabled}
        title={filterDisabled ? "Filter not available for this query" : "PocketBase Filter"}
    >
        <i class="ri-filter-line"></i>
        <span>Filter</span>
    </button>
    <button 
        class="tab-btn" 
        class:active={activeTab === "sql"}
        class:disabled={sqlDisabled}
        on:click={() => setTab("sql")}
        disabled={sqlDisabled}
        title={sqlDisabled ? "SQL not available" : "SQL Query"}
    >
        <i class="ri-database-2-line"></i>
        <span>SQL</span>
    </button>
</div>

<style>
    .query-tabs {
        display: flex;
        gap: 0;
        margin-bottom: 10px;
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .tab-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        border: none;
        background: transparent;
        color: var(--txtHintColor);
        cursor: pointer;
        font-size: 13px;
        font-weight: 500;
        border-bottom: 2px solid transparent;
        transition: all 0.2s ease;
        margin-bottom: -1px;
    }
    .tab-btn:hover:not(.disabled) {
        color: var(--txtPrimaryColor);
    }
    .tab-btn.active {
        color: var(--primaryColor);
        border-bottom-color: var(--primaryColor);
    }
    .tab-btn.disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }
    .tab-btn i {
        font-size: 16px;
    }
</style>

