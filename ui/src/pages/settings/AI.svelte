<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import AISettingsForm from "@/components/ai/AISettingsForm.svelte";

    $pageTitle = "AI Query settings";

    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let isTesting = false;
    let testError = null;
    let testSuccess = false;

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
            init(settings);
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    async function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        try {
            // Ensure temperature is a number
            const aiSettings = {
                ...formSettings.ai,
                temperature: typeof formSettings.ai.temperature === "string"
                    ? parseFloat(formSettings.ai.temperature)
                    : formSettings.ai.temperature,
            };

            const settings = await ApiClient.settings.update(
                CommonHelper.filterRedactedProps({
                    ai: aiSettings,
                })
            );
            init(settings);

            setErrors({});
            addSuccessToast("Successfully saved AI Query settings.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        formSettings = {
            ai: {
                enabled: settings?.ai?.enabled || false,
                provider: settings?.ai?.provider || "openai",
                baseUrl: settings?.ai?.baseUrl || "https://api.openai.com/v1",
                apiKey: settings?.ai?.apiKey || "",
                model: settings?.ai?.model || "gpt-4o-mini",
                temperature: typeof settings?.ai?.temperature === "number" 
                    ? settings.ai.temperature 
                    : parseFloat(settings?.ai?.temperature) || 0.1,
            },
        };

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
        testError = null;
        testSuccess = false;
    }

    function reset() {
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
        testError = null;
        testSuccess = false;
    }

    async function testConnection() {
        if (isTesting) {
            return;
        }

        // Validate required fields
        if (!formSettings.ai.baseUrl) {
            testError = "API Base URL is required";
            return;
        }

        if (formSettings.ai.provider !== "ollama" && !formSettings.ai.apiKey) {
            testError = "API Key is required for this provider";
            return;
        }

        if (!formSettings.ai.model) {
            testError = "Model is required";
            return;
        }

        isTesting = true;
        testError = null;
        testSuccess = false;

        try {
            // Temporarily save settings to test connection
            // We need to enable AI temporarily for the test
            const tempSettings = JSON.parse(JSON.stringify(originalFormSettings));
            const testAISettings = {
                ...formSettings.ai,
                enabled: true, // Temporarily enable for test
            };
            
            // Update settings temporarily for testing
            await ApiClient.settings.update(
                CommonHelper.filterRedactedProps({
                    ai: testAISettings,
                })
            );

            // Create a test request to the AI query endpoint
            // Use a simple query that should always generate a valid filter
            // "all records" should translate to an empty filter or id != ""
            const testBody = {
                collection: "_superusers", // Use a system collection that always exists
                query: "all records", // Simple query that should work - translates to empty filter or id != ""
                execute: false, // Don't execute, just test the LLM connection
            };

            // Ensure proper URL construction (avoid double slashes)
            const baseUrl = ApiClient.baseURL.endsWith('/') 
                ? ApiClient.baseURL.slice(0, -1) 
                : ApiClient.baseURL;
            const url = `${baseUrl}/api/ai/query`;
            const fetchOptions = {
                method: "POST",
                body: JSON.stringify(testBody),
                headers: {
                    "Content-Type": "application/json",
                    Authorization: ApiClient.authStore.token ? `Bearer ${ApiClient.authStore.token}` : "",
                },
            };
            const response = await fetch(url, fetchOptions);

            if (response.ok) {
                testSuccess = true;
                addSuccessToast("Connection test successful!");
                // If test succeeded, keep the settings (user might want to save)
            } else {
                const errorData = await response.json().catch(() => ({}));
                testError = errorData.message || `Connection test failed: ${response.status}`;
                // Restore original settings if test failed
                await ApiClient.settings.update(
                    CommonHelper.filterRedactedProps({
                        ai: tempSettings.ai,
                    })
                );
            }
        } catch (err) {
            testError = err.message || "Failed to test connection. Please check your settings.";
            // Try to restore original settings on error
            try {
                await ApiClient.settings.update(
                    CommonHelper.filterRedactedProps({
                        ai: originalFormSettings.ai,
                    })
                );
            } catch (restoreErr) {
                console.error("Failed to restore settings after test error:", restoreErr);
            }
        } finally {
            isTesting = false;
        }
    }

    function handleFormChange() {
        testError = null;
        testSuccess = false;
    }
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={() => save()}>
            <div class="content txt-xl m-b-base">
                <p>Configure the AI Query feature to enable natural language queries on your collections.</p>
                <p>
                    The AI will translate your natural language queries into PocketBase filter expressions.
                </p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <AISettingsForm
                    bind:settings={formSettings.ai}
                    {isTesting}
                    {testError}
                    on:change={handleFormChange}
                    on:test={testConnection}
                />

                <div class="panel-footer">
                    <button
                        type="button"
                        class="btn btn-transparent"
                        on:click={reset}
                        disabled={!hasChanges || isSaving}
                    >
                        <span class="txt">Reset</span>
                    </button>
                    <button type="submit" class="btn btn-primary" disabled={!hasChanges || isSaving}>
                        {#if isSaving}
                            <i class="ri-loader-4-line spin" aria-hidden="true"></i>
                            <span class="txt">Saving...</span>
                        {:else}
                            <i class="ri-save-line" aria-hidden="true"></i>
                            <span class="txt">Save</span>
                        {/if}
                    </button>
                </div>
            {/if}
        </form>
    </div>
</PageWrapper>

