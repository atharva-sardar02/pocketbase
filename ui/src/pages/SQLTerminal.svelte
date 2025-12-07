<script>
    import { onMount } from "svelte";
    import { pageTitle } from "@/stores/app";
    import {
        sqlQuery,
        sqlAIQuery,
        sqlMode,
        sqlResults,
        sqlColumns,
        sqlTotalRows,
        sqlRowsAffected,
        sqlLoading,
        sqlError,
        sqlSuccess,
        sqlSchema,
        sqlSchemaLoading,
        sqlNeedsConfirmation,
        sqlConfirmMessage,
        sqlPendingQuery,
        resetSQLState,
        addToHistory,
        isDestructiveQuery,
        hasResults,
        isAIMode,
        sqlIsMulti,
        sqlMultiResults,
        sqlTotalStatements,
        sqlSuccessfulCount,
        sqlFailedCount,
    } from "@/stores/sql";
    import ApiClient from "@/utils/ApiClient";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SQLEditor from "@/components/sql/SQLEditor.svelte";
    import SchemaExplorer from "@/components/sql/SchemaExplorer.svelte";
    import ResultsTable from "@/components/sql/ResultsTable.svelte";
    import QueryHistory from "@/components/sql/QueryHistory.svelte";

    $pageTitle = "SQL Terminal";

    let editorRef;

    onMount(() => {
        resetSQLState();
        loadSchema();
    });

    async function loadSchema() {
        sqlSchemaLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            const response = await fetch(`${baseUrl}/api/sql/schema`, {
                headers: {
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }

            const data = await response.json();
            sqlSchema.set(data.collections || []);
        } catch (err) {
            console.error("Failed to load schema:", err);
        } finally {
            sqlSchemaLoading.set(false);
        }
    }

    async function executeSQL(query, confirmed = false) {
        if (!query.trim()) return;

        // Check for destructive operations
        if (!confirmed && isDestructiveQuery(query)) {
            sqlNeedsConfirmation.set(true);
            sqlConfirmMessage.set(
                "This is a destructive operation. Are you sure you want to proceed?"
            );
            sqlPendingQuery.set(query);
            return;
        }

        sqlLoading.set(true);
        sqlError.set("");
        sqlSuccess.set("");
        sqlResults.set([]);
        sqlColumns.set([]);
        sqlTotalRows.set(0);
        sqlRowsAffected.set(0);
        sqlIsMulti.set(false);
        sqlMultiResults.set([]);
        sqlTotalStatements.set(0);
        sqlSuccessfulCount.set(0);
        sqlFailedCount.set(0);

        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;

            const response = await fetch(`${baseUrl}/api/sql/execute`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
                body: JSON.stringify({
                    sql: query,
                    confirm: confirmed,
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }

            // Handle multi-statement response
            if (data.isMulti) {
                sqlIsMulti.set(true);
                sqlTotalStatements.set(data.totalStatements || 0);
                sqlSuccessfulCount.set(data.successfulCount || 0);
                sqlFailedCount.set(data.failedCount || 0);
                sqlMultiResults.set(data.results || []);
            }

            // Handle response (works for both single and multi-statement)
            if (data.rows) {
                sqlResults.set(data.rows);
                sqlColumns.set(data.columns || Object.keys(data.rows[0] || {}));
                sqlTotalRows.set(data.totalRows || data.rows.length);
            }

            if (data.rowsAffected !== undefined) {
                sqlRowsAffected.set(data.rowsAffected);
            }

            if (data.message) {
                sqlSuccess.set(data.message);
            } else if (!data.rows || data.rows.length === 0) {
                sqlSuccess.set("Query executed successfully");
            }

            // Add to history
            addToHistory(query, "sql");

            // Refresh schema if it was a DDL operation
            const upper = query.toUpperCase().trim();
            if (
                upper.includes("CREATE") ||
                upper.includes("ALTER") ||
                upper.includes("DROP")
            ) {
                loadSchema();
            }
        } catch (err) {
            sqlError.set(err.message || "Failed to execute query");
        } finally {
            sqlLoading.set(false);
        }
    }

    async function executeAI() {
        if (!$sqlAIQuery.trim()) return;

        sqlLoading.set(true);
        sqlError.set("");
        sqlSuccess.set("");
        // Clear previous results
        sqlResults.set([]);
        sqlColumns.set([]);
        sqlTotalRows.set(0);
        sqlRowsAffected.set(0);
        sqlIsMulti.set(false);
        sqlMultiResults.set([]);

        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;

            const response = await fetch(`${baseUrl}/api/sql/ai`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
                body: JSON.stringify({
                    query: $sqlAIQuery,
                    execute: true,
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }

            // Set the generated SQL
            if (data.sql) {
                sqlQuery.set(data.sql);
            }

            // Handle results if executed (AI endpoint nests results in 'result')
            const result = data.result || data;
            if (result.rows && result.rows.length > 0) {
                sqlResults.set(result.rows);
                sqlColumns.set(result.columns || Object.keys(result.rows[0] || {}));
                sqlTotalRows.set(result.totalRows || result.rows.length);
                sqlSuccess.set(`Query returned ${result.rows.length} row(s)`);
            } else if (result.rowsAffected !== undefined && result.rowsAffected > 0) {
                sqlRowsAffected.set(result.rowsAffected);
                sqlSuccess.set(result.message || `${result.rowsAffected} row(s) affected`);
            } else {
                sqlSuccess.set(result.message || "Query executed - no results");
            }

            // Add to history
            addToHistory($sqlAIQuery, "ai");
            if (data.sql) {
                addToHistory(data.sql, "sql");
            }
        } catch (err) {
            sqlError.set(err.message || "Failed to execute AI query");
        } finally {
            sqlLoading.set(false);
        }
    }

    function handleExecute(event) {
        const { query } = event.detail;
        executeSQL(query);
    }

    function handleSchemaInsert(event) {
        const { text } = event.detail;
        editorRef?.insertText(text);
    }

    function handleHistorySelect(event) {
        const { query, type } = event.detail;
        if (type === "ai") {
            sqlMode.set("ai");
            sqlAIQuery.set(query);
        } else {
            sqlMode.set("sql");
            sqlQuery.set(query);
        }
    }

    function confirmExecution() {
        const query = $sqlPendingQuery;
        sqlNeedsConfirmation.set(false);
        sqlConfirmMessage.set("");
        sqlPendingQuery.set("");
        executeSQL(query, true);
    }

    function cancelExecution() {
        sqlNeedsConfirmation.set(false);
        sqlConfirmMessage.set("");
        sqlPendingQuery.set("");
    }

    function switchMode(mode) {
        sqlMode.set(mode);
    }
</script>

<PageWrapper>
    <div class="sql-terminal">
        <aside class="terminal-sidebar">
            <SchemaExplorer
                schema={$sqlSchema}
                loading={$sqlSchemaLoading}
                on:insert={handleSchemaInsert}
                on:refresh={loadSchema}
            />
        </aside>

        <main class="terminal-main">
            <div class="terminal-header">
                <h1>SQL Terminal</h1>
                <div class="header-actions">
                    <div class="mode-toggle">
                        <button
                            type="button"
                            class="mode-btn"
                            class:active={$sqlMode === "sql"}
                            on:click={() => switchMode("sql")}
                        >
                            <i class="ri-code-line"></i>
                            SQL
                        </button>
                        <button
                            type="button"
                            class="mode-btn"
                            class:active={$sqlMode === "ai"}
                            on:click={() => switchMode("ai")}
                        >
                            <i class="ri-robot-line"></i>
                            AI
                        </button>
                    </div>
                    <QueryHistory on:select={handleHistorySelect} />
                </div>
            </div>

            <div class="terminal-body">
                {#if $sqlMode === "ai"}
                    <div class="ai-input-section">
                        <div class="ai-input-wrapper">
                            <textarea
                                bind:value={$sqlAIQuery}
                                placeholder="Describe what you want to query in natural language..."
                                class="ai-textarea"
                                rows="3"
                            ></textarea>
                            <button
                                type="button"
                                class="btn btn-primary ai-execute-btn"
                                on:click={executeAI}
                                disabled={$sqlLoading || !$sqlAIQuery.trim()}
                            >
                                {#if $sqlLoading}
                                    <i class="ri-loader-4-line animate-spin"></i>
                                {:else}
                                    <i class="ri-magic-line"></i>
                                {/if}
                                Generate & Execute
                            </button>
                        </div>
                        {#if $sqlQuery && $sqlMode === "ai"}
                            <div class="generated-sql">
                                <span class="label">Generated SQL:</span>
                                <pre><code>{$sqlQuery}</code></pre>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <SQLEditor
                        bind:this={editorRef}
                        bind:value={$sqlQuery}
                        loading={$sqlLoading}
                        on:execute={handleExecute}
                    />
                {/if}

                <div class="results-section">
                    {#if $sqlIsMulti && $sqlMultiResults.length > 1}
                        <div class="multi-results-header">
                            <span class="multi-badge">
                                <i class="ri-stack-line"></i>
                                {$sqlTotalStatements} statements
                            </span>
                            <span class="multi-summary">
                                <span class="success-count">
                                    <i class="ri-check-line"></i> {$sqlSuccessfulCount} succeeded
                                </span>
                                {#if $sqlFailedCount > 0}
                                    <span class="fail-count">
                                        <i class="ri-close-line"></i> {$sqlFailedCount} failed
                                    </span>
                                {/if}
                            </span>
                        </div>
                        <div class="multi-results-all">
                            {#each $sqlMultiResults as result, idx}
                                <div class="statement-result-block" class:failed={!result.success}>
                                    <div class="statement-header">
                                        <span class="statement-num">#{idx + 1}</span>
                                        <span class="statement-type">{result.type || 'UNKNOWN'}</span>
                                        {#if result.success}
                                            <i class="ri-check-line success-icon"></i>
                                        {:else}
                                            <i class="ri-close-line error-icon"></i>
                                        {/if}
                                        <span class="statement-message">{result.message}</span>
                                    </div>
                                    {#if result.rows && result.rows.length > 0}
                                        <div class="statement-table">
                                            <table class="mini-data-table">
                                                <thead>
                                                    <tr>
                                                        {#each (result.columns || Object.keys(result.rows[0] || {})) as col}
                                                            <th>{col}</th>
                                                        {/each}
                                                    </tr>
                                                </thead>
                                                <tbody>
                                                    {#each result.rows as row}
                                                        <tr>
                                                            {#each (result.columns || Object.keys(row)) as col}
                                                                <td>{row[col] === null ? 'NULL' : row[col]}</td>
                                                            {/each}
                                                        </tr>
                                                    {/each}
                                                </tbody>
                                            </table>
                                        </div>
                                    {:else if result.lastInsertId}
                                        <div class="statement-info">
                                            <span class="info-label">Last Insert ID:</span>
                                            <code>{result.lastInsertId}</code>
                                        </div>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {:else}
                        <ResultsTable
                            columns={$sqlColumns}
                            rows={$sqlResults}
                            totalRows={$sqlTotalRows}
                            rowsAffected={$sqlRowsAffected}
                            loading={$sqlLoading}
                            error={$sqlError}
                            success={$sqlSuccess}
                        />
                    {/if}
                </div>
            </div>
        </main>
    </div>

    <!-- Confirmation Modal -->
    {#if $sqlNeedsConfirmation}
        <div class="modal-overlay" on:click={cancelExecution}>
            <div class="modal-dialog" on:click|stopPropagation>
                <div class="modal-header">
                    <i class="ri-alert-line warning-icon"></i>
                    <h3>Confirm Destructive Operation</h3>
                </div>
                <div class="modal-body">
                    <p>{$sqlConfirmMessage}</p>
                    <pre class="pending-query"><code>{$sqlPendingQuery}</code></pre>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-transparent" on:click={cancelExecution}>
                        Cancel
                    </button>
                    <button type="button" class="btn btn-danger" on:click={confirmExecution}>
                        <i class="ri-check-line"></i>
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    {/if}
</PageWrapper>

<style>
    .sql-terminal {
        display: flex;
        height: 100%;
        min-height: 0;
    }
    .terminal-sidebar {
        width: 280px;
        min-width: 200px;
        max-width: 400px;
        padding: 15px;
        border-right: 1px solid var(--baseAlt2Color);
        overflow: hidden;
        display: flex;
        flex-direction: column;
    }
    .terminal-main {
        flex: 1;
        display: flex;
        flex-direction: column;
        min-width: 0;
        padding: 20px;
        overflow: hidden;
    }
    .terminal-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 20px;
    }
    .terminal-header h1 {
        margin: 0;
        font-size: 1.4em;
    }
    .header-actions {
        display: flex;
        align-items: center;
        gap: 12px;
    }
    .mode-toggle {
        display: flex;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        padding: 3px;
    }
    .mode-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 6px 12px;
        background: none;
        border: none;
        border-radius: calc(var(--baseRadius) - 2px);
        color: var(--txtSecondaryColor);
        font-size: 13px;
        cursor: pointer;
        transition: all 0.15s;
    }
    .mode-btn:hover {
        color: var(--txtPrimaryColor);
    }
    .mode-btn.active {
        background: var(--baseColor);
        color: var(--primaryColor);
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }
    .terminal-body {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 20px;
        min-height: 0;
        overflow: visible;
    }
    .ai-input-section {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }
    .ai-input-wrapper {
        display: flex;
        gap: 12px;
    }
    .ai-textarea {
        flex: 1;
        padding: 12px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        font-size: 14px;
        color: var(--txtPrimaryColor);
        resize: vertical;
    }
    .ai-textarea:focus {
        outline: none;
        border-color: var(--primaryColor);
    }
    .ai-execute-btn {
        align-self: flex-end;
        white-space: nowrap;
    }
    .generated-sql {
        padding: 12px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        border: 1px solid var(--baseAlt2Color);
        overflow: hidden;
    }
    .generated-sql .label {
        display: block;
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
        margin-bottom: 8px;
    }
    .generated-sql pre {
        margin: 0;
        padding: 0;
        font-family: var(--monospaceFontFamily);
        font-size: 13px;
        color: var(--txtPrimaryColor);
        white-space: pre-wrap;
        word-break: break-word;
        overflow-wrap: break-word;
        overflow-x: auto;
        max-width: 100%;
    }
    .generated-sql code {
        display: block;
        white-space: pre-wrap;
        word-break: break-word;
    }
    .results-section {
        flex: 1;
        min-height: 250px;
        max-height: 500px;
        overflow: auto;
        margin: 0;
        padding: 0;
    }
    .multi-results-header {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 10px 12px;
        background: var(--baseAlt1Color);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius) var(--baseRadius) 0 0;
        margin-bottom: -1px;
    }
    .multi-badge {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        font-weight: 600;
        color: var(--primaryColor);
        background: var(--primaryAltColor);
        padding: 4px 10px;
        border-radius: 12px;
    }
    .multi-summary {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 12px;
    }
    .success-count {
        display: flex;
        align-items: center;
        gap: 4px;
        color: var(--successColor);
    }
    .fail-count {
        display: flex;
        align-items: center;
        gap: 4px;
        color: var(--dangerColor);
    }
    .multi-results-all {
        display: flex;
        flex-direction: column;
        gap: 12px;
        max-height: 100%;
        overflow-y: auto;
    }
    .statement-result-block {
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        overflow: hidden;
        background: var(--baseColor);
    }
    .statement-result-block.failed {
        border-color: var(--dangerColor);
    }
    .statement-result-block .statement-header {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 12px;
        background: var(--baseAlt1Color);
        border-bottom: 1px solid var(--baseAlt2Color);
        font-size: 12px;
    }
    .statement-table {
        max-height: 200px;
        overflow: auto;
        position: relative;
    }
    .mini-data-table {
        width: 100%;
        border-collapse: separate;
        border-spacing: 0;
        font-size: 12px;
    }
    .mini-data-table th,
    .mini-data-table td {
        padding: 6px 10px;
        text-align: left;
        border-bottom: 1px solid var(--baseAlt1Color);
        white-space: nowrap;
    }
    .mini-data-table th {
        background: var(--baseAlt1Color);
        font-weight: 600;
        color: var(--txtSecondaryColor);
        position: sticky;
        top: 0;
        z-index: 1;
    }
    .mini-data-table thead {
        position: sticky;
        top: 0;
        z-index: 2;
    }
    .mini-data-table td {
        font-family: var(--monospaceFontFamily);
        color: var(--txtPrimaryColor);
        background: var(--baseColor);
    }
    .mini-data-table tbody tr:hover td {
        background: var(--baseAlt1Color);
    }
    .statement-info {
        padding: 8px 12px;
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 8px;
    }
    .statement-info .info-label {
        color: var(--txtSecondaryColor);
    }
    .statement-info code {
        background: var(--baseAlt1Color);
        padding: 2px 6px;
        border-radius: 4px;
        font-family: var(--monospaceFontFamily);
    }
    .statement-num {
        font-weight: 600;
        color: var(--txtHintColor);
        min-width: 24px;
    }
    .statement-type {
        font-weight: 600;
        color: var(--primaryColor);
        min-width: 80px;
    }
    .statement-header .success-icon {
        color: var(--successColor);
        font-size: 14px;
    }
    .statement-header .error-icon {
        color: var(--dangerColor);
        font-size: 14px;
    }
    .statement-message {
        color: var(--txtSecondaryColor);
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    /* Modal styles */
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }
    .modal-dialog {
        background: var(--baseColor);
        border-radius: var(--baseRadius);
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
        max-width: 500px;
        width: 90%;
    }
    .modal-header {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 20px;
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .modal-header h3 {
        margin: 0;
        font-size: 1.1em;
    }
    .warning-icon {
        font-size: 24px;
        color: var(--warningColor);
    }
    .modal-body {
        padding: 20px;
    }
    .modal-body p {
        margin: 0 0 15px;
        color: var(--txtSecondaryColor);
    }
    .pending-query {
        margin: 0;
        padding: 12px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        font-family: var(--monospaceFontFamily);
        font-size: 12px;
        overflow-x: auto;
    }
    .modal-footer {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        padding: 15px 20px;
        border-top: 1px solid var(--baseAlt2Color);
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>

