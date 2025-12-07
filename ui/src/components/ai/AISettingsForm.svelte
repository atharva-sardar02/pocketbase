<script>
    import { createEventDispatcher } from "svelte";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import Toggler from "@/components/base/Toggler.svelte";

    export let settings = {};
    export let isTesting = false;
    export let testError = null;

    const dispatch = createEventDispatcher();

    const providers = [
        { label: "OpenAI", value: "openai" },
        { label: "Ollama", value: "ollama" },
        { label: "Anthropic", value: "anthropic" },
        { label: "Custom", value: "custom" },
    ];

    const defaultBaseURLs = {
        openai: "https://api.openai.com/v1",
        ollama: "http://localhost:11434/v1",
        anthropic: "https://api.anthropic.com/v1",
        custom: "",
    };

    const models = {
        openai: [
            { label: "gpt-4o-mini", value: "gpt-4o-mini" },
            { label: "gpt-4o", value: "gpt-4o" },
            { label: "gpt-3.5-turbo", value: "gpt-3.5-turbo" },
        ],
        ollama: [
            { label: "llama2", value: "llama2" },
            { label: "llama3", value: "llama3" },
            { label: "mistral", value: "mistral" },
        ],
        anthropic: [
            { label: "claude-3-5-sonnet-20241022", value: "claude-3-5-sonnet-20241022" },
            { label: "claude-3-opus-20240229", value: "claude-3-opus-20240229" },
        ],
        custom: [],
    };

    $: showApiKey = settings.provider !== "ollama";
    $: providerModels = models[settings.provider] || [];
    $: isCustomProvider = settings.provider === "custom";

    function handleProviderChange(e) {
        const newProvider = e.detail.value;
        settings.provider = newProvider;
        if (defaultBaseURLs[newProvider]) {
            settings.baseUrl = defaultBaseURLs[newProvider];
        }
        // Reset model if not available for new provider
        if (providerModels.length > 0 && !providerModels.find((m) => m.value === settings.model)) {
            settings.model = providerModels[0]?.value || "";
        }
        dispatch("change");
    }

    function handleFieldChange() {
        dispatch("change");
    }

    function handleTestConnection() {
        dispatch("test");
    }
</script>

<div class="ai-settings-form">
    <Field class="form-field form-field-toggle" name="ai.enabled" let:uniqueId>
        <input
            type="checkbox"
            id={uniqueId}
            bind:checked={settings.enabled}
            on:change={handleFieldChange}
        />
        <label for={uniqueId}>Enable AI Query</label>
    </Field>

    {#if settings.enabled}
        <Field class="form-field required" name="ai.provider" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Provider</span>
            </label>
            <ObjectSelect
                id={uniqueId}
                items={providers}
                bind:keyOfSelected={settings.provider}
                on:change={handleProviderChange}
            />
        </Field>

        <Field class="form-field required" name="ai.baseUrl" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">API Base URL</span>
            </label>
            <input
                type="text"
                id={uniqueId}
                bind:value={settings.baseUrl}
                placeholder="https://api.openai.com/v1"
                on:input={handleFieldChange}
            />
        </Field>

        {#if showApiKey}
            <Field class="form-field required" name="ai.apiKey" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">API Key</span>
                </label>
                <RedactedPasswordInput
                    id={uniqueId}
                    bind:value={settings.apiKey}
                    placeholder="Enter your API key"
                    on:input={handleFieldChange}
                />
            </Field>
        {/if}

        <Field class="form-field required" name="ai.model" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Model</span>
            </label>
            {#if isCustomProvider || providerModels.length === 0}
                <input
                    type="text"
                    id={uniqueId}
                    bind:value={settings.model}
                    placeholder="e.g., gpt-4o-mini"
                    on:input={handleFieldChange}
                />
            {:else}
                <ObjectSelect
                    id={uniqueId}
                    items={providerModels}
                    bind:keyOfSelected={settings.model}
                    on:change={handleFieldChange}
                />
            {/if}
        </Field>

        <Field class="form-field" name="ai.temperature" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Temperature</span>
                <span class="txt-hint">({settings.temperature || 0.1})</span>
            </label>
            <input
                type="range"
                id={uniqueId}
                min="0"
                max="1"
                step="0.1"
                bind:value={settings.temperature}
                on:input={(e) => {
                    settings.temperature = parseFloat(e.target.value);
                    handleFieldChange();
                }}
            />
            <div class="form-hint">
                Controls randomness: 0.0 (deterministic) to 1.0 (creative). Recommended: 0.1 for
                consistent filter generation.
            </div>
        </Field>

        <div class="form-field">
            <button
                type="button"
                class="btn btn-primary"
                on:click={handleTestConnection}
                disabled={isTesting || !settings.baseUrl || (showApiKey && !settings.apiKey)}
            >
                {#if isTesting}
                    <i class="ri-loader-4-line spin" aria-hidden="true"></i>
                    <span class="txt">Testing...</span>
                {:else}
                    <i class="ri-checkbox-circle-line" aria-hidden="true"></i>
                    <span class="txt">Test Connection</span>
                {/if}
            </button>
        </div>

        {#if testError}
            <div class="alert alert-error m-t-sm">
                <i class="ri-error-warning-line" aria-hidden="true"></i>
                <span class="txt">{testError}</span>
            </div>
        {/if}
    {/if}
</div>

<style>
    .ai-settings-form {
        padding: 20px 0;
    }

    input[type="range"] {
        width: 100%;
    }
</style>

