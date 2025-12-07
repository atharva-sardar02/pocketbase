<script>
    import { onMount, onDestroy } from "svelte";
    import {
        Chart,
        BarElement,
        BarController,
        CategoryScale,
        LinearScale,
        Tooltip,
    } from "chart.js";

    export let data = [];
    export let loading = false;

    let chartCanvas;
    let chartInst;

    $: if (chartInst && data) {
        updateChart();
    }

    function truncateEndpoint(endpoint, maxLen = 40) {
        if (!endpoint) return "";
        if (endpoint.length <= maxLen) return endpoint;
        return endpoint.slice(0, maxLen - 3) + "...";
    }

    function updateChart() {
        if (!chartInst) return;

        const labels = data.map(d => truncateEndpoint(d.endpoint));
        const counts = data.map(d => d.count);

        chartInst.data.labels = labels;
        chartInst.data.datasets[0].data = counts;
        chartInst.update();
    }

    onMount(() => {
        Chart.register(BarElement, BarController, CategoryScale, LinearScale, Tooltip);

        chartInst = new Chart(chartCanvas, {
            type: "bar",
            data: {
                labels: [],
                datasets: [
                    {
                        label: "Requests",
                        data: [],
                        backgroundColor: "rgba(99, 102, 241, 0.8)",
                        borderColor: "#6366f1",
                        borderWidth: 1,
                        borderRadius: 4,
                    },
                ],
            },
            options: {
                indexAxis: "y",
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        beginAtZero: true,
                        grid: {
                            color: "rgba(0, 0, 0, 0.05)",
                        },
                        ticks: {
                            precision: 0,
                        },
                    },
                    y: {
                        grid: {
                            display: false,
                        },
                        ticks: {
                            font: {
                                family: "var(--monospaceFontFamily)",
                                size: 11,
                            },
                        },
                    },
                },
                plugins: {
                    legend: {
                        display: false,
                    },
                    tooltip: {
                        backgroundColor: "rgba(0, 0, 0, 0.8)",
                        padding: 12,
                        callbacks: {
                            title: function(context) {
                                const idx = context[0].dataIndex;
                                return data[idx]?.endpoint || "";
                            },
                            label: function(context) {
                                const idx = context.dataIndex;
                                const item = data[idx];
                                return [
                                    `Requests: ${context.parsed.x.toLocaleString()}`,
                                    `Avg Latency: ${item?.avgLatency?.toFixed(1) || 0}ms`
                                ];
                            },
                        },
                    },
                },
            },
        });

        updateChart();
    });

    onDestroy(() => {
        chartInst?.destroy();
    });
</script>

<div class="chart-container" class:loading>
    {#if loading}
        <div class="chart-loader">
            <div class="loader"></div>
        </div>
    {/if}
    <canvas bind:this={chartCanvas}></canvas>
</div>

<style>
    .chart-container {
        position: relative;
        width: 100%;
        height: 300px;
    }
    .chart-container.loading canvas {
        opacity: 0.5;
    }
    .chart-loader {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        z-index: 10;
    }
</style>
