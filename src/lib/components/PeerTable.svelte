<script lang="ts">
  import type { Peer } from '$lib/types';
  import StatusBadge from '$lib/components/StatusBadge.svelte';
  export let peers: Peer[];
</script>

<div class="overflow-x-auto">
  <table class="table w-full">
    <thead>
      <tr>
        <th>Name</th>
        <th>Allowed IPs</th>
        <th>Status</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each peers as peer}
        <tr>
          <td>{peer.name}</td>
          <td>{peer.allowedIps.join(', ')}</td>
          <td><StatusBadge status={peer.status} /></td>
          <td>
            <button class="btn btn-sm btn-ghost">Edit</button>
            <button class="btn btn-sm btn-ghost">Delete</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
