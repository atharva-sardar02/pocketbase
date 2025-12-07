<script>
    import { createEventDispatcher } from "svelte";
    import {
        uploadedFile,
        fileFormat,
        fileDelimiter,
        importError
    } from "@/stores/import";

    const dispatch = createEventDispatcher();

    let isDragging = false;
    let fileInput;

    const acceptedTypes = [
        "text/csv",
        "application/json",
        "text/plain",
        ".csv",
        ".json"
    ];

    function handleDragOver(e) {
        e.preventDefault();
        isDragging = true;
    }

    function handleDragLeave(e) {
        e.preventDefault();
        isDragging = false;
    }

    function handleDrop(e) {
        e.preventDefault();
        isDragging = false;
        
        const files = e.dataTransfer?.files;
        if (files?.length > 0) {
            processFile(files[0]);
        }
    }

    function handleFileSelect(e) {
        const files = e.target?.files;
        if (files?.length > 0) {
            processFile(files[0]);
        }
    }

    function processFile(file) {
        importError.set("");
        
        // Validate file type
        const ext = file.name.split(".").pop()?.toLowerCase();
        if (!["csv", "json", "txt"].includes(ext)) {
            importError.set("Unsupported file type. Please upload a CSV or JSON file.");
            return;
        }

        // Validate file size (max 10MB)
        const maxSize = 10 * 1024 * 1024;
        if (file.size > maxSize) {
            importError.set("File is too large. Maximum size is 10MB.");
            return;
        }

        // Detect format
        if (ext === "json") {
            fileFormat.set("json");
        } else {
            fileFormat.set("csv");
        }

        // Read file content
        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target?.result;
            if (content) {
                uploadedFile.set({
                    name: file.name,
                    size: file.size,
                    type: file.type,
                    content: content
                });
                dispatch("fileLoaded", { file, content });
            }
        };
        reader.onerror = () => {
            importError.set("Failed to read file.");
        };
        reader.readAsText(file);
    }

    function clearFile() {
        uploadedFile.set(null);
        if (fileInput) {
            fileInput.value = "";
        }
        dispatch("fileCleared");
    }

    function formatFileSize(bytes) {
        if (bytes < 1024) return bytes + " B";
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
        return (bytes / (1024 * 1024)).toFixed(2) + " MB";
    }

    function openFileDialog() {
        fileInput?.click();
    }
</script>

<div class="file-upload">
    {#if !$uploadedFile}
        <div 
            class="drop-zone" 
            class:dragging={isDragging}
            on:dragover={handleDragOver}
            on:dragleave={handleDragLeave}
            on:drop={handleDrop}
            on:click={openFileDialog}
            role="button"
            tabindex="0"
            on:keypress={(e) => e.key === "Enter" && openFileDialog()}
        >
            <i class="ri-upload-cloud-2-line"></i>
            <p class="drop-text">
                <strong>Drop your file here</strong>
                <span>or click to browse</span>
            </p>
            <p class="file-types">Supports CSV and JSON files (max 10MB)</p>
        </div>
    {:else}
        <div class="file-info">
            <div class="file-icon">
                {#if $fileFormat === "json"}
                    <i class="ri-braces-line"></i>
                {:else}
                    <i class="ri-file-text-line"></i>
                {/if}
            </div>
            <div class="file-details">
                <p class="file-name">{$uploadedFile.name}</p>
                <p class="file-meta">
                    {formatFileSize($uploadedFile.size)} • {$fileFormat.toUpperCase()}
                </p>
            </div>
            <button 
                type="button" 
                class="btn btn-sm btn-secondary"
                on:click={clearFile}
                title="Remove file"
            >
                <i class="ri-close-line"></i>
            </button>
        </div>

        {#if $fileFormat === "csv"}
            <div class="delimiter-select">
                <label for="delimiter">Delimiter:</label>
                <select id="delimiter" bind:value={$fileDelimiter}>
                    <option value=",">Comma (,)</option>
                    <option value=";">Semicolon (;)</option>
                    <option value="\t">Tab</option>
                    <option value="|">Pipe (|)</option>
                </select>
            </div>
        {/if}
    {/if}

    <input
        type="file"
        bind:this={fileInput}
        on:change={handleFileSelect}
        accept=".csv,.json,text/csv,application/json"
        style="display: none;"
    />
</div>

<style>
    .file-upload {
        width: 100%;
    }

    .drop-zone {
        border: 2px dashed var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        padding: 40px 20px;
        text-align: center;
        cursor: pointer;
        transition: all 0.2s ease;
        background: var(--baseColor);
    }

    .drop-zone:hover,
    .drop-zone.dragging {
        border-color: var(--primaryColor);
        background: rgba(var(--primaryColorRaw), 0.05);
    }

    .drop-zone i {
        font-size: 48px;
        color: var(--txtDisabledColor);
        margin-bottom: 12px;
    }

    .drop-zone:hover i,
    .drop-zone.dragging i {
        color: var(--primaryColor);
    }

    .drop-text {
        margin: 0 0 8px;
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .drop-text strong {
        font-size: 16px;
        color: var(--txtPrimaryColor);
    }

    .drop-text span {
        font-size: 14px;
        color: var(--txtHintColor);
    }

    .file-types {
        margin: 0;
        font-size: 12px;
        color: var(--txtDisabledColor);
    }

    .file-info {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 16px;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .file-icon {
        width: 48px;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(var(--primaryColorRaw), 0.1);
        border-radius: var(--baseRadius);
    }

    .file-icon i {
        font-size: 24px;
        color: var(--primaryColor);
    }

    .file-details {
        flex: 1;
        min-width: 0;
    }

    .file-name {
        margin: 0;
        font-weight: 600;
        color: var(--txtPrimaryColor);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .file-meta {
        margin: 4px 0 0;
        font-size: 12px;
        color: var(--txtHintColor);
    }

    .delimiter-select {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-top: 12px;
        padding: 12px 16px;
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
    }

    .delimiter-select label {
        font-size: 13px;
        color: var(--txtHintColor);
    }

    .delimiter-select select {
        padding: 6px 10px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        font-size: 13px;
    }
</style>
