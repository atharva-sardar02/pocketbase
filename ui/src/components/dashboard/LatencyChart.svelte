<script>
    import { onMount, onDestroy } from "svelte";
    import {
        Chart,
        LineElement,
        PointElement,
        LineController,
        LinearScale,
        TimeScale,
        Tooltip,
        Legend,
    } from "chart.js";
    import "chartjs-adapter-luxon";

    export let data = [];
    export let loading = false;

    let chartCanvas;
    let chartInst;

    $: if (chartInst && data) {
        updateChart();
    }

    function updateChart() {
        if (!chartInst) return;

        const chartData = data.map(item => ({
            x: new Date(item.date),
            avg: item.avg,
            p50: item.p50,
            p95: item.p95,
            p99: item.p99
        }));

        chartInst.data.datasets[0].data = chartData.map(d => ({ x: d.x, y: d.p50 }));
        chartInst.data.datasets[1].data = chartData.map(d => ({ x: d.x, y: d.p95 }));
        chartInst.data.datasets[2].data = chartData.map(d => ({ x: d.x, y: d.p99 }));
        chartInst.update();
    }

    onMount(() => {
        Chart.register(LineElement, PointElement, LineController, LinearScale, TimeScale, Tooltip, Legend);

        chartInst = new Chart(chartCanvas, {
            type: "line",
            data: {
                datasets: [
                    {
                        label: "p50",
                        data: [],
                        borderColor: "#10b981",
                        backgroundColor: "transparent",
                        borderWidth: 2,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                        tension: 0.3,
                    },
                    {
                        label: "p95",
                        data: [],
                        borderColor: "#f59e0b",
                        backgroundColor: "transparent",
                        borderWidth: 2,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                        tension: 0.3,
                    },
                    {
                        label: "p99",
                        data: [],
                        borderColor: "#ef4444",
                        backgroundColor: "transparent",
                        borderWidth: 2,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                        tension: 0.3,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: {
                    intersect: false,
                    mode: "index",
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        grid: {
                            color: "rgba(0, 0, 0, 0.05)",
                        },
                        ticks: {
                            callback: function(value) {
                                if (value < 1000) return value + "ms";
                                return (value / 1000).toFixed(1) + "s";
                            },
                            maxTicksLimit: 5,
                        },
                    },
                    x: {
                        type: "time",
                        time: {
                            unit: "hour",
                            displayFormats: {
                                hour: "HH:mm"
                            },
                            tooltipFormat: "MMM d, HH:mm",
                        },
                        grid: {
                            display: false,
                        },
                        ticks: {
                            maxTicksLimit: 8,
                            maxRotation: 0,
                        },
                    },
                },
                plugins: {
                    legend: {
                        display: true,
                        position: "top",
                        align: "end",
                        labels: {
                            boxWidth: 12,
                            padding: 15,
                            usePointStyle: true,
                        },
                    },
                    tooltip: {
                        backgroundColor: "rgba(0, 0, 0, 0.8)",
                        padding: 12,
                        callbacks: {
                            label: function(context) {
                                let value = context.parsed.y;
                                if (value < 1000) return `${context.dataset.label}: ${value.toFixed(1)}ms`;
                                return `${context.dataset.label}: ${(value / 1000).toFixed(2)}s`;
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
        height: 250px;
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
