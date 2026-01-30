<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
</script>
<div class="modal modal-open">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Add new peer</h3>
    <div class="py-4">
        <div class="form-control">
            <label class="label" for="name">
                <span class="label-text">Name</span>
            </label>
            <input type="text" id="name" class="input input-bordered">
        </div>
    </div>
    <div class="modal-action">
      <button class="btn" on:click={() => dispatch('close')}>Cancel</button>
      <button class="btn btn-primary" on:click={() => dispatch('close')}>Save</button>
    </div>
  </div>
</div>
