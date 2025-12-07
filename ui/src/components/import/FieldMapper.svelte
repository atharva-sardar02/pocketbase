<script>
    import {
        previewHeaders,
        fieldMapping,
        collectionFields,
        validateLoading,
        autoDetectMapping,
        hasMappedFields,
        mappedFieldCount
    } from "@/stores/import";

    $: fieldNames = Object.keys($collectionFields);
    $: totalHeaders = $previewHeaders.length;

    function handleMappingChange(header, value) {
        fieldMapping.update(m => ({
            ...m,
            [header]: value
        }));
    }

    function handleAutoDetect() {
        autoDetectMapping($previewHeaders, $collectionFields);
    }

    function clearAllMappings() {
        const cleared = {};
        $previewHeaders.forEach(h => {
            cleared[h] = "";
        });
        fieldMapping.set(cleared);
    }

    function getFieldType(fieldName) {
        return $collectionFields[fieldName] || "";
    }

    function getFieldTypeIcon(type) {
        switch (type) {
            case "text": return "ri-text";
            case "number": return "ri-hashtag";
            case "bool": return "ri-checkbox-line";
            case "email": return "ri-mail-line";
            case "url": return "ri-link";
            case "date": return "ri-calendar-line";
            case "select": return "ri-list-check";
            case "json": return "ri-braces-line";
            case "file": return "ri-attachment-line";
            case "relation": return "ri-links-line";
            case "editor": return "ri-file-text-line";
            default: return "ri-input-field";
        }
    }
</script>

<div class="field-mapper">
    {#if $validateLoading}
        <div class="loading-state">
            <i class="ri-loader-4-line spinning"></i>
            <span>Loading collection fields...</span>
        </div>
    {:else if fieldNames.length === 0}
        <div class="empty-state">
            <i class="ri-error-warning-line"></i>
            <span>No fields found in collection</span>
        </div>
    {:else}
        <div class="mapper-header">
            <div class="mapper-info">
                <span class="mapped-count">
                    <i class="ri-link"></i>
                    {$mappedFieldCount} of {totalHeaders} mapped
                </span>
            </div>
            <div class="mapper-actions">
                <button 
                    type="button" 
                    class="btn btn-sm btn-secondary"
                    on:click={handleAutoDetect}
                    title="Auto-detect mappings based on column names"
                >
                    <i class="ri-magic-line"></i>
                    Auto-detect
                </button>
                <button 
                    type="button" 
                    class="btn btn-sm btn-secondary"
                    on:click={clearAllMappings}
                    title="Clear all mappings"
                >
                    <i class="ri-delete-bin-line"></i>
                    Clear all
                </button>
            </div>
        </div>

        <div class="mapping-list">
            <div class="mapping-header-row">
                <div class="source-col">Source Column</div>
                <div class="arrow-col"></div>
                <div class="target-col">Target Field</div>
            </div>

            {#each $previewHeaders as header}
                {@const currentMapping = $fieldMapping[header] || ""}
                {@const mappedType = currentMapping ? getFieldType(currentMapping) : ""}
                <div class="mapping-row" class:mapped={currentMapping && currentMapping !== "-"}>
                    <div class="source-col">
                        <span class="col-name">{header}</span>
                    </div>
                    <div class="arrow-col">
                        <i class="ri-arrow-right-line"></i>
                    </div>
                    <div class="target-col">
                        <select 
                            value={currentMapping}
                            on:change={(e) => handleMappingChange(header, e.target.value)}
                            class:mapped={currentMapping && currentMapping !== "-"}
                        >
                            <option value="">-- Skip this column --</option>
                            <option value="-">-- Skip this column --</option>
                            {#each fieldNames as fieldName}
                                {@const type = getFieldType(fieldName)}
                                <option value={fieldName}>
                                    {fieldName} ({type})
                                </option>
                            {/each}
                        </select>
                        {#if mappedType}
                            <span class="field-type-badge" title={mappedType}>
                                <i class={getFieldTypeIcon(mappedType)}></i>
                            </span>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>

        <div class="available-fields">
            <h4>Available Fields</h4>
            <div class="fields-grid">
                {#each fieldNames as fieldName}
                    {@const type = getFieldType(fieldName)}
                    {@const isMapped = Object.values($fieldMapping).includes(fieldName)}
                    <div class="field-tag" class:mapped={isMapped}>
                        <i class={getFieldTypeIcon(type)}></i>
                        <span>{fieldName}</span>
                        <span class="type-label">{type}</span>
                        {#if isMapped}
                            <i class="ri-check-line mapped-icon"></i>
                        {/if}
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

<style>
    .field-mapper {
        width: 100%;
    }

    .loading-state,
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 20px;
        text-align: center;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .loading-state i,
    .empty-state i {
        font-size: 36px;
        color: var(--txtDisabledColor);
        margin-bottom: 12px;
    }

    .spinning {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    .mapper-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
        padding: 12px 16px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
    }

    .mapped-count {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: var(--txtHintColor);
    }

    .mapper-actions {
        display: flex;
        gap: 8px;
    }

    .mapping-list {
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        overflow: hidden;
    }

    .mapping-header-row {
        display: grid;
        grid-template-columns: 1fr 40px 1fr;
        gap: 12px;
        padding: 12px 16px;
        background: var(--baseAlt1Color);
        font-size: 12px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
    }

    .mapping-row {
        display: grid;
        grid-template-columns: 1fr 40px 1fr;
        gap: 12px;
        padding: 12px 16px;
        background: var(--baseColor);
        border-top: 1px solid var(--baseAlt2Color);
        align-items: center;
    }

    .mapping-row.mapped {
        background: rgba(var(--successColorRaw), 0.05);
    }

    .source-col {
        font-size: 13px;
    }

    .col-name {
        display: inline-block;
        padding: 4px 10px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        font-family: monospace;
        font-size: 12px;
    }

    .arrow-col {
        text-align: center;
        color: var(--txtDisabledColor);
    }

    .target-col {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .target-col select {
        flex: 1;
        padding: 8px 12px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        font-size: 13px;
    }

    .target-col select.mapped {
        border-color: var(--successColor);
        background: rgba(var(--successColorRaw), 0.05);
    }

    .field-type-badge {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        color: var(--primaryColor);
    }

    .available-fields {
        margin-top: 20px;
    }

    .available-fields h4 {
        margin: 0 0 12px;
        font-size: 13px;
        font-weight: 600;
        color: var(--txtHintColor);
    }

    .fields-grid {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
    }

    .field-tag {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 6px 10px;
        background: var(--baseAlt1Color);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        font-size: 12px;
    }

    .field-tag.mapped {
        background: rgba(var(--successColorRaw), 0.1);
        border-color: var(--successColor);
    }

    .field-tag i {
        color: var(--txtHintColor);
    }

    .field-tag span {
        color: var(--txtPrimaryColor);
    }

    .field-tag .type-label {
        color: var(--txtHintColor);
        font-size: 11px;
    }

    .field-tag .mapped-icon {
        color: var(--successColor);
    }
</style>
