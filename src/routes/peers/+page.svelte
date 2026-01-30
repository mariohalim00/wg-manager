<script lang="ts">
  import PeerTable from '$lib/components/PeerTable.svelte';
  import PeerModal from '$lib/components/PeerModal.svelte';
  import { peers } from '$lib/stores/peers';
  import type { Peer } from '$lib/types';

  let showModal = false;

  const dummyPeers: Peer[] = [
    { id: '1', name: 'Peer 1', publicKey: '...', status: 'online', allowedIps: ['10.0.0.1/32'], latestHandshake: '...', transfer: { received: 0, sent: 0 }},
    { id: '2', name: 'Peer 2', publicKey: '...', status: 'offline', allowedIps: ['10.0.0.2/32'], latestHandshake: '...', transfer: { received: 0, sent: 0 }},
  ];

  peers.set(dummyPeers);
</script>

<div class="p-4">
  <div class="flex justify-between items-center mb-4">
    <h1 class="text-2xl font-bold">Peers</h1>
    <button class="btn btn-primary" on:click={() => showModal = true}>Add Peer</button>
  </div>

  <PeerTable peers={$peers} />

  {#if showModal}
    <PeerModal on:close={() => showModal = false} />
  {/if}
</div>
