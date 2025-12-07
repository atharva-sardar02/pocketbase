<script>
    import { onMount, onDestroy } from "svelte";
    import {
        Chart,
        LineElement,
        PointElement,
        LineController,
        LinearScale,
        TimeScale,
        Filler,
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

        chartInst.data.datasets[0].data = data.map(item => ({
            x: new Date(item.date),
            y: item.total
        }));
        chartInst.update();
    }

    onMount(() => {
        Chart.register(LineElement, PointElement, LineController, LinearScale, TimeScale, Filler, Tooltip, Legend);

        chartInst = new Chart(chartCanvas, {
            type: "line",
            data: {
                datasets: [
                    {
                        label: "Requests",
                        data: [],
                        borderColor: "#6366f1",
                        backgroundColor: "rgba(99, 102, 241, 0.1)",
                        borderWidth: 2,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                        fill: true,
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
                            precision: 0,
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
                        display: false,
                    },
                    tooltip: {
                        backgroundColor: "rgba(0, 0, 0, 0.8)",
                        padding: 12,
                        titleFont: { size: 12 },
                        bodyFont: { size: 13 },
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
