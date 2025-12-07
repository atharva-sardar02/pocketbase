<script>
    export let title = "";
    export let value = "";
    export let subtitle = "";
    export let icon = "";
    export let iconColor = "primary";
    export let loading = false;
</script>

<div class="metric-card" class:loading>
    <div class="metric-icon {iconColor}">
        <i class={icon}></i>
    </div>
    <div class="metric-content">
        <div class="metric-title">{title}</div>
        {#if loading}
            <div class="metric-value skeleton"></div>
        {:else}
            <div class="metric-value">{value}</div>
        {/if}
        {#if subtitle}
            <div class="metric-subtitle">{subtitle}</div>
        {/if}
    </div>
</div>

<style>
    .metric-card {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        transition: box-shadow 0.2s;
    }
    .metric-card:hover {
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }
    .metric-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 48px;
        height: 48px;
        border-radius: 12px;
        font-size: 24px;
    }
    .metric-icon.primary {
        background: var(--primaryAltColor);
        color: var(--primaryColor);
    }
    .metric-icon.success {
        background: rgba(var(--successColorRaw), 0.1);
        color: var(--successColor);
    }
    .metric-icon.warning {
        background: rgba(var(--warningColorRaw), 0.1);
        color: var(--warningColor);
    }
    .metric-icon.danger {
        background: rgba(var(--dangerColorRaw), 0.1);
        color: var(--dangerColor);
    }
    .metric-icon.info {
        background: rgba(var(--infoColorRaw), 0.1);
        color: var(--infoColor);
    }
    .metric-content {
        flex: 1;
        min-width: 0;
    }
    .metric-title {
        font-size: 12px;
        font-weight: 500;
        color: var(--txtHintColor);
        text-transform: uppercase;
        letter-spacing: 0.5px;
        margin-bottom: 4px;
    }
    .metric-value {
        font-size: 24px;
        font-weight: 600;
        color: var(--txtPrimaryColor);
        line-height: 1.2;
    }
    .metric-value.skeleton {
        width: 80px;
        height: 28px;
        background: var(--baseAlt1Color);
        border-radius: 4px;
        animation: pulse 1.5s infinite;
    }
    .metric-subtitle {
        font-size: 12px;
        color: var(--txtSecondaryColor);
        margin-top: 4px;
    }
    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }
</style>
