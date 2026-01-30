<script lang="ts">
    import { stats } from '$lib/stores/stats';
    import { onMount } from 'svelte';
    import type { Stats } from '$lib/stores/stats';

    onMount(() => {
        const dummyStats: Stats = {
            totalPeers: 2,
            onlinePeers: 1,
            totalDataUsage: {
                sent: 1024,
                received: 2048,
            }
        };
        stats.set(dummyStats);
    });
</script>

<div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Dashboard</h1>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Total Peers</h2>
                <p class="text-4xl font-bold">{$stats.totalPeers}</p>
            </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Online Peers</h2>
                <p class="text-4xl font-bold">{$stats.onlinePeers}</p>
            </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Data Sent</h2>
                <p class="text-4xl font-bold">{$stats.totalDataUsage.sent} KB</p>
            </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Data Received</h2>
                <p class="text-4xl font-bold">{$stats.totalDataUsage.received} KB</p>
            </div>
        </div>
    </div>
</div>
