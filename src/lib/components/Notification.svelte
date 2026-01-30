<script lang="ts">
  export let message: string;
  export let type: 'success' | 'error' = 'success';
</script>

<div class="toast toast-top toast-end">
  <div class="alert alert-{type}">
    <div>
      <span>{message}</span>
    </div>
  </div>
</div>
