<script>
    import { onMount } from "svelte";
    import { pageTitle } from "@/stores/app";
    import {
        aiQuery,
        aiFilter,
        aiSQL,
        aiRequiresSQL,
        aiCanUseFilter,
        aiActiveTab,
        aiResults,
        aiTotalItems,
        aiPage,
        aiPerPage,
        aiLoading,
        aiError,
        aiCollection,
        resetAIState,
        setDualResponse,
    } from "@/stores/ai";
    import { loadCollections, collections, changeActiveCollectionByIdOrName } from "@/stores/collections";
    import { push } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import AIQueryInput from "./AIQueryInput.svelte";
    import AIFilterDisplay from "./AIFilterDisplay.svelte";
    import AIQueryResults from "./AIQueryResults.svelte";
    import QueryTabs from "./QueryTabs.svelte";
    import EditableCodeBlock from "./EditableCodeBlock.svelte";

    $pageTitle = "AI Query";

    let executeQuery = true; // Whether to execute the filter and return results
    let useDualMode = true; // V2: Use dual output mode

    onMount(() => {
        // Reset state when component mounts
        resetAIState();
        // Load collections for the dropdown
        loadCollections();
    });

    async function handleQuerySubmit(event) {
        const { query, collection } = event.detail;

        if (!query || !collection) {
            return;
        }

        aiLoading.set(true);
        aiError.set("");
        aiFilter.set("");
        aiSQL.set("");
        aiResults.set([]);
        aiTotalItems.set(0);

        try {
            const requestBody = {
                collection: collection,
                query: query,
                execute: executeQuery,
                page: $aiPage,
                perPage: $aiPerPage,
                mode: useDualMode ? "dual" : "filter", // V2: Use dual mode
            };

            // Make API call to /api/ai/query
            // Ensure proper URL construction (avoid double slashes)
            const baseUrl = ApiClient.baseURL.endsWith('/') 
                ? ApiClient.baseURL.slice(0, -1) 
                : ApiClient.baseURL;
            const url = `${baseUrl}/api/ai/query`;
            const response = await fetch(url, {
                method: "POST",
                body: JSON.stringify(requestBody),
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token ? `Bearer ${ApiClient.authStore.token}` : "",
                },
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            // V2: Set dual response data
            setDualResponse(data);
            aiError.set("");
        } catch (err) {
            ApiClient.error(err, true, "Failed to process AI query");
            aiError.set(err.message || "An error occurred while processing your query.");
            aiFilter.set("");
            aiSQL.set("");
            aiResults.set([]);
            aiTotalItems.set(0);
        } finally {
            aiLoading.set(false);
        }
    }

    // Handle filter edit and re-execute
    async function handleFilterExecute(event) {
        const { code } = event.detail;
        if (!code || !$aiCollection) return;

        aiLoading.set(true);
        aiError.set("");

        try {
            // Execute the edited filter directly
            const baseUrl = ApiClient.baseURL.endsWith('/') 
                ? ApiClient.baseURL.slice(0, -1) 
                : ApiClient.baseURL;
            
            // Use the records API directly
            const url = `${baseUrl}/api/collections/${$aiCollection}/records?filter=${encodeURIComponent(code)}&page=${$aiPage}&perPage=${$aiPerPage}`;
            const response = await fetch(url, {
                method: "GET",
                headers: {
                    Authorization: ApiClient.authStore.token ? `Bearer ${ApiClient.authStore.token}` : "",
                },
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            aiFilter.set(code);
            aiResults.set(data.items || []);
            aiTotalItems.set(data.totalItems || 0);
            aiError.set("");
        } catch (err) {
            aiError.set(err.message || "Failed to execute filter.");
        } finally {
            aiLoading.set(false);
        }
    }

    // Handle SQL edit and re-execute
    async function handleSQLExecute(event) {
        const { code } = event.detail;
        if (!code) return;

        aiLoading.set(true);
        aiError.set("");

        try {
            const baseUrl = ApiClient.baseURL.endsWith('/') 
                ? ApiClient.baseURL.slice(0, -1) 
                : ApiClient.baseURL;
            const url = `${baseUrl}/api/sql/execute`;
            const response = await fetch(url, {
                method: "POST",
                body: JSON.stringify({ sql: code, confirm: true }),
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token ? `Bearer ${ApiClient.authStore.token}` : "",
                },
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            aiSQL.set(code);
            
            if (data.rows) {
                // Convert SQL results to display format
                aiResults.set(data.rows);
                aiTotalItems.set(data.totalRows || data.rows.length);
            }
            aiError.set("");
        } catch (err) {
            aiError.set(err.message || "Failed to execute SQL.");
        } finally {
            aiLoading.set(false);
        }
    }

    function handleCopy() {
        // Could show a toast notification here
    }

    // Get collection ID from name
    $: collectionId = $collections?.find(c => c.name === $aiCollection)?.id || $aiCollection;

    function seeInCollection() {
        if ($aiFilter && collectionId) {
            const filterParam = encodeURIComponent($aiFilter);
            changeActiveCollectionByIdOrName(collectionId);
            push(`/collections?collection=${collectionId}&filter=${filterParam}`);
        }
    }
</script>

<PageWrapper>
    <div class="page-layout">
        <PageSidebar>
            <div class="sidebar-content">
                <div class="sidebar-title">AI Query</div>
                <div class="sidebar-description">
                    Query your collections using natural language. The AI will generate both PocketBase filter
                    and SQL query for you.
                </div>
                <div class="sidebar-options">
                    <label class="form-field form-field-toggle">
                        <input type="checkbox" bind:checked={useDualMode} />
                        <span class="txt">Dual output (Filter + SQL)</span>
                    </label>
                </div>
            </div>
        </PageSidebar>

        <div class="page-content">
            <div class="page-header">
                <h1>AI Query Assistant</h1>
                <p class="page-description">
                    Describe what you're looking for in plain English. Get both PocketBase filter and SQL query.
                </p>
            </div>

            <div class="page-body">
                <div class="ai-query-container">
                    <div class="query-section">
                        <AIQueryInput on:submit={handleQuerySubmit} />
                    </div>

                    {#if $aiError}
                        <div class="alert alert-error m-t-sm">
                            <i class="ri-error-warning-line" aria-hidden="true"></i>
                            <span class="txt">{$aiError}</span>
                        </div>
                    {/if}

                    {#if $aiFilter || $aiSQL}
                        <div class="output-section m-t-md">
                            {#if useDualMode && ($aiFilter || $aiSQL)}
                                <QueryTabs 
                                    bind:activeTab={$aiActiveTab}
                                    filterDisabled={!$aiCanUseFilter}
                                    sqlDisabled={!$aiSQL}
                                />
                            {/if}

                            {#if $aiActiveTab === "filter" && $aiFilter}
                                <EditableCodeBlock 
                                    code={$aiFilter}
                                    language="filter"
                                    on:execute={handleFilterExecute}
                                    on:copy={handleCopy}
                                    on:change={(e) => aiFilter.set(e.detail.code)}
                                />
                                <div class="filter-actions m-t-xs">
                                    <button 
                                        type="button" 
                                        class="btn btn-sm btn-secondary"
                                        on:click={seeInCollection}
                                    >
                                        <i class="ri-eye-line"></i>
                                        <span class="txt">See in Collection</span>
                                    </button>
                                </div>
                                {#if $aiRequiresSQL}
                                    <div class="alert alert-warning m-t-xs">
                                        <i class="ri-information-line"></i>
                                        <span>This query may require SQL for full functionality (JOINs/aggregates).</span>
                                    </div>
                                {/if}
                            {:else if $aiActiveTab === "sql" && $aiSQL}
                                <EditableCodeBlock 
                                    code={$aiSQL}
                                    language="sql"
                                    on:execute={handleSQLExecute}
                                    on:copy={handleCopy}
                                    on:change={(e) => aiSQL.set(e.detail.code)}
                                />
                            {:else if !useDualMode && $aiFilter}
                                <div class="filter-section">
                                    <AIFilterDisplay />
                                </div>
                            {/if}
                        </div>
                    {/if}

                    {#if executeQuery && $aiResults.length > 0}
                        <div class="results-section m-t-md">
                            <AIQueryResults />
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</PageWrapper>

<style>
    .page-layout {
        display: flex;
        flex: 1;
        min-height: 0;
    }

    .sidebar-content {
        padding: 0;
    }

    .sidebar-title {
        padding: 15px 15px 5px;
        font-weight: 600;
        font-size: 1.1em;
    }

    .sidebar-description {
        padding: 0 15px 10px;
        font-size: 0.9em;
        color: var(--txtHintColor);
        line-height: 1.5;
    }

    .sidebar-options {
        padding: 10px 15px;
        border-top: 1px solid var(--baseAlt2Color);
    }

    .page-content {
        flex: 1;
        padding: 20px;
        overflow-y: auto;
    }

    .page-header {
        margin-bottom: 30px;
    }

    .page-description {
        color: var(--txtSecondaryColor);
        margin-top: 8px;
    }

    .ai-query-container {
        max-width: 900px;
    }

    .query-section {
        background: var(--baseColor);
        border: 1px solid var(--borderColor);
        border-radius: 4px;
    }

    .output-section,
    .filter-section,
    .results-section {
        margin-top: 20px;
    }

    .alert-warning {
        background: var(--warningAltColor);
        color: var(--txtPrimaryColor);
    }

    .filter-actions {
        display: flex;
        justify-content: flex-end;
        gap: 8px;
    }
</style>

