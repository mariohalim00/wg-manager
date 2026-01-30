<script lang="ts">
    export let status: 'online' | 'offline';
</script>

{#if status === 'online'}
    <div class="badge badge-success">Online</div>
{:else}
    <div class="badge badge-error">Offline</div>
{/if}
