<script lang="ts">
  import { stats } from '$lib/stores/stats';
  import type { Stats } from '$lib/stores/stats';

  const dummyStats: Stats = {
    totalPeers: 2,
    onlinePeers: 1,
    totalDataUsage: {
      sent: 1024,
      received: 2048,
    }
  };

  stats.set(dummyStats);
</script>

<div class="p-4">
  <h1 class="text-2xl font-bold mb-4">Statistics</h1>
  <div class="stats shadow">
    <div class="stat">
      <div class="stat-title">Total Peers</div>
      <div class="stat-value">{$stats.totalPeers}</div>
    </div>
    <div class="stat">
      <div class="stat-title">Online Peers</div>
      <div class="stat-value">{$stats.onlinePeers}</div>
    </div>
    <div class="stat">
      <div class="stat-title">Total Data Sent</div>
      <div class="stat-value">{$stats.totalDataUsage.sent} KB</div>
    </div>
    <div class="stat">
      <div class="stat-title">Total Data Received</div>
      <div class="stat-value">{$stats.totalDataUsage.received} KB</div>
    </div>
  </div>
</div>
