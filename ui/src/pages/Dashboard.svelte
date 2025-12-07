<script>
    import { onMount, onDestroy } from "svelte";
    import { pageTitle } from "@/stores/app";
    import {
        dashboardPeriod,
        dashboardRefreshInterval,
        overviewLoading,
        requestsLoading,
        latencyLoading,
        errorsLoading,
        endpointsLoading,
        collectionsLoading,
        overviewData,
        requestsData,
        latencyData,
        errorsData,
        endpointsData,
        collectionsData,
        dashboardError,
        formatBytes,
        formatLatency,
        formatNumber,
        formatPercent,
        resetDashboardData,
    } from "@/stores/dashboard";
    import ApiClient from "@/utils/ApiClient";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import MetricCard from "@/components/dashboard/MetricCard.svelte";
    import RequestsChart from "@/components/dashboard/RequestsChart.svelte";
    import LatencyChart from "@/components/dashboard/LatencyChart.svelte";
    import EndpointsChart from "@/components/dashboard/EndpointsChart.svelte";
    import CollectionsTable from "@/components/dashboard/CollectionsTable.svelte";

    $pageTitle = "Dashboard";

    let refreshTimer;
    let lastRefresh = null;

    const periodOptions = [
        { value: "1h", label: "Last 1 hour", interval: "5m" },
        { value: "6h", label: "Last 6 hours", interval: "15m" },
        { value: "24h", label: "Last 24 hours", interval: "1h" },
        { value: "7d", label: "Last 7 days", interval: "6h" },
    ];

    // Get interval based on selected period
    function getIntervalForPeriod(period) {
        const option = periodOptions.find(o => o.value === period);
        return option ? option.interval : "1h";
    }

    onMount(() => {
        resetDashboardData();
        loadAllData();
        startAutoRefresh();
    });

    onDestroy(() => {
        stopAutoRefresh();
    });

    function startAutoRefresh() {
        stopAutoRefresh();
        if ($dashboardRefreshInterval > 0) {
            refreshTimer = setInterval(() => {
                loadAllData();
            }, $dashboardRefreshInterval * 1000);
        }
    }

    function stopAutoRefresh() {
        if (refreshTimer) {
            clearInterval(refreshTimer);
            refreshTimer = null;
        }
    }

    function handlePeriodChange() {
        loadAllData();
    }

    async function loadAllData() {
        lastRefresh = new Date();
        dashboardError.set("");
        
        // Load all data in parallel
        await Promise.all([
            loadOverview(),
            loadRequests(),
            loadLatency(),
            loadErrors(),
            loadEndpoints(),
            loadCollections(),
        ]);
    }

    async function loadOverview() {
        overviewLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const response = await fetch(`${baseUrl}/api/metrics/overview?period=${$dashboardPeriod}`, {
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
            overviewData.set(data);
        } catch (err) {
            console.error("Failed to load overview:", err);
            dashboardError.set("Failed to load overview metrics");
        } finally {
            overviewLoading.set(false);
        }
    }

    async function loadRequests() {
        requestsLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const interval = getIntervalForPeriod($dashboardPeriod);
            const response = await fetch(`${baseUrl}/api/metrics/requests?period=${$dashboardPeriod}&interval=${interval}`, {
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
            requestsData.set(data.data || []);
        } catch (err) {
            console.error("Failed to load requests:", err);
        } finally {
            requestsLoading.set(false);
        }
    }

    async function loadLatency() {
        latencyLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const interval = getIntervalForPeriod($dashboardPeriod);
            const response = await fetch(`${baseUrl}/api/metrics/latency?period=${$dashboardPeriod}&interval=${interval}`, {
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
            latencyData.set(data.data || []);
        } catch (err) {
            console.error("Failed to load latency:", err);
        } finally {
            latencyLoading.set(false);
        }
    }

    async function loadErrors() {
        errorsLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const interval = getIntervalForPeriod($dashboardPeriod);
            const response = await fetch(`${baseUrl}/api/metrics/errors?period=${$dashboardPeriod}&interval=${interval}`, {
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
            errorsData.set(data.data || []);
        } catch (err) {
            console.error("Failed to load errors:", err);
        } finally {
            errorsLoading.set(false);
        }
    }

    async function loadEndpoints() {
        endpointsLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const response = await fetch(`${baseUrl}/api/metrics/endpoints?period=${$dashboardPeriod}&limit=10`, {
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
            endpointsData.set(data.data || []);
        } catch (err) {
            console.error("Failed to load endpoints:", err);
        } finally {
            endpointsLoading.set(false);
        }
    }

    async function loadCollections() {
        collectionsLoading.set(true);
        try {
            const baseUrl = ApiClient.baseURL.endsWith("/")
                ? ApiClient.baseURL.slice(0, -1)
                : ApiClient.baseURL;
            
            const response = await fetch(`${baseUrl}/api/metrics/collections`, {
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
            collectionsData.set(data.data || []);
        } catch (err) {
            console.error("Failed to load collections:", err);
        } finally {
            collectionsLoading.set(false);
        }
    }

    function formatLastRefresh() {
        if (!lastRefresh) return "";
        return lastRefresh.toLocaleTimeString();
    }
</script>

<PageWrapper>
    <div class="dashboard">
        <header class="dashboard-header">
            <h1>
                <i class="ri-dashboard-line"></i>
                Dashboard
            </h1>
            <div class="header-actions">
                <select 
                    class="period-select" 
                    bind:value={$dashboardPeriod} 
                    on:change={handlePeriodChange}
                >
                    {#each periodOptions as option}
                        <option value={option.value}>{option.label}</option>
                    {/each}
                </select>
                <button 
                    type="button" 
                    class="btn btn-secondary btn-sm"
                    on:click={loadAllData}
                    title="Refresh data"
                >
                    <i class="ri-refresh-line"></i>
                    Refresh
                </button>
                {#if lastRefresh}
                    <span class="last-refresh">
                        Updated: {formatLastRefresh()}
                    </span>
                {/if}
            </div>
        </header>

        {#if $dashboardError}
            <div class="alert alert-danger">
                <i class="ri-error-warning-line"></i>
                {$dashboardError}
            </div>
        {/if}

        <!-- Overview Cards -->
        <section class="metrics-grid">
            <MetricCard
                title="Total Requests"
                value={formatNumber($overviewData.totalRequests)}
                subtitle={$dashboardPeriod}
                icon="ri-send-plane-line"
                iconColor="primary"
                loading={$overviewLoading}
            />
            <MetricCard
                title="Avg Latency"
                value={formatLatency($overviewData.avgLatency)}
                subtitle="Response time"
                icon="ri-timer-line"
                iconColor="success"
                loading={$overviewLoading}
            />
            <MetricCard
                title="Error Rate"
                value={formatPercent($overviewData.errorRate)}
                subtitle="{formatNumber($overviewData.totalErrors)} errors"
                icon="ri-bug-line"
                iconColor={$overviewData.errorRate > 5 ? "danger" : $overviewData.errorRate > 1 ? "warning" : "success"}
                loading={$overviewLoading}
            />
            <MetricCard
                title="Database Size"
                value={formatBytes($overviewData.databaseSize)}
                subtitle="SQLite storage"
                icon="ri-database-2-line"
                iconColor="info"
                loading={$overviewLoading}
            />
        </section>

        <!-- Charts Row -->
        <section class="charts-row">
            <div class="chart-card">
                <h3>Requests Over Time</h3>
                <RequestsChart data={$requestsData} loading={$requestsLoading} />
            </div>
            <div class="chart-card">
                <h3>Latency Percentiles</h3>
                <LatencyChart data={$latencyData} loading={$latencyLoading} />
            </div>
        </section>

        <!-- Bottom Row -->
        <section class="bottom-row">
            <div class="chart-card">
                <h3>Top Endpoints</h3>
                <EndpointsChart data={$endpointsData} loading={$endpointsLoading} />
            </div>
            <div class="chart-card">
                <h3>Collections</h3>
                <CollectionsTable data={$collectionsData} loading={$collectionsLoading} />
            </div>
        </section>
    </div>
</PageWrapper>

<style>
    .dashboard {
        padding: 20px;
        max-width: 1600px;
        margin: 0 auto;
        min-width: 0;
    }
    .dashboard-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 24px;
        flex-wrap: wrap;
        gap: 12px;
    }
    .dashboard-header h1 {
        display: flex;
        align-items: center;
        gap: 10px;
        margin: 0;
        font-size: 1.5em;
    }
    .dashboard-header h1 i {
        color: var(--primaryColor);
    }
    .header-actions {
        display: flex;
        align-items: center;
        gap: 12px;
        flex-wrap: wrap;
    }
    .period-select {
        padding: 8px 12px;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        background: var(--baseColor);
        font-size: 13px;
        cursor: pointer;
    }
    .last-refresh {
        font-size: 12px;
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
    .metrics-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 16px;
        margin-bottom: 24px;
    }
    .charts-row,
    .bottom-row {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 20px;
        margin-bottom: 24px;
    }
    .chart-card {
        background: var(--baseColor);
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
        padding: 20px;
        min-width: 0;
        overflow: hidden;
    }
    .chart-card h3 {
        margin: 0 0 16px;
        font-size: 14px;
        font-weight: 600;
        color: var(--txtPrimaryColor);
    }

    @media (max-width: 1200px) {
        .charts-row,
        .bottom-row {
            grid-template-columns: 1fr;
        }
    }

    @media (max-width: 768px) {
        .dashboard {
            padding: 12px;
        }
        .dashboard-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 16px;
        }
        .header-actions {
            flex-wrap: wrap;
        }
        .metrics-grid {
            grid-template-columns: 1fr;
        }
        .chart-card {
            padding: 16px;
        }
    }
</style>
