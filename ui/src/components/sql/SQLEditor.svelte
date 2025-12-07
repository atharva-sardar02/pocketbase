<script>
    import { createEventDispatcher, onMount } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let value = "";
    export let placeholder = "Enter SQL query...";
    export let disabled = false;
    export let loading = false;

    const dispatch = createEventDispatcher();

    let textareaEl;
    let lineNumbers = "1";

    $: updateLineNumbers(value);

    function updateLineNumbers(text) {
        const lines = (text || "").split("\n").length;
        lineNumbers = Array.from({ length: lines }, (_, i) => i + 1).join("\n");
    }

    function handleKeydown(e) {
        // Ctrl/Cmd + Enter to execute
        if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
            e.preventDefault();
            execute();
            return;
        }

        // Tab to insert spaces
        if (e.key === "Tab") {
            e.preventDefault();
            const start = textareaEl.selectionStart;
            const end = textareaEl.selectionEnd;
            const spaces = "  ";
            value = value.substring(0, start) + spaces + value.substring(end);
            // Set cursor position after spaces
            setTimeout(() => {
                textareaEl.selectionStart = textareaEl.selectionEnd = start + spaces.length;
            }, 0);
        }
    }

    function execute() {
        if (!value.trim() || disabled || loading) return;
        dispatch("execute", { query: value });
    }

    function clear() {
        value = "";
        textareaEl?.focus();
        dispatch("clear");
    }

    export function focus() {
        textareaEl?.focus();
    }

    export function insertText(text) {
        if (!textareaEl) return;
        const start = textareaEl.selectionStart;
        const end = textareaEl.selectionEnd;
        value = value.substring(0, start) + text + value.substring(end);
        setTimeout(() => {
            textareaEl.selectionStart = textareaEl.selectionEnd = start + text.length;
            textareaEl.focus();
        }, 0);
    }

    onMount(() => {
        textareaEl?.focus();
    });
</script>

<div class="sql-editor" class:disabled class:loading>
    <div class="editor-header">
        <span class="editor-title">SQL Query</span>
        <div class="editor-actions">
            <button
                type="button"
                class="btn btn-xs btn-transparent"
                on:click={clear}
                disabled={disabled || loading || !value}
                use:tooltip={"Clear"}
            >
                <i class="ri-delete-bin-line"></i>
            </button>
            <button
                type="button"
                class="btn btn-sm btn-primary"
                on:click={execute}
                disabled={disabled || loading || !value.trim()}
                use:tooltip={"Execute (Ctrl+Enter)"}
            >
                {#if loading}
                    <i class="ri-loader-4-line animate-spin"></i>
                {:else}
                    <i class="ri-play-line"></i>
                {/if}
                Run
            </button>
        </div>
    </div>
    <div class="editor-body">
        <div class="line-numbers" aria-hidden="true">
            <pre>{lineNumbers}</pre>
        </div>
        <textarea
            bind:this={textareaEl}
            bind:value
            on:keydown={handleKeydown}
            on:input={() => dispatch("input", { value })}
            {placeholder}
            {disabled}
            spellcheck="false"
            autocomplete="off"
            autocorrect="off"
            autocapitalize="off"
            class="editor-textarea"
        ></textarea>
    </div>
    <div class="editor-footer">
        <span class="hint">
            <kbd>Ctrl</kbd>+<kbd>Enter</kbd> to execute
            &nbsp;·&nbsp;
            <kbd>Tab</kbd> to indent
        </span>
    </div>
</div>

<style>
    .sql-editor {
        display: flex;
        flex-direction: column;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        overflow: hidden;
    }
    .sql-editor.disabled,
    .sql-editor.loading {
        opacity: 0.7;
        pointer-events: none;
    }
    .editor-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 8px 12px;
        background: var(--baseAlt1Color);
        border-bottom: 1px solid var(--baseAlt2Color);
    }
    .editor-title {
        font-size: 12px;
        font-weight: 600;
        text-transform: uppercase;
        color: var(--txtHintColor);
    }
    .editor-actions {
        display: flex;
        align-items: center;
        gap: 6px;
    }
    .editor-body {
        display: flex;
        min-height: 150px;
        max-height: 300px;
    }
    .line-numbers {
        flex-shrink: 0;
        padding: 12px 0;
        background: var(--baseAlt1Color);
        border-right: 1px solid var(--baseAlt2Color);
        user-select: none;
    }
    .line-numbers pre {
        margin: 0;
        padding: 0 12px;
        font-family: var(--monospaceFontFamily);
        font-size: 13px;
        line-height: 1.5;
        color: var(--txtDisabledColor);
        text-align: right;
    }
    .editor-textarea {
        flex: 1;
        width: 100%;
        padding: 12px;
        border: none;
        background: transparent;
        font-family: var(--monospaceFontFamily);
        font-size: 13px;
        line-height: 1.5;
        color: var(--txtPrimaryColor);
        resize: none;
        overflow: auto;
    }
    .editor-textarea:focus {
        outline: none;
    }
    .editor-textarea::placeholder {
        color: var(--txtDisabledColor);
    }
    .editor-footer {
        padding: 6px 12px;
        background: var(--baseAlt1Color);
        border-top: 1px solid var(--baseAlt2Color);
    }
    .hint {
        font-size: 11px;
        color: var(--txtHintColor);
    }
    .hint kbd {
        display: inline-block;
        padding: 1px 4px;
        font-size: 10px;
        font-family: var(--monospaceFontFamily);
        background: var(--baseColor);
        border: 1px solid var(--baseAlt3Color);
        border-radius: 3px;
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>

