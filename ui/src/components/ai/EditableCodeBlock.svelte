<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let code = "";
    export let language = "sql"; // "sql" or "filter"
    export let editable = true;
    export let placeholder = "Enter code...";

    const dispatch = createEventDispatcher();

    let isEditing = false;
    let editedCode = code;
    let textareaEl;

    $: editedCode = code;

    function startEditing() {
        if (!editable) return;
        isEditing = true;
        editedCode = code;
        // Focus the textarea after it renders
        setTimeout(() => {
            if (textareaEl) {
                textareaEl.focus();
                textareaEl.select();
            }
        }, 0);
    }

    function cancelEditing() {
        isEditing = false;
        editedCode = code;
    }

    function applyEdit() {
        isEditing = false;
        dispatch("change", { code: editedCode });
    }

    function handleKeydown(e) {
        if (e.key === "Escape") {
            cancelEditing();
        } else if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
            applyEdit();
        }
    }

    function copyToClipboard() {
        navigator.clipboard.writeText(code);
        dispatch("copy");
    }

    function executeCode() {
        dispatch("execute", { code: isEditing ? editedCode : code });
    }
</script>

<div class="editable-code-block" class:editing={isEditing}>
    <div class="code-header">
        <span class="language-badge">{language === "sql" ? "SQL" : "Filter"}</span>
        <div class="code-actions">
            {#if !isEditing}
                <button 
                    class="btn btn-xs btn-transparent"
                    on:click={copyToClipboard}
                    use:tooltip={"Copy to clipboard"}
                >
                    <i class="ri-file-copy-line"></i>
                </button>
                {#if editable}
                    <button 
                        class="btn btn-xs btn-transparent"
                        on:click={startEditing}
                        use:tooltip={"Edit"}
                    >
                        <i class="ri-edit-line"></i>
                    </button>
                {/if}
                <button 
                    class="btn btn-xs btn-primary"
                    on:click={executeCode}
                    use:tooltip={"Execute (Ctrl+Enter)"}
                >
                    <i class="ri-play-line"></i>
                    Run
                </button>
            {:else}
                <button 
                    class="btn btn-xs btn-transparent btn-hint"
                    on:click={cancelEditing}
                >
                    Cancel
                </button>
                <button 
                    class="btn btn-xs btn-success"
                    on:click={applyEdit}
                    use:tooltip={"Apply and run (Ctrl+Enter)"}
                >
                    <i class="ri-check-line"></i>
                    Apply
                </button>
            {/if}
        </div>
    </div>
    <div class="code-content">
        {#if isEditing}
            <textarea
                bind:this={textareaEl}
                bind:value={editedCode}
                on:keydown={handleKeydown}
                class="code-textarea"
                {placeholder}
                spellcheck="false"
            ></textarea>
        {:else}
            <pre class="code-display" on:dblclick={startEditing}><code>{code || placeholder}</code></pre>
        {/if}
    </div>
    {#if isEditing}
        <div class="code-hint">
            Press <kbd>Ctrl</kbd>+<kbd>Enter</kbd> to apply, <kbd>Esc</kbd> to cancel
        </div>
    {/if}
</div>

<style>
    .editable-code-block {
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseAlt1Color);
        overflow: hidden;
    }
    .editable-code-block.editing {
        border-color: var(--primaryColor);
    }
    .code-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 6px 10px;
        background: var(--baseAlt2Color);
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .language-badge {
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
        background: var(--baseColor);
        padding: 2px 8px;
        border-radius: 3px;
    }
    .code-actions {
        display: flex;
        align-items: center;
        gap: 4px;
    }
    .code-content {
        padding: 0;
    }
    .code-display {
        margin: 0;
        padding: 12px;
        font-family: var(--monospaceFontFamily);
        font-size: 13px;
        line-height: 1.5;
        white-space: pre-wrap;
        word-break: break-word;
        cursor: text;
        min-height: 60px;
        color: var(--txtPrimaryColor);
    }
    .code-display:hover {
        background: var(--baseAlt2Color);
    }
    .code-display code {
        font-family: inherit;
    }
    .code-textarea {
        width: 100%;
        min-height: 100px;
        padding: 12px;
        border: none;
        background: var(--baseColor);
        font-family: var(--monospaceFontFamily);
        font-size: 13px;
        line-height: 1.5;
        color: var(--txtPrimaryColor);
        resize: vertical;
    }
    .code-textarea:focus {
        outline: none;
    }
    .code-hint {
        padding: 6px 10px;
        font-size: 11px;
        color: var(--txtHintColor);
        background: var(--baseAlt2Color);
        border-top: 1px solid var(--baseAlt2Color);
    }
    .code-hint kbd {
        display: inline-block;
        padding: 1px 4px;
        font-size: 10px;
        font-family: var(--monospaceFontFamily);
        background: var(--baseColor);
        border: 1px solid var(--baseAlt3Color);
        border-radius: 3px;
    }
</style>
