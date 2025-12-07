<script>
    import { onMount, onDestroy } from "svelte";
    import { push } from "svelte-spa-router";
    import { pageTitle } from "@/stores/app";
    import {
        importStep,
        selectedCollection,
        uploadedFile,
        fileFormat,
        fileDelimiter,
        previewHeaders,
        previewRows,
        previewTotalRows,
        previewErrors,
        fieldMapping,
        collectionFields,
        importResults,
        previewLoading,
        validateLoading,
        executeLoading,
        importError,
        hasMappedFields,
        resetImportData,
        nextStep,
        prevStep,
        goToStep,
        initializeMapping
    } from "@/stores/import";
    import ApiClient from "@/utils/ApiClient";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import FileUpload from "@/components/import/FileUpload.svelte";
    import DataPreview from "@/components/import/DataPreview.svelte";
    import FieldMapper from "@/components/import/FieldMapper.svelte";
    import ImportProgress from "@/components/import/ImportProgress.svelte";

    $pageTitle = "Import Data";

    let collections = [];
    let collectionsLoading = true;

    const steps = [
        { num: 1, label: "Select & Upload" },
        { num: 2, label: "Preview" },
        { num: 3, label: "Map Fields" },
        { num: 4, label: "Import" }
    ];

    onMount(async () => {
        resetImportData();
        await loadCollections();
    });

    onDestroy(() => {
        // Cleanup if needed
    });

    async function loadCollections() {
        collectionsLoading = true;
        try {
            const result = await ApiClient.collections.getFullList({
                $cancelKey: "importCollections",
                filter: 'type != "view"'
            });
            collections = result.filter(c => !c.system);
        } catch (err) {
            if (!err?.isAbort) {
                console.error("Failed to load collections:", err);
            }
        } finally {
            collectionsLoading = false;
        }
    }

    async function handleFileLoaded(event) {
        const { content } = event.detail;
        await previewFile(content);
    }

    async function previewFile(content) {
        previewLoading.set(true);
        previewErrors.set([]);
        importError.set("");

        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;

            const response = await fetch(`${baseUrl}/api/import/preview`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
                body: JSON.stringify({
                    data: content,
                    format: $fileFormat,
                    delimiter: $fileDelimiter
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }

            if (data.errors && data.errors.length > 0) {
                previewErrors.set(data.errors);
            } else {
                previewHeaders.set(data.headers || []);
                previewRows.set(data.sampleRows || []);
                previewTotalRows.set(data.totalRows || 0);
                initializeMapping(data.headers || []);
            }
        } catch (err) {
            console.error("Preview error:", err);
            previewErrors.set([err.message || "Failed to parse file"]);
        } finally {
            previewLoading.set(false);
        }
    }

    async function handleDelimiterChange() {
        if ($uploadedFile?.content) {
            await previewFile($uploadedFile.content);
        }
    }

    async function loadCollectionFields() {
        if (!$selectedCollection) return;

        validateLoading.set(true);
        importError.set("");

        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;

            const response = await fetch(`${baseUrl}/api/import/validate`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
                body: JSON.stringify({
                    collection: $selectedCollection,
                    mapping: {}
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }

            collectionFields.set(data.fieldTypes || {});
        } catch (err) {
            console.error("Validation error:", err);
            importError.set(err.message || "Failed to load collection fields");
        } finally {
            validateLoading.set(false);
        }
    }

    async function executeImport() {
        if (!$selectedCollection || !$uploadedFile?.content) return;

        executeLoading.set(true);
        importError.set("");
        importResults.set({
            totalRows: $previewTotalRows,
            successCount: 0,
            failureCount: 0,
            errors: []
        });

        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;

            const response = await fetch(`${baseUrl}/api/import/execute`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token
                        ? `Bearer ${ApiClient.authStore.token}`
                        : "",
                },
                body: JSON.stringify({
                    collection: $selectedCollection,
                    data: $uploadedFile.content,
                    format: $fileFormat,
                    delimiter: $fileDelimiter,
                    mapping: $fieldMapping
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || `HTTP ${response.status}`);
            }

            importResults.set({
                totalRows: data.totalRows || 0,
                successCount: data.successCount || 0,
                failureCount: data.failureCount || 0,
                errors: data.errors || []
            });
        } catch (err) {
            console.error("Import error:", err);
            importError.set(err.message || "Failed to import data");
        } finally {
            executeLoading.set(false);
        }
    }

    function handleNext() {
        if ($importStep === 2) {
            loadCollectionFields();
        }
        nextStep();
    }

    function handleBack() {
        prevStep();
    }

    function handleStartImport() {
        executeImport();
    }

    function handleStartOver() {
        resetImportData();
    }

    function handleViewCollection() {
        push(`/collections?collection=${$selectedCollection}`);
    }

    // Validation for step navigation
    $: canProceedStep1 = $selectedCollection && $uploadedFile && $previewHeaders.length > 0;
    $: canProceedStep2 = $previewHeaders.length > 0 && $previewTotalRows > 0;
    $: canProceedStep3 = $hasMappedFields;
    $: isImportComplete = !$executeLoading && $importResults.successCount + $importResults.failureCount === $importResults.totalRows && $importResults.totalRows > 0;

    // Watch for delimiter changes
    $: if ($fileDelimiter && $uploadedFile?.content) {
        // Debounce this in real app
    }
</script>

<PageWrapper>
    <div class="import-wizard">
        <header class="wizard-header">
            <h1>
                <i class="ri-upload-cloud-2-line"></i>
                Import Data
            </h1>
            <p class="subtitle">Import CSV or JSON data into your collections</p>
        </header>

        <!-- Step Indicator -->
        <div class="step-indicator">
            {#each steps as step}
                <div 
                    class="step" 
                    class:active={$importStep === step.num}
                    class:completed={$importStep > step.num}
                    class:clickable={$importStep > step.num}
                    on:click={() => $importStep > step.num && goToStep(step.num)}
                    on:keypress={(e) => e.key === "Enter" && $importStep > step.num && goToStep(step.num)}
                    role="button"
                    tabindex={$importStep > step.num ? 0 : -1}
                >
                    <div class="step-number">
                        {#if $importStep > step.num}
                            <i class="ri-check-line"></i>
                        {:else}
                            {step.num}
                        {/if}
                    </div>
                    <span class="step-label">{step.label}</span>
                </div>
                {#if step.num < steps.length}
                    <div class="step-connector" class:completed={$importStep > step.num}></div>
                {/if}
            {/each}
        </div>

        <!-- Step Content -->
        <div class="wizard-content">
            {#if $importStep === 1}
                <!-- Step 1: Select Collection & Upload File -->
                <div class="step-content">
                    <div class="form-section">
                        <label for="collection-select" class="section-label">
                            <i class="ri-database-2-line"></i>
                            Select Collection
                        </label>
                        {#if collectionsLoading}
                            <div class="loading-inline">
                                <i class="ri-loader-4-line spinning"></i>
                                Loading collections...
                            </div>
                        {:else}
                            <select 
                                id="collection-select" 
                                bind:value={$selectedCollection}
                                class="collection-select"
                            >
                                <option value="">-- Select a collection --</option>
                                {#each collections as col}
                                    <option value={col.name}>{col.name}</option>
                                {/each}
                            </select>
                        {/if}
                    </div>

                    <div class="form-section">
                        <label class="section-label">
                            <i class="ri-file-upload-line"></i>
                            Upload File
                        </label>
                        <FileUpload 
                            on:fileLoaded={handleFileLoaded}
                            on:fileCleared={() => {
                                previewHeaders.set([]);
                                previewRows.set([]);
                                previewTotalRows.set(0);
                            }}
                        />
                    </div>

                    {#if $previewHeaders.length > 0}
                        <div class="form-section preview-section">
                            <label class="section-label">
                                <i class="ri-table-line"></i>
                                Data Preview
                            </label>
                            <DataPreview />
                        </div>
                    {/if}
                </div>
            {:else if $importStep === 2}
                <!-- Step 2: Preview Data -->
                <div class="step-content">
                    <div class="preview-full">
                        <DataPreview />
                    </div>
                </div>
            {:else if $importStep === 3}
                <!-- Step 3: Map Fields -->
                <div class="step-content">
                    <div class="mapping-section">
                        <p class="mapping-instructions">
                            Map your CSV/JSON columns to the fields in <strong>{$selectedCollection}</strong> collection.
                        </p>
                        <FieldMapper />
                    </div>
                </div>
            {:else if $importStep === 4}
                <!-- Step 4: Import -->
                <div class="step-content">
                    <ImportProgress 
                        on:startOver={handleStartOver}
                        on:viewCollection={handleViewCollection}
                    />
                </div>
            {/if}
        </div>

        <!-- Error Display -->
        {#if $importError}
            <div class="alert alert-danger">
                <i class="ri-error-warning-line"></i>
                {$importError}
            </div>
        {/if}

        <!-- Navigation Buttons -->
        <div class="wizard-footer">
            {#if $importStep > 1 && !isImportComplete}
                <button 
                    type="button" 
                    class="btn btn-secondary"
                    on:click={handleBack}
                    disabled={$executeLoading}
                >
                    <i class="ri-arrow-left-line"></i>
                    Back
                </button>
            {:else}
                <div></div>
            {/if}

            <div class="footer-right">
                {#if $importStep === 1}
                    <button 
                        type="button" 
                        class="btn btn-primary"
                        on:click={handleNext}
                        disabled={!canProceedStep1}
                    >
                        Continue
                        <i class="ri-arrow-right-line"></i>
                    </button>
                {:else if $importStep === 2}
                    <button 
                        type="button" 
                        class="btn btn-primary"
                        on:click={handleNext}
                        disabled={!canProceedStep2}
                    >
                        Continue to Mapping
                        <i class="ri-arrow-right-line"></i>
                    </button>
                {:else if $importStep === 3}
                    <button 
                        type="button" 
                        class="btn btn-primary btn-expanded"
                        on:click={() => { nextStep(); handleStartImport(); }}
                        disabled={!canProceedStep3}
                    >
                        <i class="ri-upload-2-line"></i>
                        Start Import ({$previewTotalRows} rows)
                    </button>
                {:else if $importStep === 4 && !isImportComplete}
                    <button 
                        type="button" 
                        class="btn btn-primary"
                        disabled={$executeLoading}
                    >
                        <i class="ri-loader-4-line spinning"></i>
                        Importing...
                    </button>
                {/if}
            </div>
        </div>
    </div>
</PageWrapper>

<style>
    .import-wizard {
        padding: 20px;
        max-width: 900px;
        margin: 0 auto;
    }

    .wizard-header {
        text-align: center;
        margin-bottom: 32px;
    }

    .wizard-header h1 {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        margin: 0 0 8px;
        font-size: 1.5em;
    }

    .wizard-header h1 i {
        color: var(--primaryColor);
    }

    .subtitle {
        margin: 0;
        color: var(--txtHintColor);
    }

    /* Step Indicator */
    .step-indicator {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 32px;
    }

    .step {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
    }

    .step.clickable {
        cursor: pointer;
    }

    .step-number {
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        background: var(--baseAlt2Color);
        color: var(--txtHintColor);
        font-weight: 600;
        font-size: 14px;
        transition: all 0.2s ease;
    }

    .step.active .step-number {
        background: var(--primaryColor);
        color: white;
    }

    .step.completed .step-number {
        background: var(--successColor);
        color: white;
    }

    .step-label {
        font-size: 12px;
        color: var(--txtHintColor);
        white-space: nowrap;
    }

    .step.active .step-label {
        color: var(--txtPrimaryColor);
        font-weight: 600;
    }

    .step-connector {
        width: 60px;
        height: 2px;
        background: var(--baseAlt2Color);
        margin: 0 8px;
        margin-bottom: 24px;
    }

    .step-connector.completed {
        background: var(--successColor);
    }

    /* Wizard Content */
    .wizard-content {
        min-height: 400px;
        margin-bottom: 24px;
    }

    .step-content {
        animation: fadeIn 0.2s ease;
    }

    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }

    .form-section {
        margin-bottom: 24px;
    }

    .section-label {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 12px;
        font-size: 14px;
        font-weight: 600;
        color: var(--txtPrimaryColor);
    }

    .section-label i {
        color: var(--primaryColor);
    }

    .collection-select {
        width: 100%;
        padding: 12px 16px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        font-size: 14px;
    }

    .loading-inline {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px;
        color: var(--txtHintColor);
    }

    .spinning {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    .preview-section {
        margin-top: 24px;
        padding-top: 24px;
        border-top: 1px solid var(--baseAlt2Color);
    }

    .mapping-instructions {
        margin: 0 0 16px;
        color: var(--txtHintColor);
    }

    .alert {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 12px 16px;
        border-radius: var(--baseRadius);
        margin-bottom: 20px;
    }

    .alert-danger {
        background: rgba(var(--dangerColorRaw), 0.1);
        color: var(--dangerColor);
    }

    /* Footer */
    .wizard-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 20px;
        border-top: 1px solid var(--baseAlt2Color);
    }

    .footer-right {
        display: flex;
        gap: 12px;
    }

    .btn-expanded {
        padding-left: 24px;
        padding-right: 24px;
    }

    /* Responsive */
    @media (max-width: 768px) {
        .step-indicator {
            flex-wrap: wrap;
            gap: 8px;
        }
        
        .step-connector {
            display: none;
        }

        .step-label {
            font-size: 11px;
        }

        .wizard-footer {
            flex-direction: column;
            gap: 12px;
        }

        .footer-right {
            width: 100%;
            justify-content: flex-end;
        }
    }
</style>
