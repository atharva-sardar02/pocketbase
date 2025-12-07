<script>
    import { onMount } from "svelte";
    import { pageTitle } from "@/stores/app";
    import {
        aiQuery,
        aiFilter,
        aiResults,
        aiTotalItems,
        aiPage,
        aiPerPage,
        aiLoading,
        aiError,
        aiCollection,
        resetAIState,
    } from "@/stores/ai";
    import { loadCollections } from "@/stores/collections";
    import ApiClient from "@/utils/ApiClient";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import AIQueryInput from "./AIQueryInput.svelte";
    import AIFilterDisplay from "./AIFilterDisplay.svelte";
    import AIQueryResults from "./AIQueryResults.svelte";

    $pageTitle = "AI Query";

    let executeQuery = true; // Whether to execute the filter and return results

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
        aiResults.set([]);
        aiTotalItems.set(0);

        try {
            const requestBody = {
                collection: collection,
                query: query,
                execute: executeQuery,
                page: $aiPage,
                perPage: $aiPerPage,
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

            aiFilter.set(data.filter || "");
            if (executeQuery && data.results) {
                aiResults.set(data.results);
                aiTotalItems.set(data.totalItems || 0);
                aiPage.set(data.page || 1);
                aiPerPage.set(data.perPage || 30);
            }
            aiError.set("");
        } catch (err) {
            ApiClient.error(err, true, "Failed to process AI query");
            aiError.set(err.message || "An error occurred while processing your query.");
            aiFilter.set("");
            aiResults.set([]);
            aiTotalItems.set(0);
        } finally {
            aiLoading.set(false);
        }
    }
</script>

<PageWrapper>
    <div class="page-layout">
        <PageSidebar>
            <div class="sidebar-content">
                <div class="sidebar-title">AI Query</div>
                <div class="sidebar-description">
                    Query your collections using natural language. The AI will generate a PocketBase filter
                    expression for you.
                </div>
            </div>
        </PageSidebar>

        <div class="page-content">
            <div class="page-header">
                <h1>AI Query Assistant</h1>
                <p class="page-description">
                    Describe what you're looking for in plain English, and we'll generate a filter expression
                    for you.
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

                    {#if $aiFilter}
                        <div class="filter-section m-t-md">
                            <AIFilterDisplay />
                        </div>
                    {/if}

                    {#if executeQuery && ($aiResults.length > 0 || $aiTotalItems === 0)}
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

    .sidebar-description {
        padding: 15px;
        font-size: 0.9em;
        color: var(--txtHintColor);
        line-height: 1.5;
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

    .filter-section,
    .results-section {
        margin-top: 20px;
    }
</style>

